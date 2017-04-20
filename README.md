# Launchcontrol (work in progress)

[![GoDoc](https://godoc.org/github.com/funjack/launchcontrol?status.svg)](https://godoc.org/github.com/funjack/launchcontrol)
[![Go Report Card](https://goreportcard.com/badge/github.com/funjack/launchcontrol)](https://goreportcard.com/report/github.com/funjack/launchcontrol)

Launchcontrol is a server that can control a Launch, and is meant to be used
with a plugin for an external player (eg Kodi or VLC)

The goal is to support multiple haptics protocols/formats.

## Setup

See the [gatt docs](https://godoc.org/github.com/currantlabs/gatt#hdr-SETUP)
for the Bluetooth requirements/setup.

## Build

```sh
go get ./...
go build
sudo setcap 'cap_net_raw,cap_net_admin=eip' ./launchcontrol
```

To cross-compile for a Raspberry Pi 2 use `GOARCH=arm GOARM=7 go build`.

## Usage

```sh
# Start server (listening on localhost:6969 by default)
./launchcontrol
```

```sh
# Load and play script
curl -XPOST --data-ascii "{0.50:1,1.00:4,1.15:0,2.00:2}" \
	http://localhost:6969/v1/play
# Pause a playing script
curl http://localhost:6969/v1/pause
# Resume paused script
curl http://localhost:6969/v1/resume
# Jump to a position in the script
curl http://localhost:6969/v1/skip\?p=1m3s
# Stop and reset script
curl http://localhost:6969/v1/stop
# Start playing last loaded script
curl http://localhost:6969/v1/play
```

## Other examples

### Build and run example (Linux)

```sh
go build contrib/examples/playkiiroo.go
sudo setcap 'cap_net_raw,cap_net_admin=eip' ./playkiiroo
./playkiiroo -file input.txt
```

Launchcontrol is released under a [BSD-style license](./LICENSE).
