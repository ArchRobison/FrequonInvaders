package fourier

import (
	"github.com/ArchRobison/FrequonInvaders/cmplx64"
	. "github.com/ArchRobison/NimbleDraw"
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

var clut[clutSize][clutSize]Pixel

type colorMap interface {
     Color(x, y float32) (r, g, b float32) 
}
        
func SetColoring(cm colorMap) {
	for i := 0; i < clutSize; i++ {
		y := clutCoor(i)
		for j := 1; j < clutSize; j++ {
			x := clutCoor(j)
			clut[i][j] = RGB(cm.Color(x, y))
		}
	}
}

func Init(width, height int32) {
}

func Draw(pm PixMap, harmonics []Harmonic) {
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
			w[i] *= v[i]	// Rotate w by v
		}
		row := pm.Row(y)
		for x := range row {
		    const offset float32 = clutCenter + 0.5
			s := complex(offset,offset)
			for i := 0; i < n; i++ {
				s += z[i]
				z[i] *= u[i]
			}
			row[x] = clut[int(imag(s))][int(real(s))]
		}
	}
}
