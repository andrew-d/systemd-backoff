## systemd-backoff

systemd-backoff is a tool that can be used to add exponential backoff
capability to a systemd unit. It fetches the `NRestarts` property for the
current unit file and then uses that as a input to wait an amount of time
configured by the exponential backoff parameters command line arguments.
