# asensor

Inspect Android sensors from the command line. Lists available sensors, shows their static metadata, and reads a single event from any of them.

Works on Android (including Termux). Uses `purego` to call `libandroid.so` at runtime, so no cgo toolchain is needed.

## Build

```sh
go build -o asensor .
```

## Usage

```sh
# List every sensor (index and name)
asensor list

# Full details per sensor
asensor list -verbose

# Metadata for one sensor
asensor show 0

# Read one event from a sensor (5s default timeout)
asensor read 0

# Longer timeout for sensors that only fire on real-world events
asensor read -timeout 30s 40

# Stream events for 10 seconds (e.g. watch the flicker sensor)
asensor watch -duration 10s 31
```

Indices come from `asensor list` and are zero-based.