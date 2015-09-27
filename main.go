package main

import (
	"fmt"
	"github.com/ArchRobison/FrequonInvaders/coloring"
	"github.com/ArchRobison/FrequonInvaders/fall"
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/menu"
	"github.com/ArchRobison/FrequonInvaders/radar"
	"github.com/ArchRobison/FrequonInvaders/score"
	"github.com/ArchRobison/FrequonInvaders/sprite"
	"github.com/ArchRobison/FrequonInvaders/teletype"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
	"time"
)

const title = "Frequon Invaders 2.3"
const edition = "(Go Edition)"

const debugMode = true

var winTitle string = title

var winWidth, winHeight int = 1024, 768

type context struct {
}

var harmonicStorage [universe.MaxCritter]fourier.Harmonic

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
	// FIXME - use scheme from BoxFractionAndScheme when drawing the fourier view
	fracX, fracY, _ := universe.BoxFractionAndScheme()
	sizeX, sizeY := pm.Size()

	// Set up harmonics
	// (cx,cy) is center of fourier view
	cx, cy := 0.5*float32(sizeX)*fracX, 0.5*float32(sizeY)*fracY
	α, β := -0.5*cx, -0.5*cy
	const ωScale = 0.001
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
	}), h)
	if marginX != 0 || marginY != 0 {
		pm.DrawRect(nimble.Rect{Left: 0, Right: sizeX, Top: 0, Bottom: marginY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: 0, Right: sizeX, Top: sizeY - marginY, Bottom: sizeY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: 0, Right: marginX, Top: marginY, Bottom: sizeY - marginY}, nimble.Black)
		pm.DrawRect(nimble.Rect{Left: sizeX - marginX, Right: sizeX, Top: marginY, Bottom: sizeY - marginY}, nimble.Black)
	}
}

func drawFrequonsSpatial(pm nimble.PixMap, xf, yf int32) {
	// Draw in spatial domain
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

var invStorage = make([]fall.Invader, universe.MaxCritter)

func (context) KeyDown(k nimble.Key) {
	switch k {
	case nimble.KeyEscape:
		switch currentMode {
		case modeSplash, modeVanity:
			nimble.Quit()
		case modeTraining:
			youLose()
		case modeGame:
			// FIXME - need to handle as quit with current score
			youLose()
		case modeName:
			// FIXME - ask for confirmation
			nimble.Quit()
		}
	}
}

var mouseX, mouseY int32

func (context) ObserveMouse(event nimble.MouseEvent, x, y int32) {
	mouseX, mouseY = x, y
	for _, m := range menuBar {
		m.TrackMouse(event, x, y)
	}
}

var (
	fourierIsVisible = false
	fallIsVisible    = false
	radarIsVisible   = false
	radarIsRunning   = false
	scoreIsVisible   = false
	dividerCount     = 0
)

func (context) Render(pm nimble.PixMap) {
	dt := updateClock()

	// Advance the boot sequence
	advanceBootSequence(dt)

	// Draw dividers
	for i, r := range divider {
		if i < dividerCount {
			pm.DrawRect(r, nimble.White)
		} else {
			pm.DrawRect(r, nimble.Black)
		}
	}

	if fourierIsVisible {
		// Update universe
		xf, yf := fourierPort.RelativeToLeftTop(mouseX, mouseY)
		switch universe.Update(dt, xf, yf) {
		case universe.GameWin:
			youWin()
		case universe.GameLose:
			youLose()
		}
		updateZoom(dt)

		// Fourier view
		drawFrequonsFourier(pm.Intersect(fourierPort))
		drawFrequonsSpatial(pm.Intersect(fourierPort), xf, yf)
	} else {
		// Teletype view
		teletype.Draw(pm.Intersect(fourierPort))
	}

	// Fall view
	if fallIsVisible {
		inv := invStorage[0 : len(universe.Zoo)-1]
		for k := range inv {
			c := &universe.Zoo[k+1]
			inv[k] = fall.Invader{
				Progress:  c.Progress,
				Amplitude: c.Amplitude,
				Color:     Pastel[c.Id][0]}
		}
		fall.Draw(pm.Intersect(fallPort), inv)
	} else {
		pm.DrawRect(fallPort, nimble.Black)
	}

	// Radar view
	if radarIsVisible {
		radar.Draw(pm.Intersect(radarPort), radarIsRunning)
	} else {
		pm.DrawRect(radarPort, nimble.Black)
	}

	// Score
	if scoreIsVisible {
		score.Draw(pm.Intersect(scorePort), universe.NKill())
	} else {
		pm.DrawRect(scorePort, nimble.Black)
	}

	// Menu bar
	if len(menuBar) > 0 {
		menu.DrawMenuBar(pm, menuBar)
	}
}

var fallPort, radarPort, scorePort, fourierPort nimble.Rect
var divider [3]nimble.Rect
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

const NPastel = 32

var Pastel [universe.MaxCritter][NPastel]nimble.Pixel

func initPastel() {
	for k := range Pastel {
		coloring.PastelFade(Pastel[k][:], k, universe.MaxCritter)
	}
}

var screenWidth, screenHeight int32

func (context) Init(width, height int32) {
	screenWidth, screenHeight = width, height
	initCritterSprites(width, height)
	initPastel()
	setMode(modeSplash) // N.B. also causes partitionScreen to be called
	teletype.Init("Characters.png")
	teletype.Print(title + "\n")
	teletype.Print(edition + "\n")
	teletype.Print("By Arch D. Robison\n")
	if debugMode {
		teletype.Print("[debug mode]\n")
	}
}

func partitionScreen(width, height int32) {
	if width <= 0 || height <= 0 {
		panic(fmt.Sprintf("partitionScreen: width=%v height=%v", width, height))
	}

	panelWidth := width / 8 / 6 * 6

	var menuHeight int32 = 0
	if len(menuBar) > 0 {
		_, menuHeight = menuBar[0].TabSize()
	}

	fourierPort = nimble.Rect{Left: panelWidth + 1, Top: menuHeight, Right: width, Bottom: height}
	universe.Init(fourierPort.Size())

	scoreBottom := height
	scoreTop := scoreBottom - panelWidth/6
	scorePort = nimble.Rect{Left: 0, Top: scoreTop, Right: panelWidth, Bottom: scoreBottom}

	radarBottom := scoreTop - 1
	radarTop := radarBottom - panelWidth
	radarPort = nimble.Rect{Left: 0, Top: radarTop, Right: panelWidth, Bottom: radarBottom}

	fallBottom := radarTop - 1
	fallTop := menuHeight
	fallPort = nimble.Rect{Left: 0, Top: fallTop, Right: panelWidth, Bottom: fallBottom}

	divider[0] = nimble.Rect{Left: panelWidth, Top: menuHeight, Right: panelWidth + 1, Bottom: height}
	divider[1] = nimble.Rect{Left: 0, Top: fallBottom, Right: panelWidth, Bottom: radarTop}
	divider[2] = nimble.Rect{Left: 0, Top: radarBottom, Right: panelWidth, Bottom: scoreTop}

	fall.Init(fallPort.Size())
	radar.Init(radarPort.Size())
	score.Init(scorePort.Size())
	radar.SetColoring(coloring.AllBits)
	// There is no fourer.Init routine since the zoom changes its display size
	fourier.SetColoring(coloring.AllBits)
}

func youLose() {
	setZoom(zoomShrink)
	radarIsRunning = false
}

func youWin() {
	setZoom(zoomShrink)
}

type mode int8

var currentMode mode

const (
	modeSplash = mode(iota)
	modeTraining
	modeGame
	modeName
	modeVanity
)

func main() {
	nimble.SetWindowTitle(title + " " + edition)
	initMenuItem()
	rand.Seed(time.Now().UnixNano())
	nimble.AddRenderClient(context{})
	nimble.AddMouseObserver(context{})
	nimble.AddKeyObserver(context{})
	nimble.Run()
}
