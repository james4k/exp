package widget

import (
	"log"

	"j4k.co/exp/ui"
)

type TextField struct {
	ui.Box
	Text     string
	Caret    [2]int
	State    State
	OnChange func(text string)
	OnEnter  func(text string)
}

func (t *TextField) Receive(ctl *ui.Controller, event interface{}) {
	switch e := event.(type) {
	case ui.MouseEnter:
		t.State = Hot
	case ui.MouseLeave:
		if t.State != Active {
			t.State = Cold
		}
	case ui.MouseUpdate:
		t.updateSelection(e)
	case ui.FocusGained:
		t.State = Active
	case ui.FocusLost:
		t.State = Cold
	case ui.KeyDown:
		t.keyboard(e.Key)
	case ui.KeyRepeat:
		t.keyboard(e.Key)
	case ui.UnicodeTyped:
		t.insert(e.C)
	}
}

func (t *TextField) insert(char rune) {
	c0 := t.Caret[0]
	c1 := t.Caret[1]
	t.Text = t.Text[:c0] + string(char) + t.Text[c1:]
	t.Caret[0] = c0 + 1
	t.Caret[1] = c0 + 1
}

func (t *TextField) updateSelection(m ui.MouseUpdate) {
	// TODO: implement properly with a simple interface to convert from
	// pt to caret position...for now we just select all.
	if m.Left {
		if !m.Previous.Left {
			t.Caret[0] = 0
		}
		t.Caret[1] = len(t.Text)
	}
}

func (t *TextField) keyboard(k ui.Key) {
	switch k {
	case ui.Backspace:
		t.backspace()
	case ui.Delete:
		t.forwardDelete()
	case ui.Left:
		t.moveBackward()
	case ui.Right:
		t.moveForward()
	default:
		t.emacs(k)
	}
}

func (t *TextField) backspace() {
	c0 := t.Caret[0]
	c1 := t.Caret[1]
	if c0 == c1 && c0 > 0 {
		c0 -= 1
	}
	t.Text = t.Text[:c0] + t.Text[c1:]
	t.Caret[0] = c0
	t.Caret[1] = c0
}

func (t *TextField) forwardDelete() {
	log.Println("TextField.forwardDelete not implemented")
}

func (t *TextField) moveForward() {
	if t.Caret[0] < len(t.Text) {
		t.Caret[0]++
		t.Caret[1] = t.Caret[0]
	}
}

func (t *TextField) moveBackward() {
	if t.Caret[0] > 0 {
		t.Caret[0]--
		t.Caret[1] = t.Caret[0]
	}
}

func (t *TextField) emacs(key ui.Key) {
	switch key {
	case "^e":
		t.Caret[0] = len(t.Text)
		t.Caret[1] = len(t.Text)
	case "^a":
		t.Caret[0] = 0
		t.Caret[1] = 0
	case "^b":
		t.moveBackward()
	case "^f":
		t.moveForward()
	}
}

type NumberField struct {
	text TextField
}

func (n *NumberField) Receive(ctl *ui.Controller, event interface{}) {
	n.text.Receive(ctl, event)
}
