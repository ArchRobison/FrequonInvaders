// Version of Draw that relies on fast assembly-coded kernels.
//
// +build amd64,!slow

package fourier

import (
	"github.com/ArchRobison/Gophetica/cmplx64"
	"github.com/ArchRobison/Gophetica/nimble"
	"unsafe"
)

var (
	feetStorage []foot
	uStorage    []u13
	vStorage    []complex64
	wStorage    []cvec
	tmpStorage  [pixelsPerFoot]nimble.Pixel
)

// Init initializes storage used by Draw.
// It should be called once before calling Draw any number of times.
// The call to Init must specify an upper bound on the width of the PixMap
// and number of harmonics for subsequent calls to Draw.
func Init(widthMax int32, harmonicLenMax int) {
	footLenMax := (widthMax + pixelsPerFoot - 1) / pixelsPerFoot

	feetStorage = make([]foot, footLenMax)
	// Use dummy call to initialize "feet"
	feetToPixel(feetStorage, &clut, make([]nimble.Pixel, footLenMax*pixelsPerFoot))

	// The internal kernels assume an even number of harmonics,
	// so round the value provided by the user up to even.
	m := harmonicLenMax + harmonicLenMax&1
	uStorage = make([]u13, m)
	vStorage = make([]complex64, m)
	wStorage = make([]cvec, m)
}

// Draw draws a Fourier transform on the given PixMap.
// Transform values must lie on the unit circle in the complex plane.
func Draw(pm nimble.PixMap, harmonics []Harmonic, cm colorMap) {
	setColoring(cm)
	n := len(harmonics)
	// m = n rounded up to even
	m := n + n&1
	u := uStorage[:m]
	v := vStorage[:m]
	w := wStorage[:m]
	for i, h := range harmonics {
		for k := 0; k < vlen; k++ {
			c := cmplx64.Rect(h.Amplitude*clutRadius, h.Phase+float32(k+vlen)*h.Ωx)
			w[i].re[k] = real(c)
			w[i].im[k] = imag(c)
		}
		u[i].u1 = cmplx64.Rect(1, vlen*h.Ωx)
		u[i].u3 = cmplx64.Rect(1, pixelsPerFoot*h.Ωx)
		v[i] = cmplx64.Rect(1, h.Ωy)
	}
	if n < m {
		// Zero the extra element.
		w[n] = cvec{}
		u[n] = u13{}
		v[n] = 0
	}
	width, height := pm.Size()
	p := width / pixelsPerFoot // Number of whole feet
	q := p * pixelsPerFoot     // Number of pixels in whole feet
	r := width - q             // Number of pixels in partial foot
	feet := feetStorage[:(width+pixelsPerFoot-1)/pixelsPerFoot]
	for y := int32(0); y < height; y++ {
		for i := 0; i < n; i += 2 {
			accumulateToFeet(
				(*[2]cvec)(unsafe.Pointer(&w[i])),
				(*[2]u13)(unsafe.Pointer(&u[i])), feet)
		}
		rotate(w, v)
		feetToPixel(feet[:p], &clut, pm.Row(y))
		if r != 0 {
			feetToPixel(feet[p:p+1], &clut, tmpStorage[:])
			copy(pm.Row(y)[q:q+r], tmpStorage[:r])
		}
	}
}
