package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/funjack/golaunch"
	"github.com/funjack/launchcontrol/protocol/kiiroo"
)

var (
	file     = flag.String("file", "input.json", "input file (json/txt)")
	videocmd = flag.String("cmd", "", "command to run, eg video player")
	noact    = flag.Bool("noact", false, "simulate")
)

type Launch interface {
	Move(position, speed int)
}
type FakeLaunch string

func (f FakeLaunch) Move(position, speed int) {
	log.Printf("%s: Position=%d, Speed=%d", f, position, speed)
}

func main() {
	flag.Parse()

	// Create Launch by connecting to or faking one
	var l Launch
	if *noact {
		l = FakeLaunch("FakeLaunch")
	} else {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*30)
		la := golaunch.NewLaunch()
		la.HandleDisconnect(func() {
			os.Exit(0)
		})
		err := la.Connect(ctx)
		cancel()
		if err != nil {
			log.Fatal(err)
		}

		l = Launch(la)
	}

	// Create a Kiiroo script player
	player := kiiroo.NewScriptPlayer()

	// Load script from file into player
	f, err := os.Open(filepath.Clean(*file))
	if err != nil {
		log.Fatalf("opening subtitle file, err: %s\n", err)
	}
	if strings.HasSuffix(*file, ".json") {
		var v struct {
			Session struct {
				Subtitle struct {
					Name        string
					Description string
					Text        string
				}
			}
		}
		jsonParser := json.NewDecoder(f)
		if err = jsonParser.Decode(&v); err != nil {
			log.Fatalf("parsing file, err: %s\n", err)
		}
		buf := bytes.NewBufferString(v.Session.Subtitle.Text)
		err = player.Load(buf)
	} else {
		err = player.Load(f)
	}
	if err != nil {
		log.Fatal("error loading script from json")
	}

	// (Optionally) run specified command
	if *videocmd != "" {
		cmd := exec.Command("sh", "-c", *videocmd)
		defer func() {
			log.Printf("Waiting for command to finish...")
			cmd.Wait()
		}()

		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run script
	log.Println("Script started")
	for a := range player.Play() {
		l.Move(a.Position, a.Speed)
	}
	log.Println("Script ended")
}
