// Generic version of DFT kernels.
//
// These functions are not for production use, but instead serve as references
// for what the assembly language kernels (e.g. dft_amd64.s) should do.
// They also help test the test routines that test the kernels.
//
// If you need to target a platform without assembly language support,
// you are probably better off using the slow pure Go version of func Draw
// that does not require the kernels.

package fourier

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

// accumulateToFeetSlow is a pure Go reference implementation of accumulateToFeet.
func accumulateToFeetSlow(z *[2]cvec, u *[2]u13, feet []foot) {
	for i := range feet {
		f := &feet[i]
		for j := 0; j < 2; j++ {
			a := &(*z)[j]
			for k := 0; k < vlen; k++ {
				w1 := (*u)[j].u1
				ar := a.re[k]
				ai := a.im[k]
				f.a[k] += ar
				f.b[k] += ai
				f.ac[k] += ar * real(w1)
				f.bc[k] += ai * real(w1)
				f.ad[k] += ar * imag(w1)
				f.bd[k] += ai * imag(w1)
				t := complex(ar, ai) * (*u)[j].u3
				a.re[k] = real(t)
				a.im[k] = imag(t)
			}
		}
	}
}

// feetToPixelsSlow is a pure Go reference implementation of feetToPixels.
func feetToPixelSlow(feet []foot, clut *colorLookupTable, row []nimble.Pixel) {
	for i := range feet {
		f := &feet[i]
		for k := 0; k < vlen; k++ {
			row[(3*i+0)*vlen+k] = clut[int(f.bc[k]-f.ad[k])][int(f.ac[k]+f.bd[k])]
			row[(3*i+1)*vlen+k] = clut[int(f.b[k])][int(f.a[k])]
			row[(3*i+2)*vlen+k] = clut[int(f.bc[k]+f.ad[k])][int(f.ac[k]-f.bd[k])]
			const offset float32 = clutCenter + 0.5
			f.a[k] = offset
			f.b[k] = offset
			f.ac[k] = offset
			f.bc[k] = offset
			f.ad[k] = 0
			f.bd[k] = 0
		}
	}

}

// rotateSlow is a pure Go reference implementation of rotate.
func rotateSlow(w []cvec, v []complex64) {
	for i := range w {
		for k := 0; k < vlen; k++ {
			t := complex(w[i].re[k], w[i].im[k]) * v[i]
			w[i].re[k] = real(t)
			w[i].im[k] = imag(t)
		}
	}
}
