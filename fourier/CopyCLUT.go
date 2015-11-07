// +build !release

package fourier

import "github.com/ArchRobison/Gophetica/nimble"

// CopyCLUT returns a PixMap copy the color lookup table.
// It is intended only for debugging.
func CopyCLUT() (pm nimble.PixMap) {
	pm = nimble.MakePixMap(clutSize, clutSize, make([]nimble.Pixel, clutSize*clutSize), clutSize)
	for y := int32(0); y < clutSize; y++ {
		copy(pm.Row(y), clut[y][:])
	}
	return
}
