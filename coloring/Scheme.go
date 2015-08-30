package coloring

import "github.com/ArchRobison/FrequonInvaders/math32"

type Scheme int

const (
	HasReal Scheme = 1 << iota
	HasImag
	HasMagnitude
	HasPhase
	HasRed
	HasGreen
	HasBlue
	HasEverything Scheme = 1<<iota - 1
)

func (scheme Scheme) Color(x, y float32) (r, g, b float32) {
	if scheme&HasReal == 0 {
		x = 0
	}
	if scheme&HasImag == 0 {
		y = 0
	}
	θ := math32.Atan2(-y, x)
	d := float32(math32.Hypot(x, y))
	if d > 1 || scheme&HasMagnitude == 0 {
		d = 1.0
	}
	if scheme&HasPhase == 0 {
		r, g, b = d, d, d
	} else {
		r, g, b = phaseColor(θ)
		r *= d
		g *= d
		b *= d
	}
	if scheme&HasRed == 0 {
		r = 0
	}
	if scheme&HasGreen == 0 {
		g = 0
	}
	if scheme&HasBlue == 0 {
		b = 0
	}
	return
}
