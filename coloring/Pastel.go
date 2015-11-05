package coloring

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

func rgbOfHue(hue PastelHue, nHue int) (r, g, b float32) {
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

var (
	rowSize int            // Length of a row in pastel
	pastel  []nimble.Pixel // Linearized matrix with one row per hue
)

type PastelHue int8

// InitPastel initializes the pastsel pallett for m hues with n degrees of fadedness.
func InitPastels(nHue, nShade int) {
	rowSize = nShade
	pastel = make([]nimble.Pixel, nHue*nShade)
	for h := 0; h < nHue; h++ {
		r, g, b := rgbOfHue(PastelHue(h), nHue)
		scale := 1 / float32(nShade)
		for j := 0; j < nShade; j++ {
			f := float32(nShade-j) * scale
			pastel[h*nShade+j] = nimble.RGB(r*f, g*f, b*f)
		}
	}
}

// Pastel returns a pastel for the given hue and fadedness j.
func Pastel(i PastelHue, j int) nimble.Pixel {
	return pastel[int(i)*rowSize+j]
}
