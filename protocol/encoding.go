package protocol

import (
	"encoding/json"
	"time"
)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ta *TimedAction) UnmarshalJSON(in []byte) error {
	var c struct {
		At  int64
		Pos int
		Spd int
	}
	err := json.Unmarshal(in, &c)
	if err != nil {
		return err
	}
	ta.Position = c.Pos
	ta.Speed = c.Spd
	ta.Time = time.Duration(c.At) * time.Millisecond
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ta TimedAction) MarshalJSON() ([]byte, error) {
	c := struct {
		At  int64 `json:"at"`
		Pos int   `json:"pos"`
		Spd int   `json:"spd"`
	}{
		At:  ta.Time.Nanoseconds() / 1e6,
		Pos: ta.Position,
		Spd: ta.Speed,
	}
	return json.Marshal(&c)
}
