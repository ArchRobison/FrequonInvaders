package phrase

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	g := MakeGrammar([]string{
		"1: %A %N in %L",
		"A: red",
		"A: yellow",
		"L: house",
		"N: fish",
		"L: room",
		"Z: Forty = $ minus 2",
		"N: cat",
		"A: blue",
	})
	words := [][]string{
		{"red", "blue", "yellow"},
		{"cat", "fish"},
		{"in"},
		{"room", "house"}}
	var tally [4][3]int
	// Try generating phrases and check that only expected ones show up.
	for k := 0; k < 1000; k++ {
		s := strings.Split(g.Generate(rune('1')), " ")
		for i, w := range s {
			found := false
			for j, u := range words[i] {
				if u == w {
					tally[i][j]++
					found = true
				}
			}
			if !found {
				panic(w)
			}
		}
	}
	// Check that distribution of sentences is plausible.
	for i := range words {
		k := len(words[i])
		for j := range words[i] {
			if tally[i][j] < 900/k {
				panic("bad distribution")
			}
		}
	}
	if g.GenerateWithNumber('Z', 42) != "Forty = 42 minus 2" {
		t.Fail()
	}
}
