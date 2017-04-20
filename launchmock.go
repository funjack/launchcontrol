package main

import (
	"context"
	"log"
)

// launchMock implements the Launch interface but only logs called methods.
type launchMock struct {
	DisFunc func()
}

func (f launchMock) Move(position, speed int) {
	log.Printf("Move called: Position=%d, Speed=%d", position, speed)
}
func (f launchMock) Connect(ctx context.Context) error {
	log.Printf("Connect called")
	return nil
}
func (f launchMock) Disconnect() {
	log.Printf("Disconnect called")
}
func (f *launchMock) HandleDisconnect(fnc func()) {
	log.Printf("HandleFunc called")
	f.DisFunc = fnc
}
