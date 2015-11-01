// Functions for displaying the color lookup table.
// These are not for production use, but sometimes useful for debugging.

package fourier

import "github.com/ArchRobison/Gophetica/nimble"

// CLUTSize returns the size (of one axis) of the square color lookup table.
func CLUTSize() int32 {
	return clutSize
}

// DrawCLUT draws the color lookup table on the given PixMap.
// The PixMap must be at least as big as CLUTSize() x CLUTSize().
func DrawCLUT(pm nimble.PixMap) {
	for y := int32(0); y < clutSize; y++ {
		for x := int32(0); x < clutSize; x++ {
			pm.SetPixel(x, y, clut[y][x])
		}
	}
}
