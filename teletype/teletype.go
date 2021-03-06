package teletype

import (
	"github.com/ArchRobison/Gophetica/nimble"
	"math"
	"unicode"
)

// Init initializes the teletype package.
func Init(fontFilename string) {
	teletypeFont = loadFont(fontFilename)
	teletypeColor = fontColor()
	Reset()
}

// lastLine returns pointer to last line.
func lastLine() *[]byte {
	return &teletypeDisplay[len(teletypeDisplay)-1]
}

// State of teletype display.
var (
	teletypeDisplay = [][]byte{{}}
	displayCursor   bool
)

// Draw the teletype contents on the given PixMap.
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

// Reset the teletype state.
func Reset() {
	teletypeDisplay = [][]byte{{}}
}

// Print one character on the teletype.
// A '\n' goes to the next line.
func PrintChar(c rune) {

	if c == '\n' {
		teletypeDisplay = append(teletypeDisplay, []byte{})
	} else {
		r := lastLine()
		*r = append(*r, byte(c))
	}
}

// Print one character, forcing it to upper case if it's lower case.
func PrintCharUpper(c rune) {
	PrintChar(unicode.ToUpper(c))
}

// Print a string on the teletype
func Print(text string) {
	for _, c := range text {
		PrintChar(c)
	}
}

// PrintUpper prints a string on the teletype in upper case.
func PrintUpper(text string) {
	for _, c := range text {
		PrintCharUpper(c)
	}
}

// Move the cursor backwards one position if possible.
func Backup() {
	r := lastLine()
	if n := len(*r); n > 0 {
		*r = (*r)[:n-1]
	}
}

// CursorLine returns a string representation of last line.
func CursorLine() string {
	return string(*lastLine())
}

// DisplayCursor controls whether the cursor is displayed.
func DisplayCursor(display bool) {
	displayCursor = display
}
