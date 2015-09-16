package menu

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

var (
	foregroundColor         = nimble.Black
	disabledForegroundColor = nimble.Gray(0.5)
	backgroundColor         = nimble.White
	itemHilightColor        = nimble.RGB(0.875, 0.875, 1)
	tabHilightColor         = nimble.RGB(0, 0, 1)
	separatorColor          = nimble.Gray(0.75)
)

var menuFont, checkFont *nimble.Font

// Each "item" row of the menu is drawn as:
// | margin checkWidth margin | margin itemWidth margin |
var (
	checkWidth  int32 // Width of "check" column (not including margins)
	marginWidth int32 // Margin in pixels between text and border
)

// Return b ? p : q
func choose(b bool, p, q nimble.Pixel) nimble.Pixel {
	if b {
		return p
	} else {
		return q
	}
}

// Compute tab size parameters
func (m *Menu) computeTabSize() {
	h := menuFont.Height()
	w, _ := menuFont.Size(m.Label)
	m.itemHeight = uint16(h) + 1
	m.tabWidth = uint16(w)
}

// Draw draws the menu with its upper left corner at (x,y)
func (m *Menu) Draw(pm nimble.PixMap, x, y int32) {
	// Draw the tab
	if m.tabWidth == 0 {
		// Lazily compute tabWidth and itemHeight
		m.computeTabSize()
	}
	var back, fore nimble.Pixel
	if m.hilightRow != showNone {
		back = tabHilightColor
		fore = backgroundColor
	} else {
		back = backgroundColor
		fore = foregroundColor
	}
	m.tabRect = nimble.MakeRect(x, y, 2*marginWidth+int32(m.tabWidth), int32(m.itemHeight))
	pm.DrawRect(m.tabRect, back)
	pm.DrawText(x+marginWidth, y, m.Label, fore, menuFont)
	pm.DrawRect(nimble.Rect{Left: m.tabRect.Left, Top: m.tabRect.Bottom, Right: m.tabRect.Right, Bottom: m.tabRect.Bottom + 1}, separatorColor)

	if m.hilightRow != showNone {
		if m.itemWidth == 0 {
			// Lazily compute itemsWidth
			w := int32(m.tabWidth)
			for i := range m.Items {
				p := m.Items[i].GetMenuItem()
				w0, _ := menuFont.Size(p.Label)
				if w0 > w {
					w = w0
				}
			}
			m.itemWidth = uint16(checkWidth + w + 3 + 4*marginWidth)
		}
		w := int32(m.itemWidth)
		h := int32(m.itemHeight)
		m.itemsRect = nimble.MakeRect(x, m.tabRect.Bottom, w, h*int32(len(m.Items)))

		r := m.itemsRect
		r.Left += 1
		r.Right -= 1
		pm.DrawRect(r, backgroundColor)

		// Draw left border
		r.Left -= 1
		r.Right = r.Left + 1
		pm.DrawRect(r, separatorColor)

		// Draw middle border
		r.Left += 1 + checkWidth + 2*marginWidth
		r.Right = r.Left + 1
		pm.DrawRect(r, separatorColor)

		// Draw right border
		r.Left = m.itemsRect.Right - 1
		r.Right = r.Left + 1
		pm.DrawRect(r, separatorColor)

		// Draw the items
		checkX := x + 1 + marginWidth
		labelX := x + 2 + 3*marginWidth + checkWidth
		for i := range m.Items {
			mi := m.Items[i].GetMenuItem()
			yi := m.itemsRect.Top + h*int32(i)
			if i == int(m.hilightRow-hilightBase) {
				pm.DrawRect(nimble.MakeRect(x, yi, int32(m.itemWidth), h), itemHilightColor)
			}
			if i == 0 {
				pm.DrawRect(nimble.Rect{x, yi, m.itemsRect.Right - 1, yi + 1}, separatorColor)
			} else if mi.Flags&Separator != 0 || i == 0 {
				pm.DrawRect(nimble.Rect{labelX - marginWidth, yi, m.itemsRect.Right - 1, yi + 1}, separatorColor)
			}
			fore := choose(mi.Flags&Disabled != 0, disabledForegroundColor, foregroundColor)
			if mi.Check != 0 {
				pm.DrawText(checkX, yi, string(mi.Check), fore, checkFont)
			}
			pm.DrawText(labelX, yi, mi.Label, fore, menuFont)
		}
	}
}

func DrawMenuBar(pm nimble.PixMap, menuBar []*Menu) {
	// Draw menu
	x := int32(0)
	for _, m := range menuBar {
		m.Draw(pm, x, 0)
		w, _ := m.TabSize()
		x += w
	}
	h := menuFont.Height()
	pm.DrawRect(nimble.Rect{Top: 0, Left: x, Bottom: h, Right: pm.Width()}, backgroundColor)
	pm.DrawRect(nimble.Rect{Top: h, Left: x, Bottom: h + 1, Right: pm.Width()}, separatorColor)
}
