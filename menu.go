// Menus for Frequon Invaders

package main

import (
	"fmt"
	"github.com/ArchRobison/FrequonInvaders/menu"
	"github.com/ArchRobison/FrequonInvaders/universe"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
)

var (
	fileMenu     = menu.Menu{Label: "File"}
	displayMenu  = menu.Menu{Label: "Display", Items: []menu.ItemInterface{autoGain}}
	RatingsMenu  = menu.Menu{Label: "Ratings"}
	invadersMenu = menu.Menu{Label: "Invaders"}
	colorMenu    = menu.Menu{Label: "Color"}
)

var autoGain = menu.MakeCheckItem("Autogain", true, nil)

var menuBar = []*menu.Menu{}

// State of stationary/moving radio buttons.
var letFrequonsMove = menu.RadioState{OnSelect: func(value int) {
	if value == 0 {
		universe.SetVelocityMax(0)
	} else {
		universe.SetVelocityMax(30. / 1440. * math32.Sqrt(float32(screenWidth*screenHeight)))
	}
}}

// State of "maximum number of Frequons" buttons.
var maxFrequon = menu.RadioState{Value: 1, OnSelect: func(value int) {
	universe.SetNLiveMax(value)
}}

var beginGameItem, trainingItem, exitItem *menu.SimpleItem

var peek = menu.MakeCheckItem("peek", false, universe.SetShowAlways)

// FIXME -split this routine into two routines, one for setting up the menus
// and one for setting other state.
func setMode(m mode) {
	menuBarWasPresent := len(menuBar) > 0
	fourierIsVisible = false
	fallIsVisible = false
	radarIsVisible = false
	scoreIsVisible = false
	dividerCount = 0
	switch m {
	case modeSplash, modeName, modeVanity:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &RatingsMenu}
		fileMenu.Items = []menu.ItemInterface{
			beginGameItem,
			trainingItem,
			exitItem,
		}
		exitItem.Flags |= menu.Separator
	case modeTraining:
		menuBar = []*menu.Menu{&fileMenu, &displayMenu, &invadersMenu, &colorMenu}
		list := []menu.ItemInterface{
			peek,
			menu.MakeRadioItem("stationary", &letFrequonsMove, 0),
			menu.MakeRadioItem("moving", &letFrequonsMove, 1),
		}
		for k := 0; k <= 13; k++ {
			list = append(list, menu.MakeRadioItem(fmt.Sprintf("%v", k), &maxFrequon, k))
		}
		invadersMenu.Items = list
		startBootSequence()
	case modeGame:
		menuBar = menuBar[:0]
		startBootSequence()
	}
	currentMode = m
	if (len(menuBar) != 0) != menuBarWasPresent {
		// Menu bar appeared or disappeared, so repartition
		partitionScreen(screenWidth, screenHeight)
	}
}

func initMenuItem() {
	beginGameItem = menu.MakeSimpleItem("Begin Game", func() {
		setMode(modeGame)
	})
	trainingItem = menu.MakeSimpleItem("Training", func() {
		setMode(modeTraining)
	})
	exitItem = menu.MakeSimpleItem("Exit", func() {
		nimble.Quit()
	})
}
