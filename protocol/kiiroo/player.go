package kiiroo

import (
	"bytes"
	"io"

	"github.com/funjack/launchcontrol/protocol"
)

// scriptPlayer can load and play Kiiroo scripts/subtitles.
type scriptPlayer struct {
	*protocol.TimedActionsPlayer

	alg Algorithm
}

// NewScriptPlayer returns a new ScriptPlayer using the default algorithm.
func NewScriptPlayer() protocol.SkippableScriptPlayer {
	return &scriptPlayer{
		protocol.NewTimedActionsPlayer(),
		DefaultAlgorithm{},
	}
}

// Load reads Kiiroo subtitle/script format.
func (k *scriptPlayer) Load(r io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	var es Events
	err := es.UnmarshalText(buf.Bytes())
	if err != nil {
		return err
	}

	k.Script = k.alg.Actions(es)
	return nil
}
