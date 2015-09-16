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
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
	"time"
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
			setMode(modeVanity)
		case modeGame:
			// FIXME - need to handle as quit with current score
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

func (context) Render(pm nimble.PixMap) {
	dt := updateClock()

	// Update universe
	xf, yf := fourierPort.RelativeToLeftTop(mouseX, mouseY)
	universe.Update(dt, xf, yf)

	// Draw dividers
	for _, r := range divider {
		pm.DrawRect(r, nimble.White)
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
	sprite.Draw(pm.Intersect(fourierPort), xf, yf, critterSeq[0][0], nimble.White)

	// Fall view
	inv := invStorage[0 : len(universe.Zoo)-1]
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
	score.Draw(pm.Intersect(scorePort), universe.NKill())

	// Menu bar
	if len(menuBar) > 0 {
		menu.DrawMenuBar(pm, menuBar)
	}
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

var screenWidth, screenHeight int32

func (context) Init(width, height int32) {
	screenWidth, screenHeight = width, height
	initCritterSprites(width, height)
	setMode(modeSplash) // N.B. also causes partitionScreen to be called
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
	fourier.Init(fourierPort.Size())
	fourier.SetColoring(coloring.AllBits)
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

var (
	fileMenu     = menu.Menu{Label: "File"}
	displayMenu  = menu.Menu{Label: "Display"}
	RatingsMenu  = menu.Menu{Label: "Ratings"}
	invadersMenu = menu.Menu{Label: "Invaders"}
	colorMenu    = menu.Menu{Label: "Color"}
)

var menuBar = []*menu.Menu{}

type simpleItem struct {
	menu.MenuItem
	onSelect func()
}

func (m *simpleItem) OnSelect() {
	m.onSelect()
}

var beginGameItem, trainingItem, exitItem *simpleItem

func MakeItem(label string, f func()) *simpleItem {
	return &simpleItem{menu.MenuItem{Label: label}, f}
}

func setMode(m mode) {
	menuBarWasPresent := len(menuBar) > 0
	switch m {
	case modeSplash, modeName, modeVanity:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &RatingsMenu}
		fileMenu.Items = []menu.MenuItemInterface{
			beginGameItem,
			trainingItem,
			exitItem,
		}
		exitItem.Flags |= menu.Separator
	case modeTraining:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &invadersMenu, &colorMenu}
	case modeGame:
		menuBar = menuBar[:0]
	}
	currentMode = m
	if (len(menuBar) != 0) != menuBarWasPresent {
		// Menu bar appeared or disappeared, so repartition
		partitionScreen(screenWidth, screenHeight)
	}
}

func initMenuItem() {
	beginGameItem = MakeItem("Begin Game", func() {
		setMode(modeGame)
	})
	trainingItem = MakeItem("Training", func() {
		setMode(modeTraining)
	})
	exitItem = MakeItem("Quit", func() {
		nimble.Quit()
	})
}

func main() {
	initMenuItem()
	rand.Seed(time.Now().UnixNano())
	nimble.AddRenderClient(context{})
	nimble.AddMouseObserver(context{})
	nimble.AddKeyObserver(context{})
	nimble.Run()
}
