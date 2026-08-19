# asensor

Inspect Android sensors from the command line. Lists available sensors, shows their static metadata, and reads a single event from any of them.

Works on Android (including Termux). Uses `purego` to call `libandroid.so` at runtime.

## Quick start

Run directly from GitHub without installing:

```sh
go run github.com/PeronGH/asensor@latest list
```

## Build

```sh
GOOS=android GOARCH=arm64 go build -o asensor .
```

## Usage

```sh
# List every sensor (index and name)
asensor list

# Full details per sensor
asensor list -verbose

# Metadata for one sensor
asensor show 0

# Read one event (waits forever by default; stop with Ctrl+C)
asensor read 0

# Or wait a limited time for sensors that fire only on real-world events
asensor read -timeout 30s 40

# Stream indefinitely (default; stop with Ctrl+C)
asensor watch 31

# Stream for a limited time (e.g. watch the flicker sensor)
asensor watch -duration 10s 31
```

Indices come from `asensor list` and are zero-based.