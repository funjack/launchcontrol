package device

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/funjack/golaunch"
	"github.com/funjack/launchcontrol/protocol"
)

var (
	// ErrNotSupported is returned when the request operation is not
	// supported by the active scriptplayer.
	ErrNotSupported = errors.New("operation is not supported")
	// ErrNotPlaying is returned when the operation is only supported while
	// a script is playing.
	ErrNotPlaying = errors.New("not playing")
)

// ConnectionTimeout is the default timeout used per Launch connecting attempt.
var ConnectionTimeout = time.Second * 10

// LaunchManager is responsible for connecting and communicating with the
// Launch.
type LaunchManager struct {
	sync.Mutex

	launch      golaunch.Launch
	isConnected bool

	wg     sync.WaitGroup
	player protocol.Player

	tracers map[chan protocol.Action]bool

	playingMux sync.Mutex
	playing    bool
}

// NewLaunchManager creates a new manager for the given Launch.
func NewLaunchManager(l golaunch.Launch) *LaunchManager {
	lm := &LaunchManager{
		launch:  l,
		tracers: make(map[chan protocol.Action]bool),
	}
	lm.launch.HandleDisconnect(func() {
		lm.Lock()
		defer lm.Unlock()
		// TODO implement nice reconnect handling
		lm.isConnected = false
		lm.player.Stop()
		lm.wg.Wait()
	})

	return lm
}

// SetScriptPlayer switches the active ScriptPlayer. Any active script will be
// stopped.
func (m *LaunchManager) SetScriptPlayer(p protocol.Player) error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		if err := m.player.Stop(); err != nil {
			return err
		}
		m.wg.Wait()
	}
	m.player = p
	return nil
}

// connect will check if the manager has an active connection with the Launch
// or else tries to (re)connect.
func (m *LaunchManager) connect() error {
	if m.isConnected {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()
	if err := m.launch.Connect(ctx); err != nil {
		return err
	}

	m.isConnected = true
	return nil
}

// Play will start playback using the loaded scriptplayer to the connected
// Launch.
func (m *LaunchManager) Play() error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		return nil
	}

	if err := m.connect(); err != nil {
		return err
	}

	if m.player != nil {
		m.playingMux.Lock()
		m.playing = true
		m.playingMux.Unlock()
		m.wg.Add(1)
		go func() {
			for a := range m.player.Play() {
				m.launch.Move(a.Position, a.Speed)
				for t := range m.tracers {
					select {
					case t <- a:
					default:
						close(t)
						m.Lock()
						delete(m.tracers, t)
						m.Unlock()
					}
				}
			}
			m.playingMux.Lock()
			m.playing = false
			m.playingMux.Unlock()
			m.wg.Done()
		}()
	}
	return nil
}

// Stop will halt playback and reset the scriptplayer.
func (m *LaunchManager) Stop() error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		return m.player.Stop()
	}
	return nil
}

// Pause will halt playback but keep the current position.
func (m *LaunchManager) Pause() error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		if pp, ok := m.player.(protocol.Pausable); ok {
			return pp.Pause()
		}
		return ErrNotSupported
	}
	return ErrNotPlaying
}

// Resume starts playback from the paused position.
func (m *LaunchManager) Resume() error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		if pp, ok := m.player.(protocol.Pausable); ok {
			return pp.Resume()
		}
		return ErrNotSupported
	}
	return ErrNotPlaying
}

// Skip jumps playback position to the specified time.
func (m *LaunchManager) Skip(p time.Duration) error {
	m.Lock()
	defer m.Unlock()

	if m.isPlaying() {
		if pp, ok := m.player.(protocol.Skippable); ok {
			return pp.Skip(p)
		}
		return ErrNotSupported
	}
	return ErrNotPlaying
}

// Dump will return the full loaded script.
func (m *LaunchManager) Dump() (protocol.TimedActions, error) {
	m.Lock()
	defer m.Unlock()

	if pp, ok := m.player.(protocol.Dumpable); ok {
		return pp.Dump()
	}
	return nil, ErrNotSupported
}

// Trace returns a channel that receives the same actions as are send to the
// Launch.
func (m *LaunchManager) Trace() <-chan protocol.Action {
	t := make(chan protocol.Action, 8)
	m.tracers[t] = true
	return t
}

// isPlaying returns true if the loaded scriptplayer is playing.
func (m *LaunchManager) isPlaying() bool {
	m.playingMux.Lock()
	defer m.playingMux.Unlock()
	return m.playing
}

// WaitUntilStopped will return a empty struct on the done channel when the
// player is stopped.
func (m *LaunchManager) WaitUntilStopped(done chan<- struct{}) {
	go func() {
		m.wg.Wait()
		done <- struct{}{}
	}()
}
