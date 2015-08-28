package main

import (
    _ "fmt"
	"github.com/ArchRobison/FrequonInvaders/fall"
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/radar"
	"github.com/ArchRobison/FrequonInvaders/score"
	"github.com/ArchRobison/FrequonInvaders/sprite"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/FrequonInvaders/coloring"
	. "github.com/ArchRobison/FrequonInvaders/math32"
	. "github.com/ArchRobison/NimbleDraw"
)

var winTitle string = "Go-SDL2 Render"
var winWidth, winHeight int = 800, 600

type context struct {
}

const ωScale = 0.001

func drawFrequons(pm PixMap, selfX float32, selfY float32) {
    c:=universe.Zoo
    if len(c)<1 {
	    panic("universe.Zoo is empty")
    }
	c[0].Sx = selfX - 0.5*float32(pm.Width())
	c[0].Sy = selfY - 0.5*float32(pm.Height())
    c[0].Amplitude = -1
	h := make([]fourier.Harmonic, len(c))
	α, β := 0.5*float32(pm.Width()), -0.5*float32(pm.Height())
	for i := range h {
		ωx, ωy := c[i].Sx*ωScale, c[i].Sy*ωScale
		h[i].Ωx = ωx
		h[i].Ωy = ωy
		h[i].Phase = α*ωx + β*ωy
		h[i].Amplitude = c[i].Amplitude
	}
	fourier.Draw(pm, h)
}

var white = Gray(1)

var scoreCounter int
var lifetimeCounter float32
var spriteCounter int

func (context) Render(pm PixMap) {

	x, y := MouseWhere()
	xf, yf := fourierPort.RelativeToLeftTop(x, y)
	drawFrequons(pm.Intersect(fourierPort), float32(xf), float32(yf))
	for _, r := range divider {
		pm.DrawRect(r, white)
	}
	inv := make([]fall.Invader, 1)
	inv[0].Lifetime = lifetimeCounter
	lifetimeCounter += 0.01
	if lifetimeCounter > 1 {
		lifetimeCounter = 0
	}
	inv[0].Power = Sqrt(inv[0].Lifetime)
	inv[0].Color = RGB(1, 0.5, 1)
	fall.Draw(pm.Intersect(fallPort), inv)

	radar.Draw(pm.Intersect(radarPort), true)

	score.Draw(pm.Intersect(scorePort), scoreCounter>>4)
	scoreCounter++

	sprite.Draw(pm.Intersect(fourierPort), xf, yf, selfSeq, spriteCounter, Gray(1))
	spriteCounter = (spriteCounter + 1) % 60
}

var fallPort, radarPort, scorePort, fourierPort Rect
var divider [3]Rect
var selfSeq sprite.Seq

func (context) Init(width, height int32) {
	panelWidth := width / 8 / 6 * 6

    universe.Init(width, height)

	fourierPort = Rect{Left: panelWidth + 1, Top: 0, Right: width, Bottom: height}

	scoreBottom := height
	scoreTop := scoreBottom - panelWidth/6
	scorePort = Rect{Left: 0, Top: scoreTop, Right: panelWidth, Bottom: scoreBottom}

	radarBottom := scoreTop - 1
	radarTop := radarBottom - panelWidth
	radarPort = Rect{Left: 0, Top: radarTop, Right: panelWidth, Bottom: radarBottom}

	fallBottom := radarTop - 1
	fallTop := int32(0)
	fallPort = Rect{Left: 0, Top: fallTop, Right: panelWidth, Bottom: fallBottom}

	divider[0] = Rect{Left: panelWidth, Top: 0, Right: panelWidth + 1, Bottom: height}
	divider[1] = Rect{Left: 0, Top: fallBottom, Right: panelWidth, Bottom: radarTop}
	divider[2] = Rect{Left: 0, Top: radarBottom, Right: panelWidth, Bottom: scoreTop}

	fall.Init(fallPort.Size())
	radar.Init(radarPort.Size())
	score.Init(scorePort.Size())
	radar.SetColoring(coloring.HasEverything)
	fourier.Init(fourierPort.Size())
	fourier.SetColoring(coloring.HasEverything)

	selfRadius := height / 32
	selfSeq = sprite.MakeSeq(int(selfRadius), true, 60)
}

func main() {
	AddRenderClient(context{})
	Run()
}
