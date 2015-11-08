package fourier

import (
	"github.com/ArchRobison/Gophetica/math32"
)

// euler returns e raised to the power iθ.
func euler(θ float32) complex64 {
	y, x := math32.Sincos(θ)
	return complex(x, y)
}
