package kiiroo

import (
	"bytes"
	"io"
	"sort"

	"github.com/funjack/launchcontrol/protocol"
)

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
	sort.Sort(es)

	k.Script = k.alg.Actions(es)
	return nil
}
