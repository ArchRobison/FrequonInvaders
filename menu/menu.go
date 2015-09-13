package menu

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

var (
	foregroundColor  = nimble.Black
	backgroundColor  = nimble.White
	itemHilightColor = nimble.RGB(0.875, 0.875, 1)
	tabHilightColor  = nimble.RGB(0, 0, 1)
)

var menuFont *nimble.Font

func init() {
	var err error
	menuFont, err = nimble.OpenFont("Roboto-Regular.ttf", 16) // FIXME - do not hardwire
	if err != nil {
		panic(err)
	}
}

type MenuItem struct {
	Label    string
	Shortcut rune
	Disabled bool
	Check    rune
}

func (mip *MenuItem) GetMenuItem() *MenuItem {
	return mip
}

type MenuItemInterface interface {
	GetMenuItem() *MenuItem
	OnSelect()
}

type Menu struct {
	Label      string
	Items      []MenuItemInterface
	x, y       int32       // Upper left corner
	hilightRow uint8       // 0 = hide items, 1 = highlight none, 2+k = hilight row k
	itemHeight uint16      // Height of each item (in pixels) or tab
	itemWidth  uint16      // Width of widest item (in pixels)
	tabWidth   uint16      // Width of tab
	tabRect    nimble.Rect // Rectangle bounding the tab
	itemsRect  nimble.Rect // Rectangle bounding the items
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

func (m *Menu) Draw(pm nimble.PixMap, x, y int32) {
	// Draw the tab
	if m.tabWidth == 0 {
		// Lazily compute tabWidth and itemHeight
		h := menuFont.Height()
		w, _ := menuFont.Size(m.Label)
		m.itemHeight = uint16(h)
		m.tabWidth = uint16(w)
	}
	var back, fore nimble.Pixel
	if m.hilightRow != showNone {
		back = tabHilightColor
		fore = backgroundColor
	} else {
		back = backgroundColor
		fore = foregroundColor
	}
	m.tabRect = nimble.MakeRect(x, y, int32(m.tabWidth), int32(m.itemHeight))
	pm.DrawRect(m.tabRect, back)
	pm.DrawText(x, y, m.Label, fore, menuFont)

	if m.hilightRow != showNone {

		// Draw the items
		pm.DrawRect(m.itemsRect, backgroundColor)
		w := int32(m.itemWidth)
		if w == 0 {
			// Lazily compute itemsWidth
			w = int32(m.tabWidth)
			for i := range m.Items {
				p := m.Items[i].GetMenuItem()
				w0, _ := menuFont.Size(p.Label)
				if w0 > w {
					w = w0
				}
			}
			m.itemWidth = uint16(w)
		}
		h := int32(m.itemHeight)
		m.itemsRect = nimble.MakeRect(x, m.tabRect.Bottom, w, h*int32(len(m.Items)))
		for i := range m.Items {
			yi := m.itemsRect.Top + h*int32(i)
			if i == int(m.hilightRow-hilightBase) {
				pm.DrawRect(nimble.MakeRect(x, yi, int32(m.itemWidth), h), itemHilightColor)
			}
			pm.DrawText(x, yi, m.Items[i].GetMenuItem().Label, foregroundColor, menuFont)
		}
	}
}
