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

// ScriptPlayer is an interface that has the basic functions to load and play a
// script. It is supposed to be used to abstract different protocols for the
// player.
type ScriptPlayer interface {
	// Load a script from the provided reader.
	Load(io.Reader) error
	// Start playback of the loaded script the reader channel should be
	// attached to a device.
	Play() <-chan Action
}

// PausableScriptPlayer is a ScriptPlayer that can be paused and resumed.
type PausableScriptPlayer interface {
	ScriptPlayer

	// Pause playback.
	Pause()
	// Resume playback from the current position.
	Resume()
}

// SkippableScriptPlayer is a PausableScriptPlayer that can jump to a time
// location.
type SkippableScriptPlayer interface {
	PausableScriptPlayer

	// Skip (jump) to the specified position/timecode.
	Skip(position time.Duration)
}

// Mover interface provides a device that can move to a position in percent
// with a specific speed.
type Mover interface {
	Move(position, speed int)
}
