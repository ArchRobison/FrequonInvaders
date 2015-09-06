package teletype

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

func draw(pm nimble.PixMap, text [][]byte) {
	width, height := pm.Size()

	// Clear area
	pm.Fill(nimble.Black)

	// Write lines of text
	for m := range text {
		x := int32(textLeftMargin)
		for j := range text[m] {
			if x >= width {
				break
			}
			kLimit := width - x
			if kLimit > charWidth {
				kLimit = charWidth
			}
			for i, mask := range teletypeFont[text[m][j]] {
				y := int32(textTopMargin + m*textLineHeight + i)
				if y >= height {
					break
				}
				pixelRow := pm.Row(y)[x : x+kLimit]
				colorIndex := 0
				for k := range pixelRow {
					if mask&(1<<uint(k)) != 0 {
						pixelRow[k] = teletypeColor[colorIndex]
						colorIndex++
					} else {
						colorIndex = 0
					}
				}
			}
			x += charWidth
		}
	}
}
