// Functions for displaying the color lookup table.

package fourier

import "github.com/ArchRobison/NimbleDraw"

func CLUTSize() int32 {
	return clutSize
}

func DrawCLUT(pm nimble.PixMap) {
	for y := int32(0); y < clutSize; y++ {
		for x := int32(0); x < clutSize; x++ {
			pm.SetPixel(x, y, clut[y][x])
		}
	}
}
