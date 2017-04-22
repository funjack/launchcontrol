package action

import (
	"log"
	"mime"
	"net/http"
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
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			mediaType = ""
		}
		k, err := LoadScript(r.Body, mediaType)
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