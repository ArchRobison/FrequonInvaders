package teletype

import (
	"github.com/ArchRobison/Gophetica/nimble"
	"testing"
)

type context struct{}

func (*context) Init(int32, int32) {
}

var counter int

func (*context) Render(pm nimble.PixMap) {
    switch {
		case counter==0:
		    Append("THE QUICK BROWN FOX") 
		case counter==30:	
			displayCursor = true
		    Newline()
		case counter==60:
			Append("jumped over the lazy red dog.")
		case counter>=90 && counter%10==0:
		    Backup()
	}
    counter++
	Draw(pm)
}

// Requires visual inspection 
func TestTeletype(t *testing.T) {
    Init("../Characters.png")
	nimble.AddRenderClient(&context{})
	nimble.Run()
}
