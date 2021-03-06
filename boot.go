package main

import (
	"github.com/ArchRobison/FrequonInvaders/sound"
	"github.com/ArchRobison/FrequonInvaders/teletype"
)

var (
	bootSequenceIndex  = -1         // Index of next boot operation; or -1 if done booting.
	bootSequenceFrac   float32      // Boot time accumulated since last step
	bootSequencePeriod float32 = .5 // Seconds between boot steps
)

func startBootSequence() {
	bootSequenceIndex = 0
	bootSequenceFrac = 0
	teletype.DisplayCursor(false)
	teletype.Reset()
}

/* The "boot sequence" was created in the 1990's to create eye candy while
   the program was slowly computing lookup tables.  By the mid-2000s machines
   were so fast that it has no practical purpose anymore.  But to retain
   the original look of Frequon Invaders, it's done nonethless, with the
   teletype techno-babble.  It's like the flutes on concrete columns. */
func advanceBootSequence(dt float32) {
	if bootSequenceIndex < 0 || bootSequenceIndex >= 10 {
		return
	}
	bootSequenceFrac += dt
	if bootSequenceFrac < bootSequencePeriod {
		return
	}
	bootSequenceFrac = 0
	n := bootSequenceIndex
	bootSequenceIndex = n + 1
	if 1 <= n && n <= 8 {
		teletype.PrintUpper(grammar.Generate(rune('0' + n)))
		teletype.PrintChar('\n')
	}
	if 0 < n && n <= 9 {
		sound.Play(sound.Wobble, float32(n)*0.2+0.25)
	}
	switch n {
	case 1:
		break
	case 2, 3, 4:
		dividerCount = n - 1
	case 5:
		fallIsVisible = true
	case 6:
		radarIsVisible = true
		radarIsRunning = true
	case 7:
		scoreIsVisible = true
	case 9:
		// C++ original does following actions for n==8, but that hides the 8th techobabble.
		// So wait until n==9.
		fourierIsVisible = true
		setZoom(zoomGrow)
		teletype.Reset()
	}
}
