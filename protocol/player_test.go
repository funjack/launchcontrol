package protocol

import (
	"testing"
	"time"
)

var script = []TimedAction{
	{
		Action{
			Position: 5,
			Speed:    50,
		},
		time.Millisecond * 50,
	},
	{
		Action{
			Position: 50,
			Speed:    40,
		},
		time.Millisecond * 100,
	},
	{
		Action{
			Position: 90,
			Speed:    90,
		},
		time.Millisecond * 150,
	},
	{
		Action{
			Position: 30,
			Speed:    30,
		},
		time.Millisecond * 200,
	},
}

func TestPlay(t *testing.T) {
	p := NewTimedActionsPlayer()
	p.Script = script

	var eventCount int
	starttime := time.Now()
	for a := range p.Play() {
		eventtime := time.Now().Sub(starttime)
		t.Logf("Action: %s: %d,%d", eventtime, a.Position, a.Speed)
		eventCount++
	}
	playTime := time.Now().Sub(starttime)

	if eventCount != len(script) {
		t.Errorf("not all actions were played, want %d, got %d",
			len(script), eventCount)
	}
	want := script[len(script)-1].Time
	if playTime.Nanoseconds()/1e6 != want.Nanoseconds()/1e6 {
		t.Errorf("script was not played back at correct speed")
	}
}

func TestPauseResume(t *testing.T) {
	p := NewTimedActionsPlayer()
	p.Script = script

	pauseTime := time.Millisecond * 100

	go func() {
		<-time.After(time.Millisecond * 100)
		if err := p.Pause(); err != nil {
			t.Error(err)
		}
		<-time.After(pauseTime)
		if err := p.Resume(); err != nil {
			t.Error(err)
		}
	}()

	var eventCount int
	starttime := time.Now()
	for range p.Play() {
		eventCount++
	}
	playTime := time.Now().Sub(starttime)

	if eventCount != len(script) {
		t.Errorf("not all actions were played, want %d, got %d",
			len(script), eventCount)
	}
	want := script[len(script)-1].Time + pauseTime
	if playTime.Nanoseconds()/1e6 != want.Nanoseconds()/1e6 {
		t.Errorf("script was not played back at correct speed")
	}
}

func TestStop(t *testing.T) {
	p := NewTimedActionsPlayer()
	p.Script = script

	stopTime := time.Millisecond * 75

	go func() {
		<-time.After(stopTime)
		if err := p.Stop(); err != nil {
			t.Error(err)
		}
	}()

	starttime := time.Now()
	for range p.Play() {
		// pass
	}
	playTime := time.Now().Sub(starttime)

	if playTime.Nanoseconds()/1e6 != stopTime.Nanoseconds()/1e6 {
		t.Errorf("script was not stopped at the right time, want: %s, got: %s",
			stopTime, playTime)
	}
}

func TestSkip(t *testing.T) {
	cases := []struct {
		Name string
		At   time.Duration
		To   time.Duration
	}{
		{
			Name: "Forward",
			At:   time.Millisecond * 50,
			To:   time.Millisecond * 150,
		},
		{
			Name: "Back",
			At:   time.Millisecond * 50,
			To:   time.Millisecond * 150,
		},
	}

	for _, c := range cases {
		p := NewTimedActionsPlayer()
		p.Script = script

		go func() {
			<-time.After(c.At)
			if err := p.Skip(c.To); err != nil {
				t.Error(err)
			}
		}()

		starttime := time.Now()
		for range p.Play() {
			// pass
		}
		playTime := time.Now().Sub(starttime)
		want := script[len(script)-1].Time - (c.To - c.At)

		if playTime.Nanoseconds()/1e6 != want.Nanoseconds()/1e6 {
			t.Errorf("%s: player did not skip correctly, want: %s, got: %s",
				c.Name, want, playTime)
		}
	}
}
