package coloring

import "github.com/ArchRobison/Gophetica/math32"

// SchemeBits represents a color scheme for the fourier view.
type SchemeBits byte

const (
	RealBit SchemeBits = 1 << iota
	ImagBit
	MagnitudeBit
	PhaseBit
	RedBit
	GreenBit
	BlueBit
	ColorBits      = RedBit | GreenBit | BlueBit                 // Color bits
	CoordinateBits = RealBit | ImagBit | MagnitudeBit | PhaseBit // Coordinate system bits
	AllBits        = ColorBits | CoordinateBits
)

func (scheme SchemeBits) Color(x, y float32) (r, g, b float32) {
	if scheme&RealBit == 0 {
		x = 0
	}
	if scheme&ImagBit == 0 {
		y = 0
	}
	θ := math32.Atan2(-y, x)
	d := float32(math32.Hypot(x, y))
	if d > 1 || scheme&MagnitudeBit == 0 {
		d = 1.0
	}
	if scheme&PhaseBit == 0 {
		r, g, b = d, d, d
	} else {
		r, g, b = phaseColor(θ)
		r *= d
		g *= d
		b *= d
	}
	if scheme&RedBit == 0 {
		r = 0
	}
	if scheme&GreenBit == 0 {
		g = 0
	}
	if scheme&BlueBit == 0 {
		b = 0
	}
	return
}
