package math32

import (
	"testing"
)

func assert(expr bool, msg string, t *testing.T) {
	if !expr {
		t.Fail()
		panic(msg)
	}
}

func TestMath32(t *testing.T) {
	val := []float32{-2.25, -1.0, -2.25, -0.0, 0.0, 1.0, 2.25}
	for _, x := range val {
		w := Abs(x)
		assert(w == Max(x, -x), "Abs/Max", t)
		assert(-w == Min(x, -x), "Abs/Min", t)
		if x >= 0 {
			z := Sqrt(x)
			assert(z*z == x, "Sqrt", t)
		}
		for _, y := range val {
			a := Min(x, y)
			b := Max(x, y)
			assert(a <= b, "Min/Max", t)
			assert(x == y || a < b, "Min/Max", t)
			assert(a == x || a == y, "Min", t)
			assert(b == x || b == y, "Max", t)
			θ := Atan2(y, x)
			r := Hypot(x, y)
			v, u := Sincos(θ)
			assert(Sin(θ) == v, "SinCos/Sin", t)
			assert(Cos(θ) == u, "SinCos/Cos", t)
			assert(Abs(u*r-x) < 0.000001, "Atan2/Hypot/Cos", t)
			assert(Abs(v*r-y) < 0.000001, "Atan2/Hypot/Sin", t)
		}
	}
}
