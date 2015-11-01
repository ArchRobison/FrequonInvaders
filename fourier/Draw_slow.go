// Generic version of Draw for platforms without assembly-language kernels
//
// +build !amd64

package fourier

func Init(widthMax int32, harmonicLenMax int) {
	// Noop
}

// Generic version of Draw
func Draw(pm nimble.PixMap, harmonics []Harmonic, cm colorMap) {
	setColoring(cm)
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
