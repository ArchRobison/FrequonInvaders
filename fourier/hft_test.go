package fourier

import (
	"fmt"
	"math/cmplx"
	"testing"
	"time"
)

const runBenchmark = true

var test *testing.T

func feetToSum(feet []foot, sum []cvec) {
	for i := range feet {
		f := &feet[i]
		for k := 0; k < vlen; k++ {
			sum[3*i+0].re[k] = f.ac[k] + f.bd[k]
			sum[3*i+0].im[k] = f.bc[k] - f.ad[k]
			sum[3*i+1].re[k] = f.a[k]
			sum[3*i+1].im[k] = f.b[k]
			sum[3*i+2].re[k] = f.ac[k] - f.bd[k]
			sum[3*i+2].im[k] = f.bc[k] + f.ad[k]
		}
	}
}

func testAccumulateToFeet(which string, walk func(*[2]cvec, *[2]w13, []foot)) {
	φ := []float64{0.2, 0.3}
	r := []float64{1.25, 0.75}
	// fmt.Printf("ω %v %v\n", cmplx.Abs(ω), cmplx.Phase(ω))
	var a [2]cvec
	var u [2]w13
	for i := 0; i < 2; i++ {
		for k := 0; k < 4; k++ {
			a[i].put(k, cmplx.Rect(r[i], φ[i]*float64(k)))
		}
		ω := cmplx.Rect(1, φ[i]*4)
		u[i].w1 = complex64(ω)
		u[i].w3 = complex64(ω * ω * ω)
	}

	// n = number of chickenfeet
	const n = 1920 / 12
	var f [n]foot
	walk(&a, &u, f[:])

	var sum [n * 3]cvec
	for k := 0; k < 12*n; k++ {
		put(sum[:], k, 0)
	}
	feetToSum(f[:], sum[:])

	for k := 0; k < 12*n; k++ {
		expected := cmplx.Rect(r[0], φ[0]*float64(k-4)) + cmplx.Rect(r[1], φ[1]*float64(k-4))
		actual := get(sum[:], k)
		dist := cmplx.Abs(actual - expected)
		if dist > 1.1E-5 {
			test.Errorf("[%v] actual=%v expected=%v err=%v\n", k, actual, expected, dist)
		}
	}
	if runBenchmark {
		t0 := time.Now()
		trials := 100000
		for i := 0; i < trials; i++ {
			walk(&a, &u, f[:])
		}
		t1 := time.Now()
		fmt.Printf("%v: %v Gflop/sec\n", which, float64(n*4*16*2*trials)/1E9/t1.Sub(t0).Seconds())
	}
}

func testRotate(which string, rot func([]cvec, []complex64)) {
	const n = 5
	a := make([]cvec, n)
	v := make([]complex64, n)
	for i := 0; i < n; i++ {
		v[i] = complex64(cmplx.Rect(1, float64(i)*0.1))
		for j := 4 * i; j < 4*i+4; j++ {
			put(a, j, cmplx.Rect(1, float64(j)/2))
		}
	}
	b := make([]cvec, n)
	copy(b, a)
	rot(a, v)
	for i := 0; i < n; i++ {
		for j := 4 * i; j < 4*i+4; j++ {
			actual := get(a, j)
			expected := complex128(v[i]) * get(b, j)
			if cmplx.Abs(actual-expected) > 1E-7 {
				test.Errorf("%v %v %v %v %v %v\n", which, i, get(b, j), v[i], actual, actual-expected)
			}
		}
	}
}

func checkY(p pixel, z float32) {
	if z < 0 || z >= clutSize {
		panic("z too large")
	}
	y := p >> 16
	if y != pixel(z) {
		test.Errorf("Y error: y=%x z=%v\n", y, z)
	}
}

func checkX(p pixel, z float32) {
	if z < 0 || z >= clutSize {
		panic("z too large")
	}
	x := p & 0xFFFF
	if x != pixel(z) {
		test.Errorf("X error: p=%x z=%v\n", x, z)
	}
}

func checkEqual(actual float32, expected float32) {
	if actual != expected {
		test.Errorf("got %v, expected %v\n", actual, expected)
	}
}

func testFeetToPixels(which string, toPixels func([]foot, *[128][128]pixel, []pixel)) {
	n := 1920 / 12 // Number of feet

	// Initialize feet
	feet := make([]foot, n)
	for i := 0; i < n; i++ {
		f := &feet[i]
		for k := 0; k < 4; k++ {
			f.a[k] = float32((1 + 2*i ^ k) % 128)
			f.b[k] = float32((2 + 3*i ^ k) % 128)
			f.bd[k] = float32((5*i ^ k) % 32)
			f.ad[k] = float32((7*i ^ k) % 32)
			f.ac[k] = float32((11*i^k)%64 + 32)
			f.bc[k] = float32((13*i^k)%64 + 32)
		}
	}

	// Initialize lookup table
	clut := [128][128]pixel{}
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			clut[i][j] = pixel(i<<16 | j)
		}
	}

	// Must have 12 pixels per foot
	row := make([]pixel, 12*n)

	origFeet := make([]foot, len(feet))
	copy(origFeet, feet)
	toPixels(feet, &clut, row)

	// Check result
	for i := 0; i < n; i++ {
		f := &origFeet[i]
		g := &feet[i]
		for k := 0; k < 4; k++ {
			var p pixel
			p = row[12*i+0+k]
			checkX(p, f.ac[k]+f.bd[k])
			checkY(p, f.bc[k]-f.ad[k])
			p = row[12*i+4+k]
			checkX(p, f.a[k])
			checkY(p, f.b[k])
			p = row[12*i+8+k]
			checkX(p, f.ac[k]-f.bd[k])
			checkY(p, f.bc[k]+f.ad[k])
			checkEqual(g.a[k], 64.5) 
			checkEqual(g.b[k], 64.5)
			checkEqual(g.ac[k], 64.5)
			checkEqual(g.bc[k], 64.5)
			checkEqual(g.ad[k], 0)
			checkEqual(g.bd[k], 0)
		}
	}
	if runBenchmark {
		t0 := time.Now()
		trials := 100000
		for i := 0; i < trials; i++ {
			toPixels(feet, &clut, row)
		}
		t1 := time.Now()
		fmt.Printf("%v: %v Gpixel/sec\n", which, float64(12*n*trials)/1E9/t1.Sub(t0).Seconds())
	}
}

func TestRotate(c *testing.T) {
    test = c
	testRotate("slow", rotateSlow)
	testRotate("fast", rotate)
}

func TestAccumulateToFeet(c *testing.T) {
    test = c
	testAccumulateToFeet("slow", accumulateToFeetSlow)
	testAccumulateToFeet("fast", accumulateToFeet)
}

func TestFeetToPixels(c *testing.T) {
    test = c
	testFeetToPixels("slow", feetToPixelSlow)
	testFeetToPixels("fast", feetToPixel)
}
