package raw

import (
	"encoding/json"
	"io"

	"github.com/funjack/launchcontrol/protocol"
)

// Load returns a player with the waw script loaded.
func Load(r io.Reader) (protocol.Player, error) {
	p := protocol.NewTimedActionsPlayer()
	d := json.NewDecoder(r)
	err := d.Decode(&p.Script)
	return p, err
}
