package main

import (
	"github.com/ArchRobison/FrequonInvaders/universe"
)

var (
	zoomRate   float32
	zoomAmount float32
)

const (
	// 60/64 ratio derived from original C++ version assuming it ran at 60 frames/sec
	zoomGrow   = 60. / 64.
	zoomShrink = -zoomGrow
)

// Set zoom rate.  Argument should be zoomGrow or zoomShrink for normal speed.
func setZoom(dir float32) {
	zoomRate = dir
	if debugMode && benchmarkMode {
		// Accelerate zoom
		zoomRate *= 20
	}
	if dir > 0 {
		zoomAmount = 0
		// Disable birthing of Frequons until zoom completes.
		universe.SetNLiveMax(0)
	} else {
		// Use current zoomAmount
	}
}

func updateZoom(dt float32) {
	z := zoomAmount + zoomRate*dt
	if z > 1 {
		z = 1
		if zoomAmount < 1 {
			// Have reached full zoom.  Enable birthing of Frequons.
			if currentMode == modeGame {
				universe.SetNLiveMax(1)
			} else if debugMode && benchmarkMode {
				universe.SetNLiveMax(8)
			} else {
				universe.SetNLiveMax(maxFrequon.Value)
			}
		}
	} else if z < 0 {
		z = 0
	}
	zoomAmount = z
	if z != 0 {
		const min, max = 1., 16.
		universe.SetBoxFraction(min / (min + (max-min)*(1-z)))
	} else {
		fourierIsVisible = false
		endGame()
	}
}
