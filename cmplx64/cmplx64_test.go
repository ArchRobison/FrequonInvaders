package cmplx64

import (
	"math"
	"testing"
)

func assertTiny(a float32, tol float32, msg string, t *testing.T) {
	if a > tol || a < -tol {
		t.Fail()
		panic(msg)
	}
}

const π = math.Pi

func TestCmplx64(t *testing.T) {
	for _, r := range []float32{0.5, 1, 2} {
		for _, θ := range []float32{-3, -π / 2, 0, π / 2, 3} {
			a := Rect(r, θ)
			assertTiny(Abs(a)-r, 1e-6, "Rect/Abs", t)
			assertTiny(Phase(a)-θ, 1e-6, "Rect/Phase", t)
		}
	}
}
