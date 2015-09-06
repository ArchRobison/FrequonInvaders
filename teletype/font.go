package teletype

import (
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

type charRow uint32

type charMask [charHeight]charRow

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
				y := c>>4*42 + i
				if pm.Pixel(int32(x), int32(y))&0xFF00 < 0x8000 {
					word |= 1 << uint(j)
				}
			}
			font[c][i] = word
		}
	}
	return font
}

var (
	teletypeColor []nimble.Pixel
	teletypeFont  []charMask
)

func fontColor() []nimble.Pixel {
	c := make([]nimble.Pixel, 32)
	for i := range c {
		green := (math32.Exp(-float32(i)*.2) + 0.5) * (2. / 3.)
		red := (math32.Exp(-float32(i) * .4)) * 0.5
		c[i] = nimble.RGB(red, green, 0)
	}
	return c
}
