// Test for kernels

package fourier

import (
	"fmt"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/cmplx"
	"testing"
	"time"
)

const runTiming = true

var test *testing.T

func (z *cvec) put(i int, val complex128) {
	z.re[i] = float32(real(val))
	z.im[i] = float32(imag(val))
}

func (z *cvec) get(i int) complex128 {
	return complex(float64(z.re[i]), float64(z.im[i]))
}

func put(z []cvec, i int, val complex128) {
	z[i/vlen].put(i%vlen, val)
}

func get(z []cvec, i int) complex128 {
	return z[i/vlen].get(i % vlen)
}

func toe(feet []foot, i int) complex64 {
	f := &feet[i/pixelsPerFoot]
	k := i % vlen
	switch i / vlen % 3 {
	case 0:
		return complex(f.ac[k]+f.bd[k], f.bc[k]-f.ad[k])
	case 1:
		return complex(f.a[k], f.b[k])
	case 2:
		return complex(f.ac[k]-f.bd[k], f.bc[k]+f.ad[k])
	}
	panic("negative i?")

}

func testAccumulateToFeet(which string, walk func(*[2]cvec, *[2]u13, []foot)) {
	φ := []float64{0.2, 0.3}
	r := []float64{1.25, 0.75}
	// fmt.Printf("ω %v %v\n", cmplx.Abs(ω), cmplx.Phase(ω))
	var a [2]cvec
	var u [2]u13
	for i := 0; i < 2; i++ {
		for k := 0; k < vlen; k++ {
			a[i].put(k, cmplx.Rect(r[i], φ[i]*float64(k)))
		}
		ω := cmplx.Rect(1, φ[i]*vlen)
		u[i].u1 = complex64(ω)
		u[i].u3 = complex64(ω * ω * ω)
	}

	// n = number of feet
	const n = 1920 / pixelsPerFoot
	var f [n]foot
	walk(&a, &u, f[:])

	for k := 0; k < pixelsPerFoot*n; k++ {
		expected := cmplx.Rect(r[0], φ[0]*float64(k-vlen)) + cmplx.Rect(r[1], φ[1]*float64(k-vlen))
		actual := complex128(toe(f[:], k))
		dist := cmplx.Abs(actual - expected)
		if dist > 1.1E-5 {
			test.Errorf("[%v] actual=%v expected=%v err=%v\n", k, actual, expected, dist)
		}
	}
	if runTiming {
		t0 := time.Now()
		trials := 100000
		for i := 0; i < trials; i++ {
			walk(&a, &u, f[:])
		}
		t1 := time.Now()
		fmt.Printf("%v: %.2f Gflop/sec\n", which, float64(n*vlen*16*2*trials)/1E9/t1.Sub(t0).Seconds())
	}
}

func testRotate(which string, rot func([]cvec, []complex64)) {
	const n = 5
	a := make([]cvec, n)
	v := make([]complex64, n)
	for i := 0; i < n; i++ {
		v[i] = complex64(cmplx.Rect(1, float64(i)*0.1))
		for j := vlen * i; j < vlen*i+vlen; j++ {
			put(a, j, cmplx.Rect(1, float64(j)/2))
		}
	}
	b := make([]cvec, n)
	copy(b, a)
	rot(a, v)
	for i := 0; i < n; i++ {
		for j := vlen * i; j < vlen*i+vlen; j++ {
			actual := get(a, j)
			expected := complex128(v[i]) * get(b, j)
			if cmplx.Abs(actual-expected) > 1E-7 {
				test.Errorf("%v %v %v %v %v %v\n", which, i, get(b, j), v[i], actual, actual-expected)
			}
		}
	}
}

func checkY(p nimble.Pixel, z float32) {
	if z < 0 || z >= clutSize {
		panic("z too large")
	}
	y := p >> 16
	if y != nimble.Pixel(z) {
		test.Errorf("Y error: y=%x z=%v\n", y, z)
	}
}

func checkX(p nimble.Pixel, z float32) {
	if z < 0 || z >= clutSize {
		panic("z too large")
	}
	x := p & 0xFFFF
	if x != nimble.Pixel(z) {
		test.Errorf("X error: x=%x z=%v\n", x, z)
	}
}

func checkEqual(actual float32, expected float32) {
	if actual != expected {
		test.Errorf("got %v, expected %v\n", actual, expected)
	}
}

func testFeetToPixels(which string, toPixels func([]foot, *colorLookupTable, []nimble.Pixel)) {
	n := 1920 / pixelsPerFoot // Number of feet

	// Initialize feet
	feet := make([]foot, n)
	for i := 0; i < n; i++ {
		f := &feet[i]
		for k := 0; k < vlen; k++ {
			f.a[k] = float32((1 + 2*i ^ k) % 128)
			f.b[k] = float32((2 + 3*i ^ k) % 128)
			f.bd[k] = float32((5*i ^ k) % 32)
			f.ad[k] = float32((7*i ^ k) % 32)
			f.ac[k] = float32((11*i^k)%64 + 32)
			f.bc[k] = float32((13*i^k)%64 + 32)
		}
	}

	// Initialize lookup table
	clut := colorLookupTable{}
	for i := range clut {
		for j := range clut[i] {
			clut[i][j] = nimble.Pixel(i<<16 | j)
		}
	}

	row := make([]nimble.Pixel, pixelsPerFoot*n)

	origFeet := make([]foot, len(feet))
	copy(origFeet, feet)
	toPixels(feet, &clut, row)

	// Check result
	for i := 0; i < n; i++ {
		f := &origFeet[i]
		g := &feet[i]
		for k := 0; k < vlen; k++ {
			var p nimble.Pixel
			p = row[pixelsPerFoot*i+0+k]
			checkX(p, f.ac[k]+f.bd[k])
			checkY(p, f.bc[k]-f.ad[k])
			p = row[pixelsPerFoot*i+vlen+k]
			checkX(p, f.a[k])
			checkY(p, f.b[k])
			p = row[pixelsPerFoot*i+2*vlen+k]
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
	if runTiming {
		t0 := time.Now()
		trials := 100000
		for i := 0; i < trials; i++ {
			toPixels(feet, &clut, row)
		}
		t1 := time.Now()
		fmt.Printf("%v: %.2f Gpixel/sec\n", which, float64(pixelsPerFoot*n*trials)/1E9/t1.Sub(t0).Seconds())
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
