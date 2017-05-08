package raw

import (
	"bytes"
	"testing"
)

var input = `
[
  {
    "at": 100,
    "pos": 50,
    "spd": 30
  },
  {
    "at": 150,
    "pos": 70,
    "spd": 50
  },
  {
    "at": 200,
    "pos": 80,
    "spd": 20
  }
]
`

func TestLoad(t *testing.T) {
	buf := bytes.NewBufferString(input)
	_, err := Load(buf)
	if err != nil {
		t.Errorf("error loading script: %v", err)
	}
}
