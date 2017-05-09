package control

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/funjack/launchcontrol/protocol"
	"github.com/funjack/launchcontrol/protocol/kiiroo"
	"github.com/funjack/launchcontrol/protocol/raw"
)

// Loaders contains all the registered ScriptLoaders.
var Loaders = []Loader{
	{
		Loader: protocol.LoaderFunc(raw.Load),
		ContentTypes: []string{
			"application/prs.launchcontrol+json",
			"application/json",
		},
	},
	{
		Loader: protocol.LoaderFunc(kiiroo.Load),
		ContentTypes: []string{
			"text/prs.kiiroo",
			"x-text/kiiroo",
			"text/plain",
		},
	},
	{
		Loader: protocol.LoaderFunc(kiiroo.LoadText),
		ContentTypes: []string{
			"text/plain",
		},
	},
	{
		Loader: protocol.LoaderFunc(kiiroo.LoadJSON),
		ContentTypes: []string{
			"application/json",
		},
	},
}

// ErrUnsupported is returned when the script can't be loaded by any
// scriptplayer.
var ErrUnsupported = errors.New("unsupported script")

// Personalization are settings customizing a scripts behaviour
type Personalization struct {
	Latency     time.Duration
	PositionMin int // Lowest position
	PositionMax int // Highest position
	SpeedMin    int // Slowest speed to move at
	SpeedMax    int // Fastest speed to move at
}

// NewPersonalization return a Personalization with the default values.
func NewPersonalization() Personalization {
	return Personalization{
		Latency:     0,
		PositionMin: 5,
		PositionMax: 95,
		SpeedMin:    20,
		SpeedMax:    80,
	}
}

// Loader wraps a scriptloader with it's supported mediatypes.
type Loader struct {
	Loader       protocol.Loader
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
func LoadScript(r io.Reader, contentType string, p Personalization) (protocol.Player, error) {
	supportedLoaders := make([]protocol.Loader, 0, len(Loaders))
	for _, s := range Loaders {
		if contentType == "" || s.IsSupported(contentType) {
			supportedLoaders = append(supportedLoaders, s.Loader)
		}
	}
	// Just pass the reader if there is only one supported loader.
	if len(supportedLoaders) == 1 {
		return load(supportedLoaders[0], r, p)
	}
	// Make a copy of the readers contents to be used multiple times.
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	for _, loader := range supportedLoaders {
		if sp, err := load(loader, bytes.NewBuffer(data), p); err == nil {
			return sp, nil
		}
	}
	return nil, ErrUnsupported
}

// load will try to load the content of r with scriptloader l and return it's
// player.
func load(l protocol.Loader, r io.Reader, pers Personalization) (protocol.Player, error) {
	personalizeLoader(l, pers)
	p, err := l.Load(r)
	if err != nil {
		return nil, err
	}
	personalizePlayer(p, pers)
	return p, nil
}

// personalizeLoader will apply, if supported, personalized position and speed
// limits to the loader.
func personalizeLoader(l protocol.Loader, pers Personalization) {
	if pl, ok := l.(protocol.PositionLimiter); ok {
		pl.LimitPosition(pers.PositionMin, pers.PositionMax)
	}
	if sl, ok := l.(protocol.SpeedLimiter); ok {
		sl.LimitSpeed(pers.SpeedMin, pers.SpeedMax)
	}
}

// personalizePlayer will apply, if supported, personalized latency, position
// and speed limits to the loader.
func personalizePlayer(p protocol.Player, pers Personalization) {
	if lc, ok := p.(protocol.LatencyCalibrator); ok {
		lc.Latency(pers.Latency)
	}
	if pl, ok := p.(protocol.PositionLimiter); ok {
		pl.LimitPosition(pers.PositionMin, pers.PositionMax)
	}
	if sl, ok := p.(protocol.SpeedLimiter); ok {
		sl.LimitSpeed(pers.SpeedMin, pers.SpeedMax)
	}
}
