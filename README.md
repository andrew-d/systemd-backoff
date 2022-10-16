## systemd-backoff

systemd-backoff is a tool that can be used to add exponential backoff
capability to a systemd unit. It expects to be run as part of a systemd unit,
and will use `sd_pid_get_unit` to fetch the current unit name, and then fetches
the `NRestarts` property that describes the number of restarts for the current
unit. Once it has this value, it uses that as the input to wait an amount of
time configured by the exponential backoff parameters command line arguments.

### Arguments

```
Usage of systemd-backoff:
  -debug
    	print debug messages
  -factor float
    	multiplication factor for each attempt (default 1.5)
  -jitter
    	randomize backoff steps
  -max duration
    	maximum backoff duration (default 10s)
  -min duration
    	minimum backoff duration (default 100ms)
```

There are two environment variables that can be set for debugging:

* `SYSTEMD_BACKOFF_UNIT_NAME`, if set to a string, specifies the current unit name (and does not detect it at runtime).
* `SYSTEMD_BACKOFF_DEBUG_RESTARTS`, if set to an integer, fakes the number of restarts (`NRestarts`) to the provided number.

### Example Usage

In a systemd unit file; using the `+` sigil to run the backoff script as root
so it can communicate with systemd.
```
ExecStartPre=+/usr/local/bin/systemd-backoff -max 30s -factor 1.5
```
