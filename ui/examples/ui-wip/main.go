package main

import (
	"fmt"
	"image"
	"log"

	"j4k.co/exp/ui"
	"j4k.co/exp/ui/examples/internal/blendish"
	"j4k.co/exp/ui/glfwui"
)

// TODO: make ui-todo and other 'real' apps .. should expose some holes

// thoughts are occuring on doing a copy-on-write structure.. and
// changing the event model to allow for real parallelism...but yeah.
// may be a pipe dream.

func main() {
	// run our ui window
	go window()

	// listen for events on the main OS thread
	err := glfwui.ListenForEvents()
	if err != nil {
		log.Fatalln(err)
	}
}

func window() {
	wnd, err := glfwui.Open(500, 500, "ui-wip")
	if err != nil {
		log.Fatalln(err)
	}
	err = ui.DispatchEvents(wnd, &app{wnd: wnd})
	if err != nil {
		log.Fatalln(err)
	}
}

type app struct {
	ui.Box
	wnd *glfwui.Window
}

func (a *app) Run(ctl *ui.Controller) {
	ctl.Append(&clickCounter{})

	a.wnd.MakeContextCurrent()
	defer a.wnd.DetachContext()

	blendish.Init()
	draw(a.wnd, a)
	a.wnd.SwapBuffers()
	for {
		_, ok := ctl.Listen()
		if !ok {
			return
		}
		draw(a.wnd, a)
		a.wnd.SwapBuffers()
	}
}

type clickCounter struct {
	ui.Box
	btn    button
	lbl    label
	clicks int
}

func (c *clickCounter) increment() {
	c.clicks++
	c.lbl.text = fmt.Sprintf("%d click(s)", c.clicks)
}

func (c *clickCounter) Run(ctl *ui.Controller) {
	c.btn = button{text: "Click", onclick: c.increment}
	c.lbl = label{text: "0 click(s)"}
	ctl.Append(&c.btn, &c.lbl)
	{
		x, y := 10, 10
		c.btn.SetBounds(image.Rect(x, y, x+80, y+blendish.WidgetHeight))
		x += 80
		c.lbl.SetBounds(image.Rect(x, y, x+80, y+blendish.WidgetHeight))
	}
	for {
		_, ok := ctl.Listen()
		if !ok {
			return
		}
	}
}

type label struct {
	ui.Box
	text string
}

type buttonState uint8

const (
	buttonDefault buttonState = iota
	buttonHover
	buttonPressed
)

type clicky struct {
}

type button struct {
	ui.Box
	text    string
	state   buttonState
	onclick func()
}

func (b *button) Run(ctl *ui.Controller) {
	for {
		event, ok := ctl.Listen()
		if !ok {
			return
		}
		if _, ok = event.(ui.MouseEnter); ok {
			if !b.hover(ctl) {
				return
			}
			b.state = buttonDefault
		}
	}
}

func (b *button) hover(ctl *ui.Controller) bool {
	b.state = buttonHover
	for {
		event, ok := ctl.Listen()
		if !ok {
			return false
		}
		if _, ok = event.(ui.MouseLeave); ok {
			return true
		}
		if m, ok := event.(ui.MouseUpdate); ok {
			if m.Left {
				if !b.pressed(ctl) {
					return false
				}
				b.state = buttonHover
			}
		}
	}
}

func (b *button) pressed(ctl *ui.Controller) bool {
	b.state = buttonPressed
	for {
		event, ok := ctl.Listen()
		if !ok {
			return false
		}
		if m, ok := event.(ui.MouseUpdate); ok {
			if !m.Left {
				if m.Point.In(b.Bounds()) {
					if b.onclick != nil {
						b.onclick()
					}
				}
				return true
			}
		}
	}
}
