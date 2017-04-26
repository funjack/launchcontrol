package kiiroo

import (
	"bytes"
	"io"

	"github.com/funjack/launchcontrol/protocol"
)

// Load returns a Player with the Kiiroo subtitle/script loaded.
func Load(r io.Reader) (protocol.Player, error) {
	p := NewScriptPlayer()
	err := p.Load(r)
	return p, err
}

// ScriptPlayer can load and play Kiiroo scripts/subtitles.
type ScriptPlayer struct {
	*protocol.TimedActionsPlayer

	alg Algorithm
}

// NewScriptPlayer returns a new ScriptPlayer using the default algorithm.
func NewScriptPlayer() *ScriptPlayer {
	return &ScriptPlayer{
		protocol.NewTimedActionsPlayer(),
		DefaultAlgorithm{},
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

	k.Script = k.alg.Actions(es)
	return nil
}
