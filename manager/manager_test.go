package manager

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

type fakeLaunch struct {
	sync.Mutex

	DisFunc func()

	MoveCount       int
	ConnectCount    int
	DisconnectCount int
}

func (f *fakeLaunch) Move(position, speed int) {
	f.Lock()
	defer f.Unlock()
	f.MoveCount++
}
func (f *fakeLaunch) Connect(ctx context.Context) error {
	f.Lock()
	defer f.Unlock()
	f.ConnectCount++
	return nil
}
func (f *fakeLaunch) Disconnect() {
	f.Lock()
	defer f.Unlock()
	f.DisconnectCount++
}
func (f *fakeLaunch) HandleDisconnect(fnc func()) {
	f.Lock()
	defer f.Unlock()
	f.DisFunc = fnc
}

// TestManager is a basic test running through most of managers functions.
func TestManager(t *testing.T) {
	fake := &fakeLaunch{}
	lm := NewLaunchManager(fake)
	p := protocol.NewTimedActionsPlayer()
	p.Script = []protocol.TimedAction{
		{
			Action: protocol.Action{
				Position: 5,
				Speed:    50,
			},
			Time: time.Millisecond * 50,
		},
		{
			Action: protocol.Action{
				Position: 50,
				Speed:    40,
			},
			Time: time.Millisecond * 100,
		},
		{
			Action: protocol.Action{
				Position: 90,
				Speed:    90,
			},
			Time: time.Millisecond * 150,
		},
		{
			Action: protocol.Action{
				Position: 30,
				Speed:    30,
			},
			Time: time.Millisecond * 200,
		},
	}
	lm.SetScriptPlayer(p)

	if err := lm.Play(); err != nil {
		t.Error(err)
	}
	if err := lm.Skip(time.Millisecond * 5); err != nil {
		t.Error(err)
	}
	if err := lm.Pause(); err != nil {
		t.Error(err)
	}
	if err := lm.Resume(); err != nil {
		t.Error(err)
	}

	// Fake a disconnect
	fake.DisFunc()
	if err := lm.Pause(); err == nil {
		t.Error("player should have stopped on disconnect")
	}

	<-time.After(time.Millisecond * 500)
	if err := lm.Play(); err != nil {
		t.Error(err)
	}

	// Give some time for the scenario to play
	<-time.After(time.Millisecond * 500)

	fake.Lock()
	defer fake.Unlock()
	if fake.MoveCount != len(p.Script) {
		t.Errorf("launch did not move enough")
	}
	if fake.DisconnectCount > 0 {
		t.Errorf("launch disconnected")
	}
	if fake.ConnectCount != 2 {
		t.Errorf("launch did not connect enough")
	}
}
