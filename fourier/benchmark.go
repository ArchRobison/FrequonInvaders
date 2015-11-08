package fourier

import (
	"time"
	"unsafe"
)

// runBenchmark times the DFT kernel for the given number of iterations.
// For sake of comparison with classic Frequon Invaders, the kernel
// does only the accumulateToFeet and rotate steps, not the feetToPixels
// step that does a few more sums and differences to compute the
// actual Fourier transform.
//
// The current DFT kernel does only 16 floating-point operations (flops)
// per foot, whereas the old one in "classic" did 24 flops.
func runBenchmark(iterations int) (secs, flops float64) {
	const m = 4
	var (
		u [m]u13
		v [m]complex64
		w [m]cvec
	)
	for i := 0; i < m; i++ {
		ω := float32(i) * 0.01
		for k := 0; k < 4; k++ {
			c := euler(float32(k+4) * ω)
			w[i].re[k] = real(c)
			w[i].im[k] = imag(c)
		}
		u[i].u1 = euler(vlen * ω)
		u[i].u3 = euler(pixelsPerFoot * ω)
		v[i] = euler(ω)
	}
	// 140 = number for feet used for game when played in full HDTV resolution
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
				(*[2]u13)(unsafe.Pointer(&u[i])), feet)
		}
		rotate(w[:], v[:])
	}
	t1 := time.Now()
	secs = t1.Sub(t0).Seconds()
	// The 16*4*p is the floating-point operations for accumulateToFeet
	// The +8 is for "rotate"

	flops = float64((16*4*p + 8) * m * iterations)
	return
}

// Benchmark runs the CPU speed benchmark for Frequon Invaders.
// Returns floating-point operations per second.
func Benchmark() float64 {
	for n := 1; ; n *= 2 {
		secs, flops := runBenchmark(n)
		if secs > 0.1 || n >= 1<<20 {
			return flops / secs
		}
		// Time too short.  Double number of iterations and retry.
	}
}
