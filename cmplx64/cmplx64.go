package cmplx64

import (
	"github.com/ArchRobison/FrequonInvaders/math32"
)

func Abs(z complex64) float32 {
	return math32.Hypot(real(z), imag(z))
}

func Phase(z complex64) float32 {
	return math32.Atan2(imag(z), real(z))
}

func Rect(r, θ float32) complex64 {
	y, x := math32.Sincos(θ)
	return complex(r*x, r*y)
}
