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
		if t.State != Active {
			t.State = Hot
		}
	case ui.MouseLeave:
		if t.State != Active {
			t.State = Cold
		}
	case ui.MouseUpdate:
		t.mouseSelect(e)
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

func (t *TextField) mouseSelect(m ui.MouseUpdate) {
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
	var c0 int
	shift := k.Shift()
	if shift {
		c0 = t.Caret[0]
	}
	switch k.Trim() {
	case ui.Backspace:
		t.backspace()
	case ui.Delete:
		t.forwardDelete()
	case ui.Left:
		t.moveBackward(shift)
	case ui.Right:
		t.moveForward(shift)
	default:
		t.emacs(k)
	}
	if shift {
		// keep c0 rooted when holding shift
		t.Caret[0] = c0
	}
}

func (t *TextField) emacs(key ui.Key) {
	switch key.TrimShift() {
	case "^e":
		t.Caret[0] = len(t.Text)
		t.Caret[1] = len(t.Text)
	case "^a":
		t.Caret[0] = 0
		t.Caret[1] = 0
	case "^b":
		t.moveBackward(key.Shift())
	case "^f":
		t.moveForward(key.Shift())
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

func (t *TextField) normalizeCaret() {
	if t.Caret[1] < t.Caret[0] {
		t.Caret[0], t.Caret[1] = t.Caret[1], t.Caret[0]
	}
}

func (t *TextField) moveBackward(shift bool) {
	if !shift && t.Caret[0] != t.Caret[1] {
		t.normalizeCaret()
		t.Caret[1] = t.Caret[0]
	} else if t.Caret[1] > 0 {
		t.Caret[1]--
		t.Caret[0] = t.Caret[1]
	}
}

func (t *TextField) moveForward(shift bool) {
	if !shift && t.Caret[0] != t.Caret[1] {
		t.normalizeCaret()
		t.Caret[0] = t.Caret[1]
	} else if t.Caret[1] < len(t.Text) {
		t.Caret[1]++
		t.Caret[0] = t.Caret[1]
	}
}

type NumberField struct {
	text TextField
}

func (n *NumberField) Receive(ctl *ui.Controller, event interface{}) {
	n.text.Receive(ctl, event)
}
