package main

import (
	"j4k.co/ui"
	bnd "j4k.co/ui/examples/internal/blendish"
	"j4k.co/ui/glfwui"
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
	case *label:
		bnd.Label(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			-1, v.text)
	case *button:
		bnd.Button(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(),
			0, widgetState(v.state), "Default")
	}
	for i := 0; i < view.Subviews(); i++ {
		drawView(view.Sub(i))
	}
}

func widgetState(state buttonState) bnd.WidgetState {
	switch state {
	case buttonHover:
		return bnd.Hover
	case buttonPressed:
		return bnd.Active
	}
	return bnd.Default
}
