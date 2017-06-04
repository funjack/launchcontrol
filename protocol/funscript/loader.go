package funscript

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/funjack/launchcontrol/protocol"
)

// Loader loads Funscripts.
type Loader struct {
	speedMin    int
	speedMax    int
	positionMin int
	positionMax int
}

// String returns a formatted string.
func (l Loader) String() string {
	return fmt.Sprintf("funscript loader (speeds:%d-%d) (positions:%d-%d)",
		l.speedMin, l.speedMax, l.positionMin, l.positionMax)
}

// LimitPosition implements the PositionLimiter interface.
func (l *Loader) LimitPosition(lowest, highest int) {
	switch min := lowest; true {
	case min < 0:
		l.positionMin = 0
	case min > 90:
		l.positionMin = 90
	default:
		l.positionMin = min
	}
	switch max := highest; true {
	case max > 100:
		l.positionMax = 100
	case max < 10:
		l.positionMax = 10
	default:
		l.positionMax = max
	}
}

// LimitSpeed implements the SpeedLimiter interface.
func (l *Loader) LimitSpeed(slowest, fastest int) {
	switch min := slowest; true {
	case min < SpeedLimitMin:
		l.speedMin = SpeedLimitMin
	case min > SpeedLimitMax:
		l.speedMin = SpeedLimitMax
	default:
		l.speedMin = min
	}
	switch max := fastest; true {
	case max > SpeedLimitMax:
		l.speedMax = SpeedLimitMax
	case max < SpeedLimitMin:
		l.speedMax = SpeedLimitMin
	default:
		l.speedMax = max
	}
}

// Load returns a player with the Funscript loaded.
func (l Loader) Load(r io.Reader) (protocol.Player, error) {
	p := protocol.NewTimedActionsPlayer()
	var s Script
	d := json.NewDecoder(r)
	err := d.Decode(&s)
	if err != nil {
		return p, err
	}
	log.Printf("Loading Funscript: %s", l)
	var stats Stats
	p.Script, stats = s.TimedActions(
		l.speedMin,
		l.speedMax,
		l.positionMin,
		l.positionMax)
	log.Printf("Funscript stats: %s", stats)
	return p, nil
}
