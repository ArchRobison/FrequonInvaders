package main

import (
	"fmt"
	"github.com/ArchRobison/FrequonInvaders/coloring"
	"github.com/ArchRobison/FrequonInvaders/fall"
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/menu"
	"github.com/ArchRobison/FrequonInvaders/phrase"
	"github.com/ArchRobison/FrequonInvaders/radar"
	"github.com/ArchRobison/FrequonInvaders/score"
	"github.com/ArchRobison/FrequonInvaders/teletype"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
	"time"
)

const title = "Frequon Invaders 2.3"
const edition = "(Go Edition)"

const debugMode = true

type context struct {
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
	if debugMode {
		// Shortcuts for debuggin
		switch k {
		case 'b':
			bootSequencePeriod = 0
			setMode(modeGame)
		case 'e':
			youLose()
		case 's':
			universe.TallyKill()
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

var invStorage = make([]fall.Invader, universe.MaxCritter)

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
		if debugMode {
			tallyFourierFrame()
		}
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
		radar.Draw(pm.Intersect(radarPort), universe.Scheme(), radarIsRunning)
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
	teletype.Init("Characters.png")
	if debugMode && benchmarkMode {
		bootSequencePeriod = 0
		setMode(modeTraining)
	} else {
		setMode(modeSplash) // N.B. also causes partitionScreen to be called
		teletype.Print(title + "\n")
		teletype.Print(edition + "\n")
		teletype.Print("By Arch D. Robison\n")
		if debugMode {
			teletype.Print("[debug mode]\n")
		}
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
	fourier.Init(fourierPort.Width(), universe.MaxCritter)
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

func setMode(m mode) {
	fourierIsVisible = false
	fallIsVisible = false
	radarIsVisible = false
	scoreIsVisible = false
	dividerCount = 0
	switch m {
	case modeSplash, modeName, modeVanity:
	case modeTraining, modeGame:
		universe.BeginGame(m == modeTraining)
		// Temporarily set NLiveMax to 0.  End of boot sequence will set it to 1.
		universe.SetNLiveMax(0)
		startBootSequence()
	}
	currentMode = m
	setMenus(m)
}

func endGame() {
	teletype.Reset()
	n := universe.NKill()
	if n >= 64 {
		teletype.Print(phrase.Generate(rune('W')) + "\n")
	}
	/*
	       if false {
	   		// FIXME - put actions here for beating low score on vanity board
	   	} else {
	   		teletype.Print(phrase.Generate(...))
	       }
	*/
	setMode(modeVanity)
}

func main() {
	if debugMode {
		for _, fun := range profileStart() {
			defer fun()
		}
		if benchmarkMode {
			nimble.SetWindowSize(1920, 1080)
		}
	}
	nimble.SetWindowTitle(title + " " + edition)
	initMenuItem()
	rand.Seed(time.Now().UnixNano())
	nimble.AddRenderClient(context{})
	nimble.AddMouseObserver(context{})
	nimble.AddKeyObserver(context{})
	nimble.Run()
}
