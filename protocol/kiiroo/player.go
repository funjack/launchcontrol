package kiiroo

import (
	"bytes"
	"io"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

// ScriptPlayer can load and play Kiiroo scripts/subtitles.
type ScriptPlayer struct {
	alg    Algorithm
	script []TimedAction
}

// NewScriptPlayer returns a new ScriptPlayer using the default algorithm.
func NewScriptPlayer() *ScriptPlayer {
	return &ScriptPlayer{
		alg: DefaultAlgorithm{},
	}
}

// Load reads Kiiroo subtitle/script format.
func (k *ScriptPlayer) Load(r io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	var es Events
	err := es.UnmarshalText(buf.Bytes())
	if err != nil {
		return err
	}
	k.script = k.alg.Actions(es)
	return nil
}

// Play wil start executing the loaded subtitles/script.
func (k *ScriptPlayer) Play() <-chan protocol.Action {
	c := make(chan protocol.Action)
	startTime := time.Now()
	go func() {
		for _, a := range k.script {
			<-time.After(a.Time - time.Now().Sub(startTime))
			c <- protocol.Action{
				Position: a.Position,
				Speed:    a.Speed,
			}
		}
		close(c)
	}()
	return c
}
