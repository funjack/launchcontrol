package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// +buid ignore

var versionTmpl = `package main

// generated with gen-version.go
var version = "{{.}}"
`

// Generate version.go based on git tag. Requires git.
func main() {
	cmd := exec.Command("git", "describe", "--tags", "--long")
	gitOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("version.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t := template.Must(template.New("version").Parse(versionTmpl))
	t.Execute(f, strings.TrimSpace(string(gitOutput)))
}
