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

// ScriptLoader is the interface that wraps the Load method.
type ScriptLoader interface {
	// Load a script from the provided reader.
	Load(io.Reader) error
}

// ScriptPlayer is an interface that has the basic functions to play a script.
type ScriptPlayer interface {
	// Start playback of the loaded script the reader channel should be
	// attached to a device.
	Play() <-chan Action

	// Stop stops playback and resets the player.
	Stop() error
}

// PausableScriptPlayer is a ScriptPlayer that can be paused and resumed.
type PausableScriptPlayer interface {
	ScriptPlayer

	// Pause playback.
	Pause() error
	// Resume playback from the current position.
	Resume() error
}

// SkippableScriptPlayer is a PausableScriptPlayer that can jump to a time
// location.
type SkippableScriptPlayer interface {
	PausableScriptPlayer

	// Skip (jump) to the specified position/timecode.
	Skip(position time.Duration) error
}

// Mover interface provides a device that can move to a position in percent
// with a specific speed.
type Mover interface {
	Move(position, speed int)
}
