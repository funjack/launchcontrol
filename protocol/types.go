// Package protocol provides ways to playback movement.
package protocol

import (
	"io"
	"time"
)

// Action is a command that can be send to a device.
type Action struct {
	Position int
	Speed    int
}

// TimedAction wraps Action together with a timestamp.
type TimedAction struct {
	Action
	Time time.Duration
}

// Loader is the interface that wraps the Load method.
type Loader interface {
	// Load a script from the provided reader.
	Load(r io.Reader) (p Player, err error)
}

// LoaderFunc type is an adapter to allow the use of ordinary functions as
// script loader.
type LoaderFunc func(io.Reader) (Player, error)

// Load calls f(r)
func (f LoaderFunc) Load(r io.Reader) (Player, error) {
	return f(r)
}

// PositionLimiter wraps the limitposition method.
type PositionLimiter interface {
	LimitPosition(lowest, highest int)
}

// SpeedLimiter wraps the limitposition method.
type SpeedLimiter interface {
	LimitSpeed(slowest, fastest int)
}

// LatencyCalibrator wraps the Latency method.
type LatencyCalibrator interface {
	Latency(t time.Duration)
}

// Player is an interface that has the basic functions to play a script.
type Player interface {
	// Start playback of the loaded script the reader channel should be
	// attached to a device.
	Play() <-chan Action

	// Stop stops playback and resets the player.
	Stop() error
}

// Pausable is a interface that defines the pause and resume actions.
type Pausable interface {
	// Pause playback.
	Pause() error
	// Resume playback from the current position.
	Resume() error
}

// Skippable is a interface that wraps the skip method.
type Skippable interface {
	// Skip (jump) to the specified position/timecode.
	Skip(position time.Duration) error
}

// Mover interface provides a device that can move to a position in percent
// with a specific speed.
type Mover interface {
	Move(position, speed int)
}
