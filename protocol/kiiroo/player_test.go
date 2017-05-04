package kiiroo

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

// Badish input scenario, contains dups and short timings
var scenario = "{1.00:1,1.50:4,1.51:4,1.51:3,1.52:4,1.66:1,1.84:2,1.85:3,1.90:4,1.95:1,2.00:2,2.20:4,2.45:2}"

func playerwithscenario(scenario string) (protocol.Player, error) {
	b := bytes.NewBufferString(scenario)
	return Load(b)
}

type actionValidator struct {
	LastPostion int
	LastTime    time.Duration
}

// Validate takes a position and time and tests if that is allowed compared to
// previous values validated.
func (a *actionValidator) Validate(p int, t time.Duration) error {
	defer func() {
		a.LastPostion = p
		a.LastTime = t
	}()

	if p == a.LastPostion {
		return fmt.Errorf("received the same position in a row")
	}
	if a.LastTime > 0 && (t-a.LastTime) < (time.Millisecond*150) {
		return fmt.Errorf("time between events not big enough: %s", t-a.LastTime)
	}
	return nil
}

func TestPlay(t *testing.T) {
	k, err := playerwithscenario(scenario)
	if err != nil {
		t.Error(err)
	}

	av := actionValidator{}
	starttime := time.Now()
	for a := range k.Play() {
		eventtime := time.Now().Sub(starttime)
		t.Logf("Action: %s: %d,%d", eventtime, a.Position, a.Speed)
		if err := av.Validate(a.Position, eventtime); err != nil {
			t.Error(err)
		}
	}
}
