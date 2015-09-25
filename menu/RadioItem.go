package menu

// State of a "radio button" menu item.
// Normally shared by several RadioItem objects.
type RadioState struct {
	Value    int
	OnSelect func(value int)
}

// A "radio button" menu item
type RadioItem struct {
	Item
	state *RadioState
	value int
}

func (m *RadioItem) OnSelect() {
	s := m.state
	if s.Value != m.value {
		s.Value = m.value
		if s.OnSelect != nil {
			s.OnSelect(m.value)
		}
	}
}

func (m *RadioItem) GetItem() *Item {
	if m.state.Value == m.value {
		m.Item.Check = 'o'
	} else {
		m.Item.Check = 0
	}
	return &m.Item
}

func MakeRadioItem(label string, state *RadioState, value int) *RadioItem {
	return &RadioItem{
		Item{Label: label},
		state,
		value}
}
