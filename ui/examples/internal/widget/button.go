package widget

import (
	"j4k.co/exp/ui"
)

type Label struct {
	ui.Box
	Text string
}

type State uint8

const (
	Cold State = iota
	Hot
	Active
	Frozen
)

type Button struct {
	ui.Box
	Text    string
	State   State
	OnClick func()
}

func (b *Button) Receive(ctl *ui.Controller, event interface{}) {
	switch e := event.(type) {
	case ui.MouseEnter:
		b.State = Hot
	case ui.MouseLeave:
		b.State = Cold
	case ui.MouseUpdate:
		if b.State == Active {
			b.pressed(e)
		} else {
			b.hover(e)
		}
	}
}

func (b *Button) hover(m ui.MouseUpdate) {
	if m.Left {
		b.State = Active
	}
}

func (b *Button) pressed(m ui.MouseUpdate) {
	if !m.Left {
		if m.Point.In(b.Bounds()) {
			if b.OnClick != nil {
				b.OnClick()
			}
		}
		b.State = Hot
	}
}
