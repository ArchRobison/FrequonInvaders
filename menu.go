// Menus for Frequon Invaders

package main

import (
	"fmt"
	"github.com/ArchRobison/FrequonInvaders/coloring"
	"github.com/ArchRobison/FrequonInvaders/fourier"
	"github.com/ArchRobison/FrequonInvaders/menu"
	"github.com/ArchRobison/FrequonInvaders/teletype"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

var (
	fileMenu     = menu.Menu{Label: "File"}
	displayMenu  = menu.Menu{Label: "Display", Items: []menu.ItemInterface{autoGain}}
	ratingsMenu  = menu.Menu{Label: "Ratings"}
	invadersMenu = menu.Menu{Label: "Invaders"}
	colorMenu    = menu.Menu{Label: "Color"}
)

// Items for "File" menu
var (
	beginGameItem, trainingItem *menu.SimpleItem

	exitItem = menu.MakeSimpleItem("Exit", func() {
		nimble.Quit()
	})
)

// Items for "Display" menu
var autoGain = menu.MakeCheckItem("Autogain", true, nil)

// Items for "Ratings" menu
var (
	highScores *menu.SimpleItem
	cpuSpeed   = menu.MakeSimpleItem("CPU Speed", func() {
		teletype.Reset()
		teletype.PrintUpper(fmt.Sprintf("HFT SPEED = %.1f GFlops\n", fourier.Benchmark()*1e-9))
	})
)

// State of stationary/moving radio buttons on the "Invaders" menu.
var letFrequonsMove = menu.RadioState{OnSelect: func(value int) {
	if value == 0 {
		universe.SetVelocityMax(0)
	} else {
		universe.SetVelocityMax(30. / 1440. * math32.Sqrt(float32(screenWidth*screenHeight)))
	}
}}

// State of "maximum number of Frequons" buttons on the "Invaders" menu.
var maxFrequon = menu.RadioState{Value: 1, OnSelect: func(value int) {
	universe.SetNLiveMax(value)
}}

// Entry for "Color" menu
type colorSchemeItem struct {
	name    string              // Name of scheme
	missing coloring.SchemeBits // Aspect that is missing from the color scheme
}

var colorSchemeInfo = []colorSchemeItem{
	{"Complex", 0},
	{"Real", coloring.ImagBit},
	{"Imaginary", coloring.RealBit},
	{"Magnitude", coloring.PhaseBit},
	{"Phase", coloring.MagnitudeBit},
}

// State of the "Color" buttons.
var colorSchemeSelect = menu.RadioState{OnSelect: func(value int) {
	universe.SetScheme(coloring.AllBits &^ colorSchemeInfo[value].missing)
}}

var menuBar = []*menu.Menu{}

var rollPhase = menu.MakeCheckItem("phase roll", false, nil)

func setMenus(m mode) {
	menuBarWasPresent := len(menuBar) > 0
	switch m {
	case modeSplash, modeName, modeVanity:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &ratingsMenu}
		fileMenu.Items = []menu.ItemInterface{
			beginGameItem,
			trainingItem,
			exitItem,
		}
		exitItem.Flags |= menu.Separator
		ratingsMenu.Items = []menu.ItemInterface{
			highScores,
			cpuSpeed,
		}
	case modeTraining:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &invadersMenu, &colorMenu}
		// Make the "Invaders" menu
		list := []menu.ItemInterface{
			menu.MakeCheckItem("peek", false, universe.SetShowAlways),
			menu.MakeRadioItem("stationary", &letFrequonsMove, 0),
			menu.MakeRadioItem("moving", &letFrequonsMove, 1),
			rollPhase,
		}
		for k := 0; k <= 13; k++ {
			list = append(list, menu.MakeRadioItem(fmt.Sprintf("%v", k), &maxFrequon, k))
		}
		for _, k := range []int{0, 3, 4} {
			list[k].GetItem().Flags |= menu.Separator
		}
		invadersMenu.Items = list
		// Make the "Color" menu
		list = make([]menu.ItemInterface, len(colorSchemeInfo))
		for k, info := range colorSchemeInfo {
			list[k] = menu.MakeRadioItem(info.name, &colorSchemeSelect, k)
		}
		colorMenu.Items = list
	case modeGame:
		menuBar = menuBar[:0]
	}
	if (len(menuBar) != 0) != menuBarWasPresent {
		// Menu bar appeared or disappeared, so repartition
		partitionScreen(screenWidth, screenHeight)
	}
}

// Do initializations that would cause "initialization loop" if
// embedded into the respective var declarations.
func initMenuItem() {
	beginGameItem = menu.MakeSimpleItem("Begin Game", func() {
		setMode(modeGame)
	})
	trainingItem = menu.MakeSimpleItem("Training", func() {
		setMode(modeTraining)
	})
	highScores = menu.MakeSimpleItem("High Scores", func() {
		setMode(modeVanity)
	})
}
