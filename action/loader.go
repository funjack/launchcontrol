package action

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/funjack/launchcontrol/protocol"
	"github.com/funjack/launchcontrol/protocol/kiiroo"
)

// Loaders contains all the registered ScriptLoaders.
var Loaders = []Loader{
	{
		Loader: kiiroo.NewScriptPlayer(),
		ContentTypes: []string{
			"x-text/kiiroo",
			"text/plain",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
		},
	},
}

// ErrUnsupported is returned when the script can't be loaded by any
// scriptplayer.
var ErrUnsupported = errors.New("unsupported script")

// Loader wraps a scriptloader with it's supported mediatypes.
type Loader struct {
	Loader       protocol.ScriptLoader
	ContentTypes []string
}

// IsSupported checks if the loader can handle specified content type.
func (l Loader) IsSupported(contentType string) bool {
	for _, c := range l.ContentTypes {
		if strings.ToUpper(c) == strings.ToUpper(contentType) {
			return true
		}
	}
	return false
}

// LoadScript tries to load specified script with all ScriptPlayers and returns
// the first one that's succesfull.
// Loaders that are tried can be filtered by specifying the content type.
func LoadScript(r io.Reader, contentType string) (protocol.ScriptPlayer, error) {
	supportedLoaders := make([]protocol.ScriptLoader, 0, len(Loaders))
	for _, s := range Loaders {
		if contentType == "" || s.IsSupported(contentType) {
			supportedLoaders = append(supportedLoaders, s.Loader)
		}
	}
	// Just pass the reader if there is only one supported loader.
	if len(supportedLoaders) == 1 {
		return load(supportedLoaders[0], r)
	}
	// Make a copy of the readers contents to be used multiple times.
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	for _, loader := range supportedLoaders {
		if sp, err := load(loader, bytes.NewBuffer(data)); err == nil {
			return sp, nil
		}
	}
	return nil, ErrUnsupported
}

// load will try to load the content of r with scriptloader l and return it's
// player.
func load(l protocol.ScriptLoader, r io.Reader) (protocol.ScriptPlayer, error) {
	err := l.Load(r)
	if err != nil {
		return nil, err
	}

	if sp, ok := l.(protocol.ScriptPlayer); ok {
		return sp, nil
	}
	return nil, ErrUnsupported
}
