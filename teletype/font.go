package teletype

import (
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

const (
	charWidth      = 24
	charHeight     = 32
	textLineHeight = 40
	textTopMargin  = 24
	textLeftMargin = 24
)

type (
	charRow  uint32
	charMask [charHeight]charRow
)

var (
	teletypeColor []nimble.Pixel
	teletypeFont  []charMask
)

func loadFont(filename string) []charMask {
	pm, err := nimble.ReadPixMap(filename)
	if err != nil {
		panic(err)
	}
	font := make([]charMask, 128)
	for c := range font {
		for i := range font[c] {
			word := charRow(0)
			for j := 0; j < charWidth; j++ {
				x := c&0xF*24 + j
				y := c>>4*42 + i - 4
				if y >= 0 && pm.Pixel(int32(x), int32(y))&0xFF00 < 0x8000 {
					word |= 1 << uint(j)
				}
			}
			font[c][i] = word
		}
	}
	return font
}

func fontColor() []nimble.Pixel {
	c := make([]nimble.Pixel, 32)
	for i := range c {
		green := (math32.Exp(-float32(i)*.2) + 0.5) * (2./3.)
        // Classic version had 0.5 for red
		red := (math32.Exp(-float32(i) * .4)) * 0.4
		c[i] = nimble.RGB(red, green, 0)
	}
	return c
}
