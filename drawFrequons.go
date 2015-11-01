package main

import (
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/sprite"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

var critterSeq [universe.MaxCritter][]sprite.Sprite

func initCritterSprites(width, height int32) {
	// The factor .75/64. should yield the same size as the original game
	// for a 1920x1080 display
	radius := int32(math32.Sqrt(float32(width*height)) * (.75 / 64.))
	const frameCount = 60
	for k := range critterSeq {
		critterSeq[k] = sprite.MakeAnimation(int(radius), k == 0, frameCount)
	}
}

var harmonicStorage [universe.MaxCritter]fourier.Harmonic

// drawFrequeonsFourier draws the frequency-domain representation of Frequons.
func drawFrequonsFourier(pm nimble.PixMap) {
	c := universe.Zoo
	h := harmonicStorage[:len(c)]

	var ampScale float32
	if autoGain.Value {
		// Compute L1 norm of amplitudes
		norm := float32(0)
		for i := range c {
			norm += math32.Abs(c[i].Amplitude)
		}
		ampScale = 1 / norm
	} else {
		ampScale = 1 / float32(len(c))
	}
	fracX, fracY := universe.BoxFraction()
	fracX *= zoomCompression
	fracY *= zoomCompression
	sizeX, sizeY := pm.Size()

	// Set up harmonics
	// (cx,cy) is center of fourier view
	cx, cy := 0.5*float32(sizeX)*fracX, 0.5*float32(sizeY)*fracY
	α, β := -0.5*cx, -0.5*cy
	const ωScale = 0.0005 // FIXME - should vary with screen size
	for i := range h {
		ωx := (c[i].Sx - cx) * ωScale
		ωy := (c[i].Sy - cy) * ωScale
		h[i].Ωx = ωx
		h[i].Ωy = ωy
		h[i].Phase = α*ωx + β*ωy
		// Scale amplitude so that DFT values fit within domain of color lookup table.
		h[i].Amplitude = c[i].Amplitude * ampScale
	}

	marginX := int32(math32.Round(0.5 * float32(sizeX) * (1 - fracX)))
	marginY := int32(math32.Round(0.5 * float32(sizeY) * (1 - fracY)))
	fourier.Draw(pm.Intersect(nimble.Rect{
		Left:   marginX,
		Right:  sizeX - marginX,
		Top:    marginY,
		Bottom: sizeY - marginY,
	}), h, universe.Scheme())
	if marginX != 0 || marginY != 0 {
		pm.DrawRect(nimble.Rect{Left: 0, Right: sizeX, Top: 0, Bottom: marginY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: 0, Right: sizeX, Top: sizeY - marginY, Bottom: sizeY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: 0, Right: marginX, Top: marginY, Bottom: sizeY - marginY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: sizeX - marginX, Right: sizeX, Top: marginY, Bottom: sizeY - marginY}, nimble.Black)
	}
}

// drawFrequeonsSpatial draws the spatial-domain representation of Frequons.
func drawFrequonsSpatial(pm nimble.PixMap, xf, yf int32) {
	for k := 1; k < len(universe.Zoo); k++ {
		c := &universe.Zoo[k]
		if c.Show {
			i := c.ImageIndex()
			if i < len(critterSeq[k]) {
				// FIXME - draw only if close
				// FIXME - fade according to distance
				sprite.Draw(pm, int32(math32.Round(c.Sx)), int32(math32.Round(c.Sy)), critterSeq[k][i], Pastel[c.Id][0])
			}
		}
	}
	sprite.Draw(pm, xf, yf, critterSeq[0][0], nimble.White)
}
