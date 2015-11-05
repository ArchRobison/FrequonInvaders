package main

import (
	"github.com/ArchRobison/FrequonInvaders/coloring"
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

var phaseRoll float32

func updatePhaseRoll(dt float32) {
	const rate = 2 * math32.Pi * (0.25) // rate in cycles per second
	phaseRoll += rate * dt
	for phaseRoll > math32.Pi {
		phaseRoll -= 2 * math32.Pi
	}
}

// drawFrequonsFourier draws the frequency-domain representation of Frequons.
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
	ωScale := 0.15 / math32.Sqrt(float32(pm.Width()*pm.Height()))
	for i := range h {
		ωx := (c[i].Sx - cx) * ωScale
		ωy := (c[i].Sy - cy) * ωScale
		h[i].Ωx = ωx
		h[i].Ωy = ωy
		h[i].Phase = α*ωx + β*ωy + phaseRoll
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
		d := int(math32.Hypot(float32(xf)-c.Sx, float32(yf)-c.Sy))
		if c.Show || d < NPastel {
			i := c.ImageIndex()
			if i < len(critterSeq[k]) {
				j := 0
				if !c.Show {
					j = d
				}
				sprite.Draw(pm, int32(math32.Round(c.Sx)), int32(math32.Round(c.Sy)), critterSeq[k][i], coloring.Pastel(c.Id, j))
			}
		}
	}
	if fourierPort.Contains(mouseX, mouseY) {
		sprite.Draw(pm, xf, yf, critterSeq[0][0], nimble.White)
	}
}
