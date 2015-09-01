// Package fall implements the "Amplitude View" view, where
// the player sees tics falling on a parabolic curve.
package fall

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

const (
	tickHalfWidth  = 8
	tickHalfHeight = 2
)

type Invader struct {
	Progress  float32 // Verical coordinate of tic mark
	Amplitude float32 // Horizonal coordinate of tic mark
	Color     nimble.Pixel
}

var background nimble.PixMap

var black nimble.Pixel = nimble.Gray(0)

func Init(width, height int32) {
	black = nimble.Gray(0)
	background = nimble.MakePixMap(width, height, make([]nimble.Pixel, height*width), width)
	background.Fill(black)
}

var lastDotTime float64

const dotTimeInterval = 0.2

func Draw(pm nimble.PixMap, invaders []Invader) {
	if pm.Width() != background.Width() || pm.Height() != background.Height() {
		panic("fall.Draw: pm and background differ")
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
		background.DrawRect(r, black)
		if drawDot {
			if doty := r.Top - 1; background.Contains(x, doty) {
				background.SetPixel(x, doty, inv.Color)
			}
		}
	}
}
