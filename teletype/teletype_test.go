package teletype

import (
	"github.com/ArchRobison/Gophetica/nimble"
	"testing"
)

type context struct{}

func (*context) Init(int32, int32) {
}

func (*context) KeyDown(k nimble.Key) {
	displayCursor = true
	if 0x20 <= k && k < 0x7F {
		PrintChar(rune(k))
	} else {
		switch k {
		case nimble.KeyReturn:
			PrintChar('\n')
		case nimble.KeyEscape:
			nimble.Quit()
		case nimble.KeyBackspace, nimble.KeyDelete:
			Backup()
		}
	}
}

var flag bool

func (*context) Render(pm nimble.PixMap) {
	if !flag {
		Print("Type some text and 'enter'.\n")
		Print("Try backspace and del.\n" +
			"Press Esc to quit.")
		flag = true
	}
	Draw(pm)
}

// Requires visual inspection
func TestTeletype(t *testing.T) {
	Init("../Characters.png")
	nimble.AddRenderClient(&context{})
	nimble.AddKeyObserver(&context{})
	nimble.Run()
}
