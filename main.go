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
	"github.com/ArchRobison/FrequonInvaders/vanity"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
	"time"
)

const title = "Frequon Invaders 2.3"
const edition = "(Go Edition)"

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
		case modeTraining, modeGame:
			youLose()
		case modeName:
			// FIXME - ask for confirmation
			nimble.Quit()
		}
	}
	if currentMode == modeName {
		if 0x20 <= k && k < 0x7F {
			teletype.PrintCharUpper(rune(k))
		} else {
			switch k {
			case nimble.KeyReturn:
				acceptScore(uint8(universe.NKill()), teletype.CursorLine())
			case nimble.KeyBackspace, nimble.KeyDelete:
				teletype.Backup()
			}
		}
	}

	if devConfig {
		// Shortcuts for debugging
		switch k {
		case 'b':
			// Begin game
			bootSequencePeriod = 0
			setMode(modeGame)
		case 'e':
			// End game
			youLose()
		case 'r':
			// Reset score file
			records = make([]vanity.Record, 0)
			vanity.WriteToFile(records)
		case 's':
			// Score a point
			universe.TallyKill()
		case 't':
			// Begin training
			bootSequencePeriod = 0
			setMode(modeTraining)
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

var cursorIsVisible = true

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
		if rollPhase.Value {
			updatePhaseRoll(dt)
		}
		drawFrequonsFourier(pm.Intersect(fourierPort))
		drawFrequonsSpatial(pm.Intersect(fourierPort), xf, yf)
		tallyFourierFrame()
	} else {
		// Teletype view
		teletype.Draw(pm.Intersect(fourierPort))
	}

	// Fall view
	if fallIsVisible {
		inv := invStorage[:len(universe.Zoo)-1]
		for k := range inv {
			c := &universe.Zoo[k+1]
			inv[k] = fall.Invader{
				Progress:  c.Progress,
				Amplitude: c.Amplitude,
				Color:     pastels.Pixel(0, int32(c.Id)),
			}
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

	// Cursor
	showCursor := true
	switch currentMode {
	case modeGame:
		showCursor = false
	case modeTraining:
		showCursor = !fourierPort.Contains(mouseX, mouseY)
	}
	if showCursor != cursorIsVisible {
		nimble.ShowCursor(showCursor)
		cursorIsVisible = showCursor
	}
}

var fallPort, radarPort, scorePort, fourierPort nimble.Rect
var divider [3]nimble.Rect

var screenWidth, screenHeight int32

var pastels nimble.PixMap

func (context) Init(width, height int32) {
	screenWidth, screenHeight = width, height
	nShade := int32(math32.Round(math32.Sqrt(float32(width*height)) * (32. / 1440)))
	initCritterSprites(width, height)
	pastels = coloring.PastelPallette(universe.MaxCritter, nShade)
	teletype.Init("Characters.png")
	if benchmarking {
		bootSequencePeriod = 0
		setMode(modeTraining)
	} else {
		setMode(modeSplash) // N.B. also causes partitionScreen to be called
		if devConfig {
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

var records []vanity.Record

func acceptScore(score uint8, name string) {
	records = vanity.Insert(records, score, name)
	err := vanity.WriteToFile(records)
	setMode(modeVanity)
	if err != nil {
		teletype.Print(err.Error())
	}
}

func setMode(m mode) {
	fourierIsVisible = false
	fallIsVisible = false
	radarIsVisible = false
	scoreIsVisible = false
	dividerCount = 0
	switch m {
	case modeSplash:
		teletype.PrintUpper(title + "\n")
		teletype.PrintUpper(edition + "\n")
		teletype.PrintUpper("By Arch D. Robison\n")
		var err error
		records, err = vanity.ReadFromFile()
		if err != nil {
			teletype.Print(err.Error())
		}
	case modeName:
		teletype.Reset()
		teletype.PrintUpper(phrase.GenerateWithNumber(rune('H'), universe.NKill()) + "\n")
		teletype.PrintUpper("Please enter your name:\n")
		teletype.DisplayCursor(true)
	case modeTraining, modeGame:
		universe.BeginGame(m == modeTraining)
		// Temporarily set NLiveMax to 0.  End of boot sequence will set it to 1.
		universe.SetNLiveMax(0)
		startBootSequence()
	case modeVanity:
		teletype.Reset()
		teletype.DisplayCursor(false)
		teletype.PrintUpper("Supreme Freqs\n\nScore Player\n")
		for _, r := range records {
			teletype.PrintUpper(fmt.Sprintf("%5d %s\n", r.Score, r.Name))
		}
	}
	currentMode = m
	setMenus(m)
}

func endGame() {
	teletype.Reset()
	n := universe.NKill()
	if n >= 64 {
		teletype.PrintUpper(phrase.Generate(rune('W')) + "\n")
	}
	if vanity.IsWorthyScore(records, uint8(n)) {
		setMode(modeName)
	} else {
		setMode(modeVanity)
	}
}

type windowSpec struct {
	width, height int32
	title         string
}

func (w *windowSpec) Size() (width, height int32) {
	return w.width, w.height
}

func (w *windowSpec) Title() string {
	return w.title
}

func main() {
	var winSpec nimble.WindowSpec = nil
	if devConfig {
		for _, fun := range profileStart() {
			defer fun()
		}
		if benchmarking {
			winSpec = &windowSpec{1920, 1080, "benchmark"}
		} else {
			winSpec = &windowSpec{1024, 600, "debug"}
		}
	}
	initMenuItem()
	rand.Seed(time.Now().UnixNano())
	nimble.AddRenderClient(context{})
	nimble.AddMouseObserver(context{})
	nimble.AddKeyObserver(context{})
	nimble.Run(winSpec)
}
