package main

import (
	"j4k.co/exp/ui"
	bnd "j4k.co/exp/ui/examples/internal/blendish"
	"j4k.co/exp/ui/examples/internal/widget"
	"j4k.co/exp/ui/glfwui"
)

func draw(wnd *glfwui.Window, body ui.View) {
	w, h, ratio := wnd.Size()
	bnd.BeginFrame(w, h, ratio)
	bnd.Background(0, 0, w, h)
	drawView(body)
	bnd.EndFrame()
}

func drawView(view ui.View) {
	rect := view.Bounds()
	switch v := view.(type) {
	case *widget.Label:
		bnd.Label(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			-1, v.Text)
	case *widget.Button:
		bnd.Button(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			0, widgetState(v.State), v.Text)
	case *widget.TextField:
		// TODO: calculate the amount of text that will fit in the
		// dimensions, account for the caret position
		bnd.TextField(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			0, widgetState(v.State), v.Text, v.Caret[0], v.Caret[1])
	}
	for i := 0; i < view.Subviews(); i++ {
		drawView(view.Sub(i))
	}
}

func widgetState(state widget.State) bnd.WidgetState {
	switch state {
	case widget.Hot:
		return bnd.Hover
	case widget.Active:
		return bnd.Active
	}
	return bnd.Default
}
