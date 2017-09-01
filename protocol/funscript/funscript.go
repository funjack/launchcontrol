package funscript

import (
	"fmt"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

const (
	// SpeedLimitMin is the slowest movement command possible. The Launch
	// crashes on very slow speeds.
	SpeedLimitMin = 20
	// SpeedLimitMax is the fasts movement command possible. The Launch
	// makes weird 'clicking' noises when moving at very fast speeds.
	SpeedLimitMax = 80
	// PositionMin is the lowest position possible.
	PositionMin = 5
	// PositionMax is the hight position possible.
	PositionMax = 95
	// Threshold is the minimum amount of time between actions.
	Threshold = 100 * time.Millisecond
)

// Script is the Funscript container type holding Launch data.
type Script struct {
	// Version of Launchscript
	Version string `json:"version"`
	// Inverted causes up and down movement to be flipped.
	Inverted bool `json:"inverted,omitempty"`
	// Range is the percentage of a full stroke to use.
	Range Range `json:"range,omitempty"`
	// Actions are the timed moves.
	Actions []Action `json:"actions"`
}

// Action is a move at a specific time.
type Action struct {
	// At time in milliseconds the action should fire.
	At int64 `json:"at"`
	// Pos is the place in percent to move to.
	Pos int `json:"pos"`
}

// Range in percent.
type Range int

// Position returns the position for p limited within the range.
func (r Range) Position(p int) int {
	if r > 0 {
		return int(float64(p) / 100 * float64(r))
	}
	return p
}

// Stats when generated script.
type Stats struct {
	Count             uint64 // Amount of actions generated.
	DistanceTotal     uint64 // Total distance that will be traveled.
	SpeedTotal        uint64 // Accumulation of all commands speed param.
	SpeedOverrideFast int    // Number of times generated speed was too fast.
	SpeedOverrideSlow int    // Number of times generated speed was too slow.
	Delayed           int    // Number of actions that will be thresholded.
}

// String returns a formatted.
func (s Stats) String() string {
	var fastPct, slowPct float64
	overrideTotal := s.SpeedOverrideFast + s.SpeedOverrideSlow
	if overrideTotal > 0 {
		fastPct = float64(s.SpeedOverrideFast) /
			float64(overrideTotal) * 100
		slowPct = float64(s.SpeedOverrideSlow) /
			float64(overrideTotal) * 100
	}
	var avgSpeed uint64
	if s.Count > 0 {
		avgSpeed = s.SpeedTotal / s.Count
	}
	return fmt.Sprintf("actions=%d (avgspeed=%d%%), delayed=%d, "+
		"speedoverrides=%d (fast=%.2f%%,slow=%.2f%%)",
		s.Count, avgSpeed, s.Delayed, overrideTotal,
		fastPct, slowPct)
}

// TimedActions creates timed Launch actions from the Scripts timed positions.
// The minspd/maxspd arguments are Launch speed limits in percent. The
// minpos/maxpos specify the position limits in percent.
// The second return value are statistics on the script generation.
func (fs Script) TimedActions(minspd, maxspd, minpos, maxpos int) (s protocol.TimedActions, stat Stats) {
	if minspd < SpeedLimitMin {
		minspd = SpeedLimitMin
	}
	if maxspd > SpeedLimitMax {
		maxspd = SpeedLimitMax
	}
	if minpos < PositionMin {
		minpos = PositionMin
	}
	if maxpos > PositionMax {
		maxpos = PositionMax
	}
	r := Range(maxpos - minpos)
	if fs.Range != 0 && r > fs.Range {
		r = fs.Range
	}

	s = make(protocol.TimedActions, 1, len(fs.Actions)+1)
	s[0].Time = 0
	if fs.Inverted {
		s[0].Position = maxpos // Init at top
	} else {

		s[0].Position = minpos // Init at bottom
	}
	s[0].Speed = SpeedLimitMin

	previousPosition := s[0].Position
	previous := Action{
		At:  0,
		Pos: 0,
	}
	for _, e := range fs.Actions {
		if e.Pos == previous.Pos {
			previous = e
			continue
		}
		timediff := time.Duration(e.At-previous.At) * time.Millisecond
		if timediff < Threshold {
			stat.Delayed++
		}
		position := e.Pos
		if fs.Inverted {
			position = 100 - e.Pos
		}
		position = r.Position(position) + minpos
		distance := position - previousPosition
		if distance < 0 {
			distance = -distance
		}
		stat.DistanceTotal = stat.DistanceTotal + uint64(distance)
		speed := Speed(distance, timediff)
		if speed > maxspd {
			speed = maxspd
			stat.SpeedOverrideFast++
		} else if speed < minspd {
			speed = minspd
			stat.SpeedOverrideSlow++
		}
		stat.SpeedTotal = stat.SpeedTotal + uint64(speed)
		ta := protocol.TimedAction{
			Time: time.Duration(previous.At) * time.Millisecond,
			Action: protocol.Action{
				Position: position,
				Speed:    speed,
			},
		}
		s = append(s, ta)
		stat.Count++
		previous = e
		previousPosition = ta.Position
	}
	return
}
