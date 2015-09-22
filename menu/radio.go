package menu

// State of a "radio button" menu item
type RadioState struct {
	Value    int
	OnSelect func(value int)
}

// A "radio button" menu item
type RadioItem struct {
	Item
	target *RadioState
	value  int
}

func (m *RadioItem) OnSelect() {
	t := m.target
	if t.Value != m.value {
		t.Value = m.value
		t.OnSelect(m.value)
	}
}

func (m *RadioItem) GetItem() *Item {
	if m.target.Value == m.value {
		m.Item.Check = 'o'
	} else {
		m.Item.Check = 0
	}
	return &m.Item
}

func MakeRadioItem(label string, target *RadioState, value int) *RadioItem {
	return &RadioItem{
		Item{Label: label},
		target,
		value}
}
