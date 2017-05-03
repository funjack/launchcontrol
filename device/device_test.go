package device

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

var testScript = []protocol.TimedAction{
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
	p.Script = testScript
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

func TestTrace(t *testing.T) {
	fake := &fakeLaunch{}
	lm := NewLaunchManager(fake)
	p := protocol.NewTimedActionsPlayer()
	p.Script = testScript
	lm.SetScriptPlayer(p)

	trace := lm.Trace()
	if err := lm.Play(); err != nil {
		t.Error(err)
	}

	breakAtAction := 2
	breakTimer := time.After(testScript[breakAtAction-1].Time +
		time.Millisecond*25)
	traced := make([]protocol.Action, 0, len(testScript))
STOP:
	for {
		select {
		case <-breakTimer:
			break STOP
		case a := <-trace:
			traced = append(traced, a)
		}
	}
	// Make sure at least one more action is written to the unread channel
	time.Sleep(testScript[len(testScript)-1].Time - testScript[breakAtAction].Time)

	if breakAtAction != len(traced) {
		t.Errorf("did not trace all actions: want %d, got %d",
			breakAtAction, len(traced))
	}
}

func TestStop(t *testing.T) {
	fake := &fakeLaunch{}
	lm := NewLaunchManager(fake)
	if err := lm.Stop(); err != nil {
		t.Errorf("stop on empty player did return an error")
	}
	p := protocol.NewTimedActionsPlayer()
	p.Script = testScript
	lm.SetScriptPlayer(p)

	if err := lm.Play(); err != nil {
		t.Error(err)
	}

	stopAtAction := 2
	time.Sleep(testScript[stopAtAction-1].Time + time.Millisecond*25)
	if err := lm.Stop(); err != nil {
		t.Error(err)
	}
	fake.Lock()
	defer fake.Unlock()
	if fake.MoveCount > stopAtAction {
		t.Error("manager did not stop playing")
	}
}

func TestDump(t *testing.T) {
	fake := &fakeLaunch{}
	lm := NewLaunchManager(fake)

	if _, err := lm.Dump(); err != ErrNotSupported {
		t.Errorf("dump on empty player did not return error")
	}

	p := protocol.NewTimedActionsPlayer()
	p.Script = testScript
	lm.SetScriptPlayer(p)

	dump, err := lm.Dump()
	if err != nil {
		t.Error(err)
	}

	if len(testScript) != len(dump) {
		t.Errorf("dump is not complete: want %d, got %d",
			len(testScript), len(dump))
	}

}
