package menu

// A "check button" menu item
type CheckItem struct {
	Item
	Value   bool
	Handler func(value bool)
}

func (m *CheckItem) OnSelect() {
	m.Value = !m.Value
	m.Handler(m.Value)
}

func (m *CheckItem) GetItem() *Item {
	if m.Value {
		m.Item.Check = 'c'
	} else {
		m.Item.Check = 0
	}
	return &m.Item
}

func MakeCheckItem(label string, value bool, onSelect func(bool)) *CheckItem {
	return &CheckItem{
		Item:    Item{Label: label},
		Value:   value,
		Handler: onSelect}
}
