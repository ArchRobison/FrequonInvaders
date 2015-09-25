package menu

type SimpleItem struct {
	Item
	onSelect func()
}

func (m *SimpleItem) OnSelect() {
	m.onSelect()
}

func MakeSimpleItem(label string, f func()) *SimpleItem {
	return &SimpleItem{Item{Label: label}, f}
}
