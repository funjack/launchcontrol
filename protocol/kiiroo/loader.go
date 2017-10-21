package kiiroo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/funjack/launchcontrol/protocol"
)

// Load returns a Player with the Kiiroo subtitle/script loaded.
func Load(r io.Reader) (protocol.Player, error) {
	p := NewScriptPlayer()
	err := p.Load(r)
	if err == nil {
		log.Printf("Kiiroo stats: %d actions", len(p.Script))
	}
	return p, err
}

// LoadText loads a VRP txt file and returns a script player.
func LoadText(r io.Reader) (protocol.Player, error) {
	var inKiirooBlock bool
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "[Kiiroo]" {
			inKiirooBlock = true
			continue
		} else if len(line) > 0 && line[0] == '[' {
			inKiirooBlock = false
			continue
		}

		if inKiirooBlock && strings.HasPrefix(line, "onyx=") {
			line = strings.TrimPrefix(line, "onyx=")
			line = strings.Replace(line, ",", ":", -1)
			line = strings.Replace(line, ";", ",", -1)
			return Load(bytes.NewBufferString("{" + line + "}"))
		}
	}
	return nil, ErrEventFormat
}

// LoadJSON loads the FlMe JSON and returns a script player.
func LoadJSON(r io.Reader) (protocol.Player, error) {
	var format struct {
		Text string `json:"text"`
	}
	d := json.NewDecoder(r)
	err := d.Decode(&format)
	if err != nil {
		return nil, err
	}
	return Load(bytes.NewBufferString(format.Text))
}
