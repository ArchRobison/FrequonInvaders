package fourier

import (
	"github.com/ArchRobison/Gophetica/cmplx64"
	"time"
	"unsafe"
)

// runBenchmark times the HFT kernel for the given number of iterations.
// For sake of comparison with classic Frequon Invaders, the kernel
// does only the accumulateToFeet and rotate steps, not the feetToPixels
// step that does a few more sums and differences to compute the
// actual Fourier transform.
func runBenchmark(iterations int) (secs, flops float64) {
	const m = 4
	var (
		u [m]w13
		v [m]complex64
		w [m]cvec
	)
	for i := 0; i < m; i++ {
		ω := float32(i) * 0.01
		for k := 0; k < 4; k++ {
			c := cmplx64.Rect(1.0, float32(k+4)*ω)
			w[i].re[k] = real(c)
			w[i].im[k] = imag(c)
		}
		u[i].w1 = cmplx64.Rect(1, 4*ω)
		u[i].w3 = cmplx64.Rect(1, 12*ω)
		v[i] = cmplx64.Rect(1, ω)
	}
	const p = 140
	feet := make([]foot, p)
	var t0 time.Time
	for t := -1; t < iterations; t++ {
		// t==-1 warms up caches
		if t == 0 {
			t0 = time.Now()
		}
		for i := 0; i < m; i += 2 {
			accumulateToFeet(
				(*[2]cvec)(unsafe.Pointer(&w[i])),
				(*[2]w13)(unsafe.Pointer(&u[i])), feet)
			rotate(w[:], v[:])
		}
	}
	t1 := time.Now()
	secs = t1.Sub(t0).Seconds()
	flops = float64(16 * 4 * p * m * iterations)
	return
}

// Return gigaflops
func Benchmark() float64 {
	for n := 1; ; n *= 2 {
		secs, flops := runBenchmark(n)
		if secs > 0.1 || n >= 1<<20 {
			return flops / secs * 1E-9
		}
	}
}
