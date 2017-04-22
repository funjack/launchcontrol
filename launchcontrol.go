package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/funjack/golaunch"
	"github.com/funjack/launchcontrol/action"
	"github.com/funjack/launchcontrol/manager"
)

var (
	listen = flag.String("listen", "127.0.0.1:6969", "listen address")
	noact  = flag.Bool("noact", false, "simulate launch on console")
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s (%s)", r.URL.Path, r.Method)
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	log.Println("Launchcontrol: Get ready for the Launch")

	var l golaunch.Launch
	if *noact {
		l = &launchMock{}
	} else {
		l = golaunch.NewLaunch()
		defer l.Disconnect()
	}

	lm := manager.NewLaunchManager(l)
	c := action.NewController(lm)

	http.Handle("/v1/play", logger(http.HandlerFunc(c.PlayHandler)))
	http.Handle("/v1/stop", logger(http.HandlerFunc(c.StopHandler)))
	http.Handle("/v1/pause", logger(http.HandlerFunc(c.PauseHandler)))
	http.Handle("/v1/resume", logger(http.HandlerFunc(c.ResumeHandler)))
	http.Handle("/v1/skip", logger(http.HandlerFunc(c.SkipHandler)))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		log.Println("Shutting down...")
		l.Disconnect()
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()

	log.Printf("Listening on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
