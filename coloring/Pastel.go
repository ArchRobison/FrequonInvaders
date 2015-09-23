package coloring

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

func PastelFade(p []nimble.Pixel, k, maxCritter int) {
	var r, g, b float32
	if k == 0 {
		// Self
		r = 1
		g = 1
		b = 1
	} else if k == 1 {
		r = 0.5
		g = 0.5
		b = 0.5
	} else {
		r, g, b = phaseColor(float32(k-2) * 2 * π / float32(maxCritter))
		if k&1 != 0 {
			r *= 0.5
			g *= 0.5
			b *= 0.5
		}
		r = (r + 1) * 0.5
		g = (g + 1) * 0.5
		b = (b + 1) * 0.5
	}
	n := len(p)
	scale := 1 / float32(n)
	for j := range p {
		f := float32(n-j) * scale
		p[j] = nimble.RGB(r*f, g*f, b*f)
	}
}
