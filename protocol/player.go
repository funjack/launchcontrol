package protocol

import (
	"errors"
	"time"
)

var (
	// ErrPlaying is the error returned when an action is requested that is
	// not allowed while playing.
	ErrPlaying = errors.New("action not allowed while playing")
	// ErrStopped is the error returned when an action is requested that is
	// not allowed while stopped.
	ErrStopped = errors.New("action not allowed while stopped")
)

// TimedActionsPlayer can playback an array of TimeActions. It can be used by
// protocols that can pre-calculate TimeActions.
//
// All of the SkippableScriptPlayer methods are implemented except for
// ScriptLoader. Protocols only need to implement ScriptLoader themselves and
// set the Script field with their result.
type TimedActionsPlayer struct {
	Script []TimedAction

	isPlaying      bool
	cancelPlaying  chan bool
	stoppedPlaying chan bool

	startPosition time.Duration
	pausePosition time.Duration

	startTime  time.Time
	actionChan chan Action
}

// NewTimedActionsPlayer returns a new TimedActionsPlayer.
func NewTimedActionsPlayer() *TimedActionsPlayer {
	return &TimedActionsPlayer{
		cancelPlaying:  make(chan bool),
		stoppedPlaying: make(chan bool),
	}
}

// Play will start executing the loaded script from the start.
func (k *TimedActionsPlayer) Play() <-chan Action {
	return k.PlayFrom(0)
}

// PlayFrom will start executing the loaded script from the given position.
func (k *TimedActionsPlayer) PlayFrom(p time.Duration) <-chan Action {
	if k.isPlaying == true {
		return k.actionChan
	}

	k.actionChan = make(chan Action)
	k.startTime = time.Now()
	k.startPosition = p
	go k.startPlaying()
	return k.actionChan
}

// Position will return the current position in the script.
func (k *TimedActionsPlayer) Position() time.Duration {
	if !k.isPlaying {
		return k.pausePosition
	}
	// Now()+startPosition - startTime
	return time.Now().Add(k.startPosition).Sub(k.startTime)
}

// Stop stops playback and resets player.
func (k *TimedActionsPlayer) Stop() error {
	if !k.isPlaying {
		return ErrStopped
	}
	k.cancelPlaying <- true
	<-k.stoppedPlaying
	k.reset()
	return nil
}

// Pause will stop playback at the current position.
func (k *TimedActionsPlayer) Pause() error {
	if !k.isPlaying {
		return ErrStopped
	}
	k.pausePosition = k.Position()
	k.cancelPlaying <- true
	<-k.stoppedPlaying
	return nil
}

// Resume will continue playback from the paused location.
func (k *TimedActionsPlayer) Resume() error {
	if k.isPlaying {
		return ErrPlaying
	}
	k.startPosition = k.pausePosition
	k.startTime = time.Now()
	go k.startPlaying()
	return nil
}

// Skip will jump to a specific position.
func (k *TimedActionsPlayer) Skip(p time.Duration) error {
	if !k.isPlaying {
		return ErrStopped
	}
	k.Pause()
	k.pausePosition = p
	k.Resume()
	return nil
}

// startPlaying is the actual play loop sending out actions called as a
// goroutine.
func (k *TimedActionsPlayer) startPlaying() {
	k.isPlaying = true
	for _, a := range k.Script {
		if a.Time < k.startPosition {
			continue
		}
		select {
		case <-k.cancelPlaying:
			k.isPlaying = false
			k.stoppedPlaying <- true
			return
		case <-time.After(a.Time - k.Position()):
			k.actionChan <- Action{
				Position: a.Position,
				Speed:    a.Speed,
			}
		}
	}
	// end of script reached
	k.reset()
	k.isPlaying = false
}

// reset rewinds the player when playback has finished or is stopped.
func (k *TimedActionsPlayer) reset() {
	close(k.actionChan)
	k.pausePosition = 0 // Be kind rewind
	k.isPlaying = false
	k.startTime = time.Time{}
}
