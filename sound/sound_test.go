package sound

import (
	"github.com/ArchRobison/Gophetica/nimble"
	"testing"
)

type context struct{}

func (*context) Init(int32, int32) {
}

var counter int

func (*context) Render(pm nimble.PixMap) {
	switch counter {
	case 0:
		nimble.PlaySound(Wobble, 1, 1)
	case 30:
		nimble.PlaySound(Twang, 1, 1.5)
	case 60:
		nimble.PlaySound(Bell, 1, 1)
	case 90:
		nimble.PlaySound(Bell, 1, 0.75)
	case 120:
		nimble.PlaySound(AntiTwang, 1, 1)
	case 150:
		nimble.PlaySound(Broken, 1, 1)
	case 180:
		nimble.Quit()
	}
	counter++
}

// Requires visual inspection of stdout
func TestTwang(t *testing.T) {
	nimble.AddRenderClient(&context{})
	var winSpec nimble.WindowSpec = nil
	nimble.Run(winSpec)
}
