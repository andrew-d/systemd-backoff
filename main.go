package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/util"
	"github.com/jpillora/backoff"
)

var (
	debug = flag.Bool("debug", false, "print debug messages")

	min    = flag.Duration("min", 100*time.Millisecond, "minimum backoff duration")
	max    = flag.Duration("max", 10*time.Second, "maximum backoff duration")
	factor = flag.Float64("factor", 1.5, "multiplication factor for each attempt")
	jitter = flag.Bool("jitter", false, "randomize backoff steps")
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("systemd-backoff: ")
	log.SetFlags(0)
	flag.Parse()

	var (
		unit string
		err  error
	)
	if ss, ok := os.LookupEnv("SYSTEMD_BACKOFF_UNIT_NAME"); ok {
		unit = ss
	} else {
		unit, err = util.CurrentUnitName()
		if err != nil {
			log.Fatalf("error getting current unit name: %v", err)
		}
	}
	if *debug {
		log.Printf("current unit: %s", unit)
	}

	var restarts uint32
	if ss, ok := os.LookupEnv("SYSTEMD_BACKOFF_DEBUG_RESTARTS"); ok {
		var i uint64
		i, err = strconv.ParseUint(ss, 10, 32)
		if err == nil {
			restarts = uint32(i)
		}
	} else {
		restarts, err = getRestarts(unit)
	}
	if err != nil {
		log.Fatalf("error getting restarts: %v", err)
	}
	if *debug {
		log.Printf("NRestarts: %d", restarts)
	}

	// If this is our first start, then just continue.
	if restarts == 0 {
		return
	}

	// If we have more than one restart, then calculate the sleep and do it.
	b := backoff.Backoff{
		Factor: *factor,
		Jitter: *jitter,
		Min:    *min,
		Max:    *max,
	}
	dur := b.ForAttempt(float64(restarts))
	log.Printf("waiting for: %v", dur)
	time.Sleep(dur)
}

func getRestarts(unit string) (uint32, error) {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dbus.NewSystemConnectionContext(): %w", err)
	}
	defer conn.Close()

	prop, err := conn.GetServicePropertyContext(ctx, unit, "NRestarts")
	if err != nil {
		return 0, fmt.Errorf(`GetServicePropertyContext(%q, "NRestarts"): %w`, unit, err)
	}

	v := prop.Value.Value()
	if ii, ok := v.(uint32); ok {
		return ii, nil
	}

	return 0, fmt.Errorf("unknown property type %T", v)
}
