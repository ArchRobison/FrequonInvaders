package fourier

import (
	"github.com/ArchRobison/Gophetica/cmplx64"
	"github.com/ArchRobison/Gophetica/nimble"
	"unsafe"
)

type Harmonic struct {
	Ωx, Ωy    float32 // Angular velocities
	Phase     float32 // Phase at (0,0)
	Amplitude float32 // Amplitude
}

const (
	clutSize   = 128            // Size of Clut along either axis.  Power of 2 to speed up indexing.
	clutCenter = clutSize / 2   // Clut indices corresponding to (0,0)
	clutRadius = clutCenter - 1 // Distance from center representing magnitude of 1.
)

func clutCoor(k int) (z float32) {
	const (
		clutScale  = 1.0 / clutRadius
		clutOffset = -clutCenter * clutScale
	)
	return float32(k)*clutScale + clutOffset
}

var clut [clutSize][clutSize]nimble.Pixel

type colorMap interface {
	Color(x, y float32) (r, g, b float32)
}

func SetColoring(cm colorMap) {
	for i := 0; i < clutSize; i++ {
		y := clutCoor(i)
		for j := 1; j < clutSize; j++ {
			x := clutCoor(j)
			clut[i][j] = nimble.RGB(cm.Color(x, y))
		}
	}
}

// Generic version of Draw
func DrawSlow(pm nimble.PixMap, harmonics []Harmonic) {
	n := len(harmonics)
	w := make([]complex64, n)
	u := make([]complex64, n)
	v := make([]complex64, n)
	z := make([]complex64, n)
	for i, h := range harmonics {
		w[i] = cmplx64.Rect(h.Amplitude*clutRadius, h.Phase)
		u[i] = cmplx64.Rect(1, h.Ωx)
		v[i] = cmplx64.Rect(1, h.Ωy)
	}
	for y := int32(0); y < pm.Height(); y++ {
		for i := 0; i < n; i++ {
			z[i] = w[i]
			w[i] *= v[i] // Rotate w by v
		}
		row := pm.Row(y)
		for x := range row {
			const offset float32 = clutCenter + 0.5
			s := complex(offset, offset)
			for i := 0; i < n; i++ {
				s += z[i]
				z[i] *= u[i]
			}
			row[x] = clut[int(imag(s))][int(real(s))]
		}
	}
}

var (
	feetStorage []foot
	uStorage    []w13
	vStorage    []complex64
	wStorage    []cvec
	tmpStorage  [12]nimble.Pixel
)

// Set max. width and max. harmonics
func Init(widthMax int32, harmonicLenMax int) {
	footLenMax := (widthMax + 11) / 12

	feetStorage = make([]foot, footLenMax)
	// Use dummy call to initialize "feet"
	feetToPixel(feetStorage, &clut, make([]nimble.Pixel, footLenMax*12))

	m := harmonicLenMax + harmonicLenMax&1
	uStorage = make([]w13, m)
	vStorage = make([]complex64, m)
	wStorage = make([]cvec, m)
}

func Draw(pm nimble.PixMap, harmonics []Harmonic) {
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
		u[i].w1 = cmplx64.Rect(1, 4*h.Ωx)
		u[i].w3 = cmplx64.Rect(1, 12*h.Ωx)
		v[i] = cmplx64.Rect(1, h.Ωy)
	}
	if n < m {
		// Zero the extra element.
		w[n] = cvec{}
		u[n] = w13{}
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
				(*[2]w13)(unsafe.Pointer(&u[i])), feet)
		}
		rotate(w, v)
		feetToPixel(feet[0:p], &clut, pm.Row(y))
		if q != 0 {
			feetToPixel(feet[p:p+1], &clut, tmpStorage[0:12])
			copy(pm.Row(y)[12*p:12*p+q], tmpStorage[0:q])
		}
	}
}
