package fourier

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

// Size of color-lookup table (CLUT) along either axis.
// It is a power of 2 to speed up indexing.
// If you change this number, you must update lgClutSize in gen_*.jl.
const clutSize = 128

const (
	clutCenter = clutSize / 2   // CLUT index corresponding to 0.
	clutRadius = clutCenter - 1 // CLUT distance from center representing magnitude of 1.
)

// clutCoor returns the real/imag value corresponding to a CLUT index.
func clutCoor(k int) (z float32) {
	const (
		clutScale  = 1.0 / clutRadius
		clutOffset = -clutCenter * clutScale
	)
	return float32(k)*clutScale + clutOffset
}

// colorLookupTable maps the complex plane onto colors
type colorLookupTable [clutSize][clutSize]nimble.Pixel

// theClut is the color lookup table (CLUT)
var clut colorLookupTable

type colorMap interface {
	Color(x, y float32) (r, g, b float32)
}

// currentMap is the current color map used to generate the CLUT.
var currentMap colorMap

// setColoring sets the CLUT to given colorMap.
func setColoring(cm colorMap) {
	if cm == currentMap {
		return
	}
	currentMap = cm

	for i := range clut {
		y := clutCoor(i)
		for j := range clut[i] {
			x := clutCoor(j)
			clut[i][j] = nimble.RGB(cm.Color(x, y))
		}
	}
}
