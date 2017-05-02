package action

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/funjack/launchcontrol/manager"
)

// Controller translates http requests into manager actions.
type Controller struct {
	manager *manager.LaunchManager
}

// NewController returns a new controller for the given manager.
func NewController(m *manager.LaunchManager) *Controller {
	return &Controller{
		manager: m,
	}
}

// PlayHandler is a http.Handler to load and play scripts.
func (c *Controller) PlayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pers := parsePlayParams(r.URL.Query())
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			mediaType = ""
		}
		k, err := LoadScript(r.Body, mediaType, pers)
		if err == ErrUnsupported {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		} else if err != nil {
			log.Printf("Error loading script: %s\n", err)
			internalServerError(w)
			return
		}
		if err := c.manager.SetScriptPlayer(k); err != nil {
			handleManagerError(w, err)
			return
		}
	}
	handleManagerError(w, c.manager.Play())
}

// StopHandler is a http.Handler to stop playback.
func (c *Controller) StopHandler(w http.ResponseWriter, r *http.Request) {
	handleManagerError(w, c.manager.Stop())
}

// PauseHandler is a http.Handler to pause playback.
func (c *Controller) PauseHandler(w http.ResponseWriter, r *http.Request) {
	handleManagerError(w, c.manager.Pause())
}

// ResumeHandler is a http.Handler to resume playback.
func (c *Controller) ResumeHandler(w http.ResponseWriter, r *http.Request) {
	handleManagerError(w, c.manager.Resume())
}

// SkipHandler is a http.Handler to jump to a given timecode.
func (c *Controller) SkipHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := time.ParseDuration(r.Form.Get("p"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	handleManagerError(w, c.manager.Skip(p))
}

// DumpHandler is a http.Handler to dump the current script.
func (c *Controller) DumpHandler(w http.ResponseWriter, r *http.Request) {
	script, err := c.manager.Dump()
	if err != nil {
		handleManagerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	e := json.NewEncoder(w)
	err = e.Encode(&script)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
	return
}

// handleManagerError writes a http response based on a manager error.
func handleManagerError(w http.ResponseWriter, err error) {
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	case manager.ErrNotSupported:
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("operation not supported by loaded script type\n"))
	case manager.ErrNotPlaying:
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("operation cannot be executed when not playing\n"))
	default:
		log.Printf("Internal server error, %s\n", err)
		internalServerError(w)
	}
}

// internalServerError returns a status 500 with message to a ResponseWriter.
func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error\n"))
}

// parsePlayParams extracts personalization values from the query params.
func parsePlayParams(q url.Values) Personalization {
	p := NewPersonalization()
	if i, err := strconv.Atoi(q.Get("latency")); err == nil {
		p.Latency = time.Duration(i) * time.Millisecond
	}
	if i, err := strconv.Atoi(q.Get("positionmin")); err == nil {
		p.PositionMin = i
	}
	if i, err := strconv.Atoi(q.Get("positionmax")); err == nil {
		p.PositionMax = i
	}
	if i, err := strconv.Atoi(q.Get("speedmin")); err == nil {
		p.SpeedMin = i
	}
	if i, err := strconv.Atoi(q.Get("speedmax")); err == nil {
		p.SpeedMax = i
	}
	return p
}
