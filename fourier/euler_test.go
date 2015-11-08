package fourier

import (
	"math"
	"testing"
)

func assertTiny(a float64, tol float64, msg string, t *testing.T) {
	if a > tol || a < -tol {
		t.Fail()
		panic(msg)
	}
}

const π = math.Pi

func TestEuler(t *testing.T) {
	for _, θ := range []float32{-3, -π / 2, 0, π / 2, 3} {
		a := complex128(euler(θ))
		assertTiny(math.Hypot(real(a), imag(a))-1, 1e-6, "euler/hypot", t)
		assertTiny(math.Atan2(imag(a), real(a))-float64(θ), 1e-6, "euler/atan2", t)
	}
}
