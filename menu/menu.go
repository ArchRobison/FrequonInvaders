package menu

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

func init() {
	var err error
	menuFont, err = nimble.OpenFont("Roboto-Regular.ttf", 24) // FIXME - do not hardwire
	if err != nil {
		panic(err)
	}
	checkFont, err = nimble.OpenFont("unicons.1.0.ttf", 24) // FIXME - do not hardwire
	if err != nil {
		panic(err)
	}
	checkWidth, _ = checkFont.Size(" ")
	marginWidth = checkWidth/4 + 1
}

type menuItemFlag uint8

const (
	Disabled = menuItemFlag(1 << iota)
	Separator
)

type Item struct {
	Label string
	Flags menuItemFlag
	Check rune
}

func (i *Item) GetItem() *Item {
	return i
}

type ItemInterface interface {
	GetItem() *Item
	OnSelect()
}

func Add(i ItemInterface, f menuItemFlag) ItemInterface {
	i.GetItem().Flags |= f
	return i
}

type Menu struct {
	Label      string
	Items      []ItemInterface
	x, y       int32       // Upper left corner
	hilightRow uint8       // 0 = hide items, 1 = highlight none, 2+k = hilight row k
	itemHeight uint16      // Height of each item (in pixels) or tab
	itemWidth  uint16      // Width of widest item (in pixels)
	tabWidth   uint16      // Width of tab
	tabRect    nimble.Rect // Rectangle bounding the tab
	itemsRect  nimble.Rect // Rectangle bounding the items
}

func (m *Menu) TabSize() (width, height int32) {
	if m.tabWidth == 0 {
		m.computeTabSize()
	}
	return int32(m.tabWidth) + 2*marginWidth, int32(m.itemHeight) + 1
}

const (
	showNone    = 0
	hilightNone = 1
	hilightBase = 2
)

func (m *Menu) TrackMouse(e nimble.MouseEvent, x, y int32) bool {
	if m.hilightRow != showNone {
		if m.itemsRect.Contains(x, y) {
			row := (y - m.itemsRect.Top) / int32(m.itemHeight)
			switch e {
			case nimble.MouseUp:
				m.Items[row].OnSelect()
				m.hilightRow = showNone
			case nimble.MouseMove, nimble.MouseDown, nimble.MouseDrag:
				m.hilightRow = hilightBase + uint8(row)
			}
			return true
		} else {
			if e == nimble.MouseDown {
				m.hilightRow = showNone
				return m.tabRect.Contains(x, y)
			}
		}
	}
	if m.tabRect.Contains(x, y) {
		if e == nimble.MouseDown {
			if m.hilightRow == showNone {
				m.hilightRow = hilightNone
			} else {
				m.hilightRow = showNone
			}
		}
		return true
	}
	return false
}
