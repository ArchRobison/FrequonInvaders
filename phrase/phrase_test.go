package phrase

import (
	"fmt"
	"testing"
)

func assert(expr bool, msg string, t *testing.T) {
	if !expr {
		t.Fail()
		panic(msg)
	}
}

// Requires visual inspection of stdout
func TestGenerate(t *testing.T) {
	for i := int('1'); i <= int('8'); i++ {
		fmt.Printf("%s\n", Generate(rune(i)))
	}
}
