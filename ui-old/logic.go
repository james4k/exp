package ui

import "image"

type HitTester interface {
	HitTest(*Control, image.Point) bool
}

type MouseListener interface {
	Mouse(ctl *Control, m, prevm Mouse)
}

type MouseFocusListener interface {
	MouseFocus(*Control, bool)
}

type KeyboardListener interface {
	Keyboard(ctl *Control, k Keyboard, prevk Keyboard)
}

var defaultHitTester = rectHitTester{}

type rectHitTester struct {
}

func (rectHitTester) HitTest(ctl *Control, pt image.Point) bool {
	return pt.In(ctl.Rect())
}
