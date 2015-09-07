package teletype

import (
	"github.com/ArchRobison/Gophetica/nimble"
	"math"
)

const (
	charWidth      = 24
	charHeight     = 32
	textLineHeight = 40
	textTopMargin  = 8
	textLeftMargin = 8
)

func Init(fontFilename string) {
	teletypeFont = loadFont(fontFilename)
	teletypeColor = fontColor()
	Reset()
}

// Return pointer to last line
func lastLine() *[]byte {
	return &teletypeDisplay[len(teletypeDisplay)-1]
}

// State of teletype display
var (
	teletypeDisplay [][]byte
	displayCursor   bool
)

func Draw(pm nimble.PixMap) {
	var r *[]byte = nil
	if displayCursor && math.Mod(nimble.Now(), 1) >= 0.5 {
		r = lastLine()
	}
	if r != nil {
		*r = append(*r, 0)
	}
	draw(pm, teletypeDisplay[:])
	if r != nil {
		*r = (*r)[:len(*r)-1]
	}
}

// Reset the teletype state
func Reset() {
	teletypeDisplay = [][]byte{{}}
}

// Print one character on the teletype.
func PrintChar(c rune) {
	if c == '\n' {
		teletypeDisplay = append(teletypeDisplay, []byte{})
	} else {
		r := lastLine()
		*r = append(*r, byte(c))
	}
}

// Print a string on the teletype
func Print(text string) {
	for _, c := range text {
		PrintChar(c)
	}
}

// Move the cursor backwards one position if possible.
func Backup() {
	r := lastLine()
	if n := len(*r); n > 0 {
		*r = (*r)[:n-1]
	}
}

// Return string representation of last line.
func CursorLine() string {
	return string(*lastLine())
}
