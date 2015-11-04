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
	tmpStorage  [12]nimble.Pixel
)

// Init initializes storage used by Draw.  The parameters specify the
// maximum width of the PixMap and maximum number of harmonics that
// will be passed to Draw.
func Init(widthMax int32, harmonicLenMax int) {
	footLenMax := (widthMax + 11) / 12

	feetStorage = make([]foot, footLenMax)
	// Use dummy call to initialize "feet"
	feetToPixel(feetStorage, &clut, make([]nimble.Pixel, footLenMax*12))

	m := harmonicLenMax + harmonicLenMax&1
	uStorage = make([]u13, m)
	vStorage = make([]complex64, m)
	wStorage = make([]cvec, m)
}

// Draw draws a Fourier transform on the given PixMap.
func Draw(pm nimble.PixMap, harmonics []Harmonic, cm colorMap) {
	setColoring(cm)
	n := len(harmonics)
	// m = n rounded up to even
	m := n + n&1
	u := uStorage[0:m]
	v := vStorage[0:m]
	w := wStorage[0:m]
	for i, h := range harmonics {
		for k := 0; k < 4; k++ {
			c := cmplx64.Rect(h.Amplitude*clutRadius, h.Phase+float32(k+4)*h.Ωx)
			w[i].re[k] = real(c)
			w[i].im[k] = imag(c)
		}
		u[i].u1 = cmplx64.Rect(1, 4*h.Ωx)
		u[i].u3 = cmplx64.Rect(1, 12*h.Ωx)
		v[i] = cmplx64.Rect(1, h.Ωy)
	}
	if n < m {
		// Zero the extra element.
		w[n] = cvec{}
		u[n] = u13{}
		v[n] = 0
	}
	width := pm.Width()
	p := width / 12
	q := width % 12
	feet := feetStorage[0 : (width+11)/12]
	for y := int32(0); y < pm.Height(); y++ {
		for i := 0; i < n; i += 2 {
			accumulateToFeet(
				(*[2]cvec)(unsafe.Pointer(&w[i])),
				(*[2]u13)(unsafe.Pointer(&u[i])), feet)
		}
		rotate(w, v)
		feetToPixel(feet[0:p], &clut, pm.Row(y))
		if q != 0 {
			feetToPixel(feet[p:p+1], &clut, tmpStorage[0:12])
			copy(pm.Row(y)[12*p:12*p+q], tmpStorage[0:q])
		}
	}
}
