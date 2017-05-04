package protocol

import (
	"reflect"
	"testing"
	"time"
)

var jsonTestData = []struct {
	JSON        string
	TimedAction TimedAction
}{
	{
		JSON: `{"at":100,"pos":50,"spd":20}`,
		TimedAction: TimedAction{
			Time: time.Millisecond * 100,
			Action: Action{
				Position: 50,
				Speed:    20,
			},
		},
	},
	{
		JSON: `{"at":150,"pos":100,"spd":50}`,
		TimedAction: TimedAction{
			Time: time.Millisecond * 150,
			Action: Action{
				Position: 100,
				Speed:    50,
			},
		},
	},
}

func TestTimedActionUnmarshalJSON(t *testing.T) {
	for i, c := range jsonTestData {
		var ta TimedAction
		err := (&ta).UnmarshalJSON([]byte(c.JSON))
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}
		if !reflect.DeepEqual(ta, c.TimedAction) {
			t.Errorf("case %d: not equal: want %+v, got %+v",
				i, c.TimedAction, ta)
		}
	}

}

func TestTimedActionMarshalJSON(t *testing.T) {
	for i, c := range jsonTestData {
		out, err := c.TimedAction.MarshalJSON()
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}
		var ta TimedAction
		err = (&ta).UnmarshalJSON(out)
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}
		if !reflect.DeepEqual(ta, c.TimedAction) {
			t.Errorf("case %d: not equal: want %+v, got %+v",
				i, c.TimedAction, ta)
		}
	}
}
