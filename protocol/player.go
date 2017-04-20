package protocol

import (
	"errors"
	"sync"
	"time"
)

// ErrTimeout is the error returned when a requested operation could not be
// performed in time.
var ErrTimeout = errors.New("operation timed out")

type command int

const (
	cmdStop   = iota // stop playback
	cmdPause         // pause playback
	cmdResume        // resume playback from paused position
	cmdSkip          // skip/jump to position

	commandTimeout = time.Second
)

// control struct is the internal control structure to control the playback
// loop routine.
type control struct {
	Command  command
	Position time.Duration
}

// TimedActionsPlayer can playback an array of TimeActions. It can be used by
// protocols that can pre-calculate TimeActions.
//
// All of the SkippableScriptPlayer methods are implemented except for
// ScriptLoader. Protocols only need to implement ScriptLoader themselves and
// set the Script field with their result.
type TimedActionsPlayer struct {
	// Script that the player will use.
	Script []TimedAction

	wg   sync.WaitGroup
	ctrl chan control
}

// NewTimedActionsPlayer returns a new TimedActionsPlayer.
func NewTimedActionsPlayer() *TimedActionsPlayer {
	return &TimedActionsPlayer{
		ctrl: make(chan control),
	}
}

// Play will start executing the loaded script from the start.
func (ta *TimedActionsPlayer) Play() <-chan Action {
	// Only play one script at a time
	ta.wg.Wait()
	ta.wg.Add(1)
	out := make(chan Action)
	go ta.playbackLoop(out, ta.ctrl)
	return out
}

// sendCommand to the playbackLoop with a timeout.
func (ta *TimedActionsPlayer) sendCommand(c control) error {
	select {
	case ta.ctrl <- c:
		return nil
	case <-time.After(commandTimeout):
		return ErrTimeout
	}
}

// Stop stops playback and resets player.
func (ta *TimedActionsPlayer) Stop() error {
	return ta.sendCommand(control{
		Command: cmdStop,
	})
}

// Pause will halt playback at the current position.
func (ta *TimedActionsPlayer) Pause() error {
	return ta.sendCommand(control{
		Command: cmdPause,
	})
}

// Resume will continue playback from the paused location.
func (ta *TimedActionsPlayer) Resume() error {
	return ta.sendCommand(control{
		Command: cmdResume,
	})
}

// Skip will jump to a specific position.
func (ta *TimedActionsPlayer) Skip(p time.Duration) error {
	return ta.sendCommand(control{
		Command:  cmdSkip,
		Position: p,
	})
}

// playbackLoop will play the loaded script to out and can be controlled using
// ctrl.
func (ta *TimedActionsPlayer) playbackLoop(out chan<- Action, ctrl <-chan control) {
	defer func() {
		ta.wg.Done()
		close(out)
	}()

	var (
		cursor        int           // event position in script
		startTime     = time.Now()  // time playback started/resumed
		startPosition time.Duration // timecode where playback started
		paused        bool
	)

	for cursor < len(ta.Script) {
		a := ta.Script[cursor]
		if a.Time < startPosition {
			cursor++
			continue
		}

		var nextEventTime <-chan time.Time
		if !paused {
			nextEventTime = time.After(
				a.Time - calcPosition(startTime, startPosition))
		}

		select {
		case cmd := <-ctrl:
			switch cmd.Command {
			case cmdStop:
				return
			case cmdPause:
				if !paused {
					paused = true
					startPosition = calcPosition(
						startTime,
						startPosition,
					)
				}
			case cmdResume:
				if paused {
					paused = false
					startTime = time.Now()
					continue
				}
			case cmdSkip:
				startTime = time.Now()
				startPosition = cmd.Position
				cursor = 0
				continue
			}
		case <-nextEventTime:
			if !paused {
				out <- Action{
					Position: a.Position,
					Speed:    a.Speed,
				}
				cursor++
			}
		}

	}
}

// calcPosition will return the current timecode in the script based on start
// time and starting position.
func calcPosition(startTime time.Time, startPosition time.Duration) time.Duration {
	return time.Now().Add(startPosition).Sub(startTime)
}
