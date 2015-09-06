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

func getCursorChar() byte {
	if math.Mod(nimble.Now(), 1) >= 0.5 {
		return 0 // Cursor character
	} else {
		return 0x20 // ASCII space
	}
}

func setChar(c byte) {
	r := &teletypeDisplay[len(teletypeDisplay)-1]
	if cursorCol >= len(*r) {
		*r = append(*r, byte(c))
	} else {
		(*r)[cursorCol] = c
	}
}

// State of teletype display
var (
	teletypeDisplay [][]byte
	cursorCol       int
	displayCursor   bool
)

func Draw(pm nimble.PixMap) {
	if displayCursor {
		setChar(getCursorChar())
	}
	draw(pm, teletypeDisplay[:])
}

// Reset the teletype state
func Reset() {
	teletypeDisplay = [][]byte{{}}
	cursorCol = 0
}

// Append one character and advance the cursor
func AppendChar(c byte) {
	setChar(c)
	cursorCol++
}

// Append a string
func Append(text string) {
	r := &teletypeDisplay[len(teletypeDisplay)-1]
	*r = append((*r)[:cursorCol], text...)
	cursorCol = len(*r)
}

// Advance to next line
func Newline() {
	teletypeDisplay = append(teletypeDisplay, []byte{})
	cursorCol = 0
}

// Move the cursor backwards one position if possible.
func Backup() {
	if cursorCol > 0 {
		setChar(' ')
		cursorCol--
	}
}

// Return string representation of the line that cursor is on.
func CursorLine() string {
	return string(teletypeDisplay[len(teletypeDisplay)-1][:cursorCol])
}
