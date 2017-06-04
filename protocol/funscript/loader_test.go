package funscript

import (
	"bytes"
	"io"
	"testing"
)

type LoaderTestCase struct {
	Input  io.Reader
	Limits []int // spdmin,spdmax,posmin,posmax
	Error  error
}

var script = `{
	"version":"1.0",
	"inverted":false,
	"range":90,
	"actions":[
		{"at":100,"pos":0},
		{"at":600,"pos":100},
		{"at":1200,"pos":0}
	]
}`

var LoaderTests = []LoaderTestCase{
	{
		Input:  bytes.NewBufferString(script),
		Limits: []int{20, 80, 5, 95},
		Error:  nil,
	},
	{
		Input:  bytes.NewBufferString(script),
		Limits: []int{0, 1000, 0, 1000},
		Error:  nil,
	},
	{
		Input:  bytes.NewBufferString(script),
		Limits: []int{1000, 0, 1000, 0},
		Error:  nil,
	},
}

func TestLoaderLoad(t *testing.T) {
	for i, c := range LoaderTests {
		var l Loader
		l.LimitSpeed(c.Limits[0], c.Limits[1])
		l.LimitPosition(c.Limits[2], c.Limits[3])
		_, err := l.Load(c.Input)
		if err != c.Error {
			t.Errorf("case %d: errors did not match, want %v, got %v",
				i, c.Error, err)
		}
	}
}

func TestLoaderString(t *testing.T) {
	var l Loader
	want := "funscript loader (speeds:0-0) (positions:0-0)"
	got := l.String()
	if got != want {
		t.Errorf("strings do not match, want %q, got %q", want, got)
	}
	l.LimitSpeed(20, 80)
	l.LimitPosition(5, 95)
	want = "funscript loader (speeds:20-80) (positions:5-95)"
	got = l.String()
	if got != want {
		t.Errorf("strings do not match, want %q, got %q", want, got)
	}
}
