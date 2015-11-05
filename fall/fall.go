package fall

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

const (
	tickHalfWidth   = 8
	tickHalfHeight  = 2
	dotTimeInterval = 0.3
)

// Description of Frequon needed by Draw.
type Invader struct {
	Progress  float32      // Verical coordinate of tic mark
	Amplitude float32      // Horizonal coordinate of tic mark
	Color     nimble.Pixel // color of tick mark
}

var background nimble.PixMap

// Init initializes state used by Draw.  Init should be called once with width and
// height values matching the size of the PixMap passed to future calls to Draw.
func Init(width, height int32) {
	background = nimble.MakePixMap(width, height, make([]nimble.Pixel, height*width), width)
	background.Fill(nimble.Black)
}

// Time the most recent "trailing dots" were drawn.
var lastDotTime float64

// Draw the "fall" view onto the given PixMap.
func Draw(pm nimble.PixMap, invaders []Invader) {
	if pm.Width() != background.Width() || pm.Height() != background.Height() {
		panic("fall.Draw: pm and background differ in size")
	}

	drawDot := false
	time := nimble.Now()
	if time-lastDotTime >= dotTimeInterval {
		lastDotTime = time
		drawDot = true
	}

	pm.Copy(0, 0, &background)
	xScale := float32(pm.Width() - tickHalfWidth)   // Avoid clipping tic on right side
	yScale := float32(pm.Height() - tickHalfHeight) // Avoid clippling tic on bottom
	xOffset := float32(0)
	yOffset := float32(0)
	for _, inv := range invaders {
		x := int32(inv.Amplitude*xScale + xOffset)
		y := int32(inv.Progress*yScale + yOffset)
		r := nimble.Rect{
			Left:   x - tickHalfWidth,
			Top:    y - tickHalfHeight,
			Right:  x + tickHalfWidth,
			Bottom: y + tickHalfHeight,
		}
		pm.DrawRect(r, inv.Color)
		background.DrawRect(r, nimble.Black)
		if drawDot {
			if doty := r.Top - 1; background.Contains(x, doty) {
				background.SetPixel(x, doty, inv.Color)
			}
		}
	}
}
