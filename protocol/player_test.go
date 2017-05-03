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
		<-time.After(time.Millisecond * 75)
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
		{
			Name: "Soon",
			At:   time.Millisecond * 0,
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

func TestLimits(t *testing.T) {
	type TestCase struct {
		Low, High  int
		Slow, Fast int
	}

	cases := []TestCase{
		{0, 100, 0, 0},
		{10, 90, 0, 0},
		{30, 50, 0, 0},
		{50, 90, 0, 0},
		{60, 80, 0, 0},
		{0, 0, 20, 80},
		{0, 0, 30, 60},
		{0, 0, 50, 90},
		{0, 0, 10, 20},
		{0, 0, 0, 50},
		{0, 100, 20, 80},
		{10, 90, 30, 60},
		{30, 50, 50, 90},
		{50, 90, 10, 20},
		{60, 80, 0, 50},
	}
	for i, c := range cases {
		p := NewTimedActionsPlayer()
		p.Script = script
		if c.Low != 0 && c.High != 0 {
			p.LimitPosition(c.Low, c.High)
		}
		if c.Slow != 0 && c.Fast != 0 {
			p.LimitSpeed(c.Slow, c.Fast)
		}

		lowest, highest := c.Low, c.High
		slowest, fastest := c.Slow, c.Fast

		for a := range p.Play() {
			if a.Position < lowest {
				lowest = a.Position
			} else if a.Position > highest {
				highest = a.Position
			}
			if a.Position < slowest {
				slowest = a.Speed
			} else if a.Position > highest {
				fastest = a.Speed
			}
		}

		if c.Low != 0 && c.High != 0 {
			if lowest < c.Low {
				t.Errorf("case %d: went lower than allowed, %d < %d", i, lowest, c.Low)
			}
			if highest > c.High {
				t.Errorf("case %d: went higher than allowed, %d > %d", i, highest, c.High)
			}
		}

		if c.Slow != 0 && c.Fast != 0 {
			if slowest < c.Slow {
				t.Errorf("case %d: went slower than allowed, %d < %d", i, slowest, c.Slow)
			}
			if fastest > c.Fast {
				t.Errorf("case %d: went higher than allowed, %d > %d", i, fastest, c.Fast)
			}
		}

	}
}
