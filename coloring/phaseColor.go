package coloring

import (
	"github.com/ArchRobison/FrequonInvaders/math32"
)

const π = math32.Pi

func phaseColor(angle float32) (r, g, b float32) {
	// j = index of hextant
	if angle < 0 {
		angle += 2 * π
	}
	j := int(3 / π * angle)
	θ := angle
	if j == 6 {
		j = 0
		θ = 0
	}

	// Red
	switch j {
	case 0, 1:
		r = 1
	case 2:
		r = math32.Cos(1.5*θ - π)
	case 3:
		r = 0
	default:
		r = math32.Cos(1.5*π - .75*θ)
	}

	// Green
	switch j {
	case 0, 1, 2:
		g = math32.Cos(π/2 - 0.5*θ)
	case 3:
		g = math32.Cos(1.5*θ - 3*π/2)
	default:
		g = 0
	}

	// Blue
	switch j {
	case 0, 1, 2:
		b = 0
	case 3:
		b = math32.Cos(1.5*θ - 2*π)
	default:
		b = math32.Cos(.75*θ - π)
	}
	return
}
