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
	// Threshold is the minimum amount of time between actions.
	Threshold = 100 * time.Millisecond
)

// Script is the Funscript container type holding Launch data.
type Script struct {
	// Version of Launchscript
	Version string
	// Inverted causes up and down movement to be flipped.
	Inverted bool `json:"inverted,omitempty"`
	// Range is the percentage of a full stroke to use.
	Range Range `json:"range,omitempty"`
	// Actions are the timed moves.
	Actions []Action
}

// Action is a move at a specific time.
type Action struct {
	// At time in milliseconds the action should fire.
	At int64
	// Pos is the place in percent to move to.
	Pos int
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
	Count              uint64 // Amount of actions generated.
	DistanceTotal      uint64 // Total distance that will be traveled.
	SpeedTotal         uint64 // Accumulation of all commands speed param.
	SpeedOverrideTotal int    // Number of times speed had to be limited.
	SpeedOverrideFast  int    // Number of times generated speed was too fast.
	SpeedOverrideSlow  int    // Number of times generated speed was too slow.
	Delayed            int    // Number of actions that will be thresholded.
}

// String returns a formatted.
func (s Stats) String() string {
	var fastPct, slowPct float64
	if s.SpeedOverrideTotal > 0 {
		fastPct = float64(s.SpeedOverrideFast) /
			float64(s.SpeedOverrideTotal) * 100
		slowPct = float64(s.SpeedOverrideSlow) /
			float64(s.SpeedOverrideTotal) * 100
	}
	avgSpeed := s.SpeedTotal / s.Count
	return fmt.Sprintf("actions=%d (avgspeed=%d%%), delayed=%d, "+
		"speedoverrides=%d (fast=%.2f%%,slow=%.2f%%)",
		s.Count, avgSpeed, s.Delayed, s.SpeedOverrideTotal,
		fastPct, slowPct)
}

// TooFastInc increments the SpeedOveride values
func (s *Stats) TooFastInc() {
	s.SpeedOverrideFast++
	s.SpeedOverrideTotal++
}

// TooSlowInc increments the SpeedOveride values
func (s *Stats) TooSlowInc() {
	s.SpeedOverrideSlow++
	s.SpeedOverrideTotal++
}

// TimedActions creates timed Launch actions from the Scripts timed positions.
// The minspd/maxspd arguments are Launch speed limits in percent. The
// minpos/maxpos specify the position limits in percent.
// The seconds return value are statistics on the script generation.
func (fs Script) TimedActions(minspd, maxspd, minpos, maxpos int) (s protocol.TimedActions, stat Stats) {
	if minspd < SpeedLimitMin {
		minspd = SpeedLimitMin
	}
	if maxspd > SpeedLimitMax {
		maxspd = SpeedLimitMax
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
			stat.TooFastInc()
		} else if speed < minspd {
			speed = minspd
			stat.TooSlowInc()
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
