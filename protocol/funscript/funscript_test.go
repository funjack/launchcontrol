package funscript

import (
	"reflect"
	"testing"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

type TATestCase struct {
	Name         string
	Script       Script
	TimedActions protocol.TimedActions
	Stats        Stats
}

var TATests = []TATestCase{
	{
		Name:   "Empty script",
		Script: Script{},
		TimedActions: protocol.TimedActions{
			{
				Time: 0,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    SpeedLimitMin,
				},
			},
		},
		Stats: Stats{},
	},
	{
		Name: "Slow overrides",
		Script: Script{
			Version:  "1.0",
			Inverted: false,
			Range:    90,
			Actions: []Action{
				{
					At:  1000,
					Pos: 0,
				},
				{
					At:  2000,
					Pos: 100,
				},
				{
					At:  3000,
					Pos: 0,
				},
			},
		},
		TimedActions: protocol.TimedActions{
			{
				Time: 0,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    SpeedLimitMin,
				},
			},
			{
				Time: 1000 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMax,
					Speed:    SpeedLimitMin,
				},
			},
			{
				Time: 2000 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    SpeedLimitMin,
				},
			},
		},
		Stats: Stats{
			Count:             2,
			DistanceTotal:     (PositionMax - PositionMin) * 2,
			SpeedTotal:        SpeedLimitMin * 2,
			SpeedOverrideFast: 0,
			SpeedOverrideSlow: 2,
			Delayed:           0,
		},
	},
	{
		Name: "Fast overrides",
		Script: Script{
			Version:  "1.0",
			Inverted: false,
			Range:    90,
			Actions: []Action{
				{
					At:  1000,
					Pos: 0,
				},
				{
					At:  1150,
					Pos: 100,
				},
				{
					At:  1300,
					Pos: 0,
				},
			},
		},
		TimedActions: protocol.TimedActions{
			{
				Time: 0,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    SpeedLimitMin,
				},
			},
			{
				Time: 1000 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMax,
					Speed:    SpeedLimitMax,
				},
			},
			{
				Time: 1150 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    SpeedLimitMax,
				},
			},
		},
		Stats: Stats{
			Count:             2,
			DistanceTotal:     (PositionMax - PositionMin) * 2,
			SpeedTotal:        SpeedLimitMax * 2,
			SpeedOverrideFast: 2,
			SpeedOverrideSlow: 0,
			Delayed:           0,
		},
	},
	{
		Name: "Inverted",
		Script: Script{
			Version:  "1.0",
			Inverted: true,
			Range:    90,
			Actions: []Action{
				{
					At:  600,
					Pos: 0,
				},
				{
					At:  1200,
					Pos: 100,
				},
				{
					At:  1800,
					Pos: 0,
				},
			},
		},
		TimedActions: protocol.TimedActions{
			{
				Time: 0,
				Action: protocol.Action{
					Position: PositionMax,
					Speed:    SpeedLimitMin,
				},
			},
			{
				Time: 600 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMin,
					Speed:    30,
				},
			},
			{
				Time: 1200 * time.Millisecond,
				Action: protocol.Action{
					Position: PositionMax,
					Speed:    30,
				},
			},
		},
		Stats: Stats{
			Count:             2,
			DistanceTotal:     (PositionMax - PositionMin) * 2,
			SpeedTotal:        60,
			SpeedOverrideFast: 0,
			SpeedOverrideSlow: 0,
			Delayed:           0,
		},
	},
}

func TestScriptTimedActions(t *testing.T) {
	arguments := [][]int{
		{SpeedLimitMin, SpeedLimitMax, 0, 100},
		{0, 1000, 0, 1000},
	}
	for _, c := range TATests {
		for _, args := range arguments {
			ta, stats := c.Script.TimedActions(args[0], args[1], args[2], args[3])
			if !reflect.DeepEqual(ta, c.TimedActions) {
				t.Errorf("case %s (%v): TimedActions don't match, want %+v, got %+v",
					c.Name, args, c.TimedActions, ta)
			}
			if !reflect.DeepEqual(stats, c.Stats) {
				t.Errorf("case %s (%v): Stats don't match, want %+v, got %+v",
					c.Name, args, c.Stats, stats)
			}
		}
	}
}

func TestRangePosition(t *testing.T) {
	var r Range
	if r.Position(100) != 100 {
		t.Errorf("ranged position failed, 100")
	}
	r = 80
	if r.Position(100) != 80 {
		t.Errorf("ranged position failed, 80")
	}
}

func TestStatsString(t *testing.T) {
	want := "actions=0 (avgspeed=0%), delayed=0, speedoverrides=0 (fast=0.00%,slow=0.00%)"
	var stat Stats
	got := stat.String()
	if got != want {
		t.Errorf("output does not match: want %q, got %q", want, got)
	}
	want = "actions=10 (avgspeed=35%), delayed=1, speedoverrides=4 (fast=25.00%,slow=75.00%)"
	stat = Stats{
		Count:             10,
		DistanceTotal:     500,
		SpeedTotal:        350,
		SpeedOverrideFast: 1,
		SpeedOverrideSlow: 3,
		Delayed:           1,
	}
	got = stat.String()
	if got != want {
		t.Errorf("output does not match: want %q, got %q", want, got)
	}
}
