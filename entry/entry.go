package entry

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DateEntry struct {
	widget.Entry
	focusLostFunc    func()
	enterPressedFunc func(key *fyne.KeyEvent)
}

func NewDateEntry() *DateEntry {
	entry := &DateEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *DateEntry) SetOnFocusLost(f func()) {
	e.focusLostFunc = f
}

func (e *DateEntry) FocusLost() {
	e.focusLostFunc()
	e.Entry.FocusLost()
}

func (e *DateEntry) SetOnEnter(f func(key *fyne.KeyEvent)) {
	e.enterPressedFunc = f
}

func (e *DateEntry) TypedKey(key *fyne.KeyEvent) {
	e.enterPressedFunc(key)
}
