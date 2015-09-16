package menu

import (
	"fmt"
	"github.com/ArchRobison/Gophetica/nimble"
	"os"
	"testing"
)

type context struct{}

var theMenu Menu

func (*Menu) ObserveMouse(event nimble.MouseEvent, x, y int32) {
	_ = theMenu.TrackMouse(event, x, y)
}

func (*context) Init(int32, int32) {
}

func (*context) Render(pm nimble.PixMap) {
	pm.Fill(nimble.Gray(0.1))
	theMenu.Draw(pm, 50, 100)
}

type FruitItem struct {
	MenuItem
}

func (f *FruitItem) OnSelect() {
	fmt.Fprintf(os.Stderr, "%v\n", f.Label)
}

// Requires visual inspection
func TestMenu(t *testing.T) {
	i0 := FruitItem{MenuItem{Label: "Apple"}}
	i1 := FruitItem{MenuItem{Label: "Banana", Check: 'c', Flags: Separator}}
	i2 := FruitItem{MenuItem{Label: "Cherry", Check: 'o', Flags: Separator}}
	i3 := FruitItem{MenuItem{Label: "Date", Flags: Disabled}}
	theMenu = Menu{Label: "Fruits",
		Items: []MenuItemInterface{&i0, &i1, &i2, &i3}}
	nimble.AddRenderClient(&context{})
	nimble.AddMouseObserver(&theMenu)
	nimble.Run()
}
