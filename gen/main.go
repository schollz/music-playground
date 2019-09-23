// This program generates frequencies.go. It can be invoked by running
// go generate ./theory/note/
package main

import (
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// sanitizeVarName replaces forbidden characters for variable names, e.g.
// A#4 => Asharp4; C-1 => C_1
func sanitizeVarName(name string) string {
	return strings.Replace(
		strings.Replace(name, "#", "sharp", 1),
		"-1", "_1", 1)
}

func main() {
	f, err := os.Create("./frequencies.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	type note struct {
		VarName   string
		Name      string
		KeyNumber int
		Freq      float64
	}

	notes := []note{}

	var names = []string{
		"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	for i := 0; i <= 127; i++ {
		octave := int(i/12) - 1
		name := names[i%12] + strconv.Itoa(octave)

		// MIDI key 69 is used for A4, 440 Hz in standard tuning
		// If n = number of semitones between the note and A4
		// then pitch = 440 * 2^(n/12)
		distance := float64(i) - 69
		freq := 440 * math.Pow(2, (distance/12))

		notes = append(notes, note{
			VarName:   sanitizeVarName(name),
			Name:      name,
			KeyNumber: i,
			Freq:      freq,
		})
	}

	tpl.Execute(f, struct {
		Notes []note
	}{
		Notes: notes,
	})
}

var tpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package note

const (
{{- range .Notes}}
	{{ printf "%-8v pitchValue = %v" .VarName .KeyNumber }}
{{- end}}
)

var pitchValues = map[pitchValue]struct {
	name      string
	frequency float64
}{
{{- range .Notes}}
	{{ printf "%-9v {%q, %v}," (printf "%v:" .VarName) .Name .Freq }}
{{- end}}
}
`))
