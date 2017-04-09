package kiiroo

import (
	"bytes"
	"testing"
	"time"
)

func TestPlay(t *testing.T) {
	b := bytes.NewBufferString("{1.00:1,1.50:4,1.51:3,1.52:4,1.66:1,1.84:2,1.85:3,1.90:4,1.95:1,2.00:2,2.20:4,2.45:2}")
	k := NewScriptPlayer()
	err := k.Load(b)
	if err != nil {
		t.Error(err)
	}

	var (
		lastPosition  int
		lastEventTime time.Duration
	)

	starttime := time.Now()
	for a := range k.Play() {
		eventtime := time.Now().Sub(starttime)
		t.Logf("Action: %s: %d,%d", eventtime, a.Position, a.Speed)
		if a.Position == lastPosition {
			t.Error("received the same position in a row")
		}

		if lastEventTime > 0 && (eventtime-lastEventTime) < (time.Millisecond*150) {
			t.Errorf("time between events not big enough: %s", eventtime-lastEventTime)
		}

		lastPosition = a.Position
		lastEventTime = eventtime
	}
}
