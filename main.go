package main

import (
	_ "fmt"
	"github.com/ArchRobison/FrequonInvaders/coloring"
	"github.com/ArchRobison/FrequonInvaders/fall"
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/radar"
	"github.com/ArchRobison/FrequonInvaders/score"
	"github.com/ArchRobison/FrequonInvaders/sprite"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

var winTitle string = "Go-SDL2 Render"
var winWidth, winHeight int = 800, 600

type context struct {
}

var harmonicStorage [universe.MaxCritter]fourier.Harmonic

func drawFrequons(pm nimble.PixMap) {
	c := universe.Zoo
	h := harmonicStorage[:len(c)]

	// Compute L1 norm of amplitudes
	norm := float32(0)
	for i := range h {
		norm += math32.Abs(c[i].Amplitude)
	}
	invNorm := 1 / norm

	// Set up harmonics
	// (cx,cy) is center of fourier view
	cx, cy := 0.5*float32(pm.Width()), 0.5*float32(pm.Height())
	α, β := -0.5*cx, -0.5*cy
	const ωScale = 0.001
	for i := range h {
		ωx := (c[i].Sx - cx) * ωScale
		ωy := (c[i].Sy - cy) * ωScale
		h[i].Ωx = ωx
		h[i].Ωy = ωy
		h[i].Phase = α*ωx + β*ωy
		// Scale amplitude so that DFT values fit within domain of color lookup table.
		h[i].Amplitude = c[i].Amplitude * invNorm
	}
	fourier.Draw(pm, h)
}

var white = nimble.Gray(1)

var scoreCounter int
var lastTime float64

func updateClock() (dt float32) {
	t := nimble.Now()
	if lastTime > 0 {
		dt = float32(t - lastTime)
	} else {
		dt = 0
	}
	lastTime = t
	return
}

func (context) Render(pm nimble.PixMap) {
	dt := updateClock()

	// Update universe
	x, y := nimble.MouseWhere()
	xf, yf := fourierPort.RelativeToLeftTop(x, y)
	universe.Update(dt, xf, yf)

	// Draw dividers
	for _, r := range divider {
		pm.DrawRect(r, white)
	}

	// Fourier view
	drawFrequons(pm.Intersect(fourierPort))
	for k := 1; k < len(universe.Zoo); k++ {
		c := &universe.Zoo[k]
		if c.Show {
			i := c.ImageIndex()
			if i < len(critterSeq[k]) {
				// FIXME - draw only if close
				sprite.Draw(pm.Intersect(fourierPort), int32(math32.Round(c.Sx)), int32(math32.Round(c.Sy)), critterSeq[k][i], nimble.Gray(0.5)) // FIXME - use pastel instead of gray
			}
		}
	}
	sprite.Draw(pm.Intersect(fourierPort), xf, yf, critterSeq[0][0], white)

	// Fall view
	// FIXME - use storage buffer instead of creating new array each time?
	inv := make([]fall.Invader, len(universe.Zoo)-1)
	for k := range inv {
		c := &universe.Zoo[k+1]
		// FIXME - use pastel for color
		inv[k] = fall.Invader{
			Progress:  c.Progress,
			Amplitude: c.Amplitude,
			Color:     nimble.RGB(1, 0.5, 1)}
	}
	fall.Draw(pm.Intersect(fallPort), inv)

	// Radar view
	radar.Draw(pm.Intersect(radarPort), true)

	// Score
	score.Draw(pm.Intersect(scorePort), scoreCounter>>4)
	scoreCounter++ // FIXME - temporary hack
}

var fallPort, radarPort, scorePort, fourierPort nimble.Rect
var divider [3]nimble.Rect
var critterSeq [universe.MaxCritter][]sprite.Sprite

func initCritterSprites(width, height int32) {
	radius := height / 32
	const frameCount = 60
	for k := range critterSeq {
		critterSeq[k] = sprite.MakeAnimation(int(radius), k == 0, frameCount)
	}
}

func (context) Init(width, height int32) {
	panelWidth := width / 8 / 6 * 6

	universe.Init(width, height)

	fourierPort = nimble.Rect{Left: panelWidth + 1, Top: 0, Right: width, Bottom: height}

	scoreBottom := height
	scoreTop := scoreBottom - panelWidth/6
	scorePort = nimble.Rect{Left: 0, Top: scoreTop, Right: panelWidth, Bottom: scoreBottom}

	radarBottom := scoreTop - 1
	radarTop := radarBottom - panelWidth
	radarPort = nimble.Rect{Left: 0, Top: radarTop, Right: panelWidth, Bottom: radarBottom}

	fallBottom := radarTop - 1
	fallTop := int32(0)
	fallPort = nimble.Rect{Left: 0, Top: fallTop, Right: panelWidth, Bottom: fallBottom}

	divider[0] = nimble.Rect{Left: panelWidth, Top: 0, Right: panelWidth + 1, Bottom: height}
	divider[1] = nimble.Rect{Left: 0, Top: fallBottom, Right: panelWidth, Bottom: radarTop}
	divider[2] = nimble.Rect{Left: 0, Top: radarBottom, Right: panelWidth, Bottom: scoreTop}

	fall.Init(fallPort.Size())
	radar.Init(radarPort.Size())
	score.Init(scorePort.Size())
	radar.SetColoring(coloring.HasEverything)
	fourier.Init(fourierPort.Size())
	fourier.SetColoring(coloring.HasEverything)
	initCritterSprites(width, height)
}

func main() {
	nimble.AddRenderClient(context{})
	nimble.Run()
}
