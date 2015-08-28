// Package fall implements the "Amplitude View" view, where
// the player sees tics falling on a parabolic curve.

package fall

import . "github.com/ArchRobison/NimbleDraw"

const (
	tickHalfWidth  = 8
	tickHalfHeight = 2
)

type Interface interface {
	Len() int
	Power(id int) float64
	TickColor(id int) Pixel
}

type Invader struct {
	Power    float32
	Lifetime float32
	Color    Pixel
}

var background PixMap
 
var black Pixel

func Init(width, height int32) {
    black = Gray(0)
    background = MakePixMap(width,height,make([]Pixel,height*width),width)
    background.Fill(black)
}

func Draw(pm PixMap, invaders []Invader) {
    if pm.Width() != background.Width() || pm.Height() != background.Height() {
	    panic("fall.Draw: pm and background differ")
    }
	pm.Copy(0,0,&background)

    xScale := float32(pm.Width()-2*tickHalfWidth)
    xOffset := float32(tickHalfWidth)
    yScale := float32(pm.Height()-2*tickHalfHeight)
    yOffset := float32(tickHalfHeight)
	for _, inv := range invaders {
		x := int32(inv.Power * xScale + xOffset)
		y := int32(inv.Lifetime * yScale + yOffset)
		r := Rect{x - tickHalfWidth, y - tickHalfHeight, x + tickHalfWidth, y + tickHalfHeight}
		pm.DrawRect(r, inv.Color)
	    background.DrawRect(r, black)
		if background.Contains(x,y) {
		    background.SetPixel(x,y,inv.Color)
	    }
	}
}
