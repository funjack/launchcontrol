package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/funjack/golaunch"
	"github.com/funjack/launchcontrol/control"
	"github.com/funjack/launchcontrol/device"
)

// Update license.go
//go:generate govendor license -o licenses.go -template tools/gen-license.template

// Update version.go containing the "version" variable
//go:generate go run tools/gen-version.go

var (
	listen   = flag.String("listen", "127.0.0.1:6969", "listen address")
	buttplug = flag.String("buttplug", "", "buttplug.io websocket server address (eg ws://localhost:12345/buttplug)")
	ca       = flag.String("ca", "", "certificate authority in PEM format")
	insecure = flag.Bool("insecure", false, "skip certificate verification")
	noact    = flag.Bool("noact", false, "simulate launch on console")
	lics     = flag.Bool("licenses", false, "show licenses")
	ver      = flag.Bool("version", false, "show version")
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s (%s)", r.URL.Path, r.Method)
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}
	if *lics {
		fmt.Println(licenses)
		os.Exit(0)
	}

	log.Println("Launchcontrol: Get ready for the Launch")

	var l golaunch.Launch
	if *noact {
		l = &launchMock{}
	} else if *buttplug != "" {
		tlscfg, err := createTLSConfig(*ca, *insecure)
		if err != nil {
			log.Fatalf("error creating tls config: %v", err)
		}
		ctx := context.Background()
		l = golaunch.NewButtplugLaunch(ctx, *buttplug, "Launchcontrol", tlscfg)
	} else {
		l = golaunch.NewLaunch()
		defer l.Disconnect()
	}

	lm := device.NewLaunchManager(l)
	c := control.NewController(lm)

	http.Handle("/v1/play", logger(http.HandlerFunc(c.PlayHandler)))
	http.Handle("/v1/stop", logger(http.HandlerFunc(c.StopHandler)))
	http.Handle("/v1/pause", logger(http.HandlerFunc(c.PauseHandler)))
	http.Handle("/v1/resume", logger(http.HandlerFunc(c.ResumeHandler)))
	http.Handle("/v1/skip", logger(http.HandlerFunc(c.SkipHandler)))
	http.Handle("/v1/dump", logger(http.HandlerFunc(c.DumpHandler)))
	http.Handle("/v1/socket", logger(http.HandlerFunc(c.WebsocketHandler)))
	http.Handle("/", logger(http.FileServer(assetFS())))

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

// createTLSConfig creates a configuration trusting certs signed by the ca or
// skip all tls verification.
func createTLSConfig(ca string, insecure bool) (*tls.Config, error) {
	tlscfg := new(tls.Config)
	pool := x509.NewCertPool()
	if ca != "" {
		cacert, err := loadPEMFile(ca)
		if err != nil {
			return tlscfg, err
		}
		pool.AddCert(cacert)
		tlscfg.RootCAs = pool
	}
	if insecure {
		tlscfg.InsecureSkipVerify = true
	}
	return tlscfg, nil
}

// loadCaFile will load the certificate from specified PEM file.
func loadPEMFile(file string) (cert *x509.Certificate, err error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	p, _ := pem.Decode(f)
	if p == nil {
		return cert, errors.New("no pem data found")
	}
	cert, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
