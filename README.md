# Launchcontrol (work in progress)

[![GoDoc](https://godoc.org/github.com/funjack/launchcontrol?status.svg)](https://godoc.org/github.com/funjack/launchcontrol)
[![Go Report Card](https://goreportcard.com/badge/github.com/funjack/launchcontrol)](https://goreportcard.com/report/github.com/funjack/launchcontrol)

Launchcontrol is a server that can control a Launch, and is meant to be used
with a plugin for an external player (eg Kodi or VLC)

The goal is to support multiple haptics protocols/formats.

## Setup

See the [gatt docs](https://godoc.org/github.com/currantlabs/gatt#hdr-SETUP)
for the Bluetooth requirements/setup.

## Examples

### Build and run example (Linux)

```sh
go build contrib/examples/playkiiroo.go
sudo setcap 'cap_net_raw,cap_net_admin=eip' ./playkiiroo
./playkiiroo -file input.txt
```

Launchcontrol is released under a [BSD-style license](./LICENSE).
