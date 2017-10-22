package kiiroo

import (
	"bytes"
	"testing"
)

func TestLoadText(t *testing.T) {
	var inputText = []string{`[Player]
zoom=0

[VideoInfo]
name=Test
version=2

[Kiiroo]
onyx=1.00,4;2.50,1;8.12,3
`,
		`[Player]` + "\r\r" + `
zoom=0` + "\r\r" + `
` + "\r\r" + `
[VideoInfo]` + "\r\r" + `
name=Test` + "\r\r" + `
version=2` + "\r\r" + `
` + "\r\r" + `
[Kiiroo]` + "\r\r" + `
onyx=1.00,4;2.50,1;8.12,3` + "\r\r" + `
`,
	}

	for i, c := range inputText {
		p, err := LoadText(bytes.NewBufferString(c))
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}
		if _, ok := p.(*ScriptPlayer); !ok {
			t.Errorf("case %d: did not return a kiiroo player", i)
		}
	}
}

func TestLoadJSON(t *testing.T) {
	var inputJson = []string{`{
	"text": "{1.00:4,2.50:1}"
}`,
		`{
	"text": " {1.00:4,2.50:1}"
}`,
		`{"subs":{"text":" {1.00:4,2.50:1}"}}`,
	}
	for i, c := range inputJson {
		p, err := LoadJSON(bytes.NewBufferString(c))
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}
		sp, ok := p.(*ScriptPlayer)
		if !ok {
			t.Errorf("case %d: did not return a kiiroo player", i)
		}
		actions, err := sp.Dump()
		if err != nil {
			t.Errorf("case %d: could not dump script", i)
		}
		if len(actions) < 2 {
			t.Errorf("case %d: not enough actions generated", i)
		}
	}
}
