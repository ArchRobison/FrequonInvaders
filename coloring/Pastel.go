package coloring

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

func rgbOfHue(hue int32, nHue int32) (r, g, b float32) {
	switch hue {
	case 0:
		// Self
		r = 1
		g = 1
		b = 1
	case 1:
		// Gray
		r = 0.5
		g = 0.5
		b = 0.5
	default:
		r, g, b = phaseColor(float32(hue-2) * 2 * π / float32(nHue))
		if hue&1 != 0 {
			r *= 0.5
			g *= 0.5
			b *= 0.5
		}
		r = (r + 1) * 0.5
		g = (g + 1) * 0.5
		b = (b + 1) * 0.5
	}
	return
}

// PastelPallette generates a pallette of pastel colors.
// The pallette is returned as a PixMap with one row for each hue.
func PastelPallette(nHue, nShade int32) (pm nimble.PixMap) {
	pm = nimble.MakePixMap(nShade, nHue, make([]nimble.Pixel, nHue*nShade), nShade)
	for h := int32(0); h < nHue; h++ {
		r, g, b := rgbOfHue(h, nHue)
		scale := 1 / float32(nShade)
		row := pm.Row(h)
		for j := int32(0); j < nShade; j++ {
			f := float32(nShade-j) * scale
			row[j] = nimble.RGB(r*f, g*f, b*f)
		}
	}
	return
}
