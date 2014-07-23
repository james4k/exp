package main

import (
	"fmt"
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
	wnd, err := glfwui.Open(500, 500, "ui-zoo")
	if err != nil {
		log.Fatalln(err)
	}
	err = ui.DispatchEvents(wnd, &zooApp{wnd: wnd})
	if err != nil {
		log.Fatalln(err)
	}
}

type zooApp struct {
	ui.Box
	wnd *glfwui.Window
}

func (z *zooApp) Run(ctl *ui.Controller) {
	ctl.Append(&clickCounter{})

	z.wnd.MakeContextCurrent()
	defer z.wnd.DetachContext()

	blendish.Init()
	draw(z.wnd, z)
	z.wnd.SwapBuffers()
	for {
		_, ok := ctl.Listen()
		if !ok {
			return
		}
		draw(z.wnd, z)
		z.wnd.SwapBuffers()
	}
}

type clickCounter struct {
	ui.Box
}

func (c *clickCounter) Run(ctl *ui.Controller) {
	var (
		btn    = button{text: "Click", onclick: ctl}
		lbl    = label{text: "0 click(s)"}
		clicks = 0
	)
	ctl.Append(&btn, &lbl)
	for {
		event, ok := ctl.Listen()
		if !ok {
			return
		}
		if _, ok = event.(clicky); ok {
			clicks++
			lbl.text = fmt.Sprintf("%d click(s)", clicks)
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
	onclick interface{}
}

func (b *button) Run(ctl *ui.Controller) {
	for {
		event, ok := ctl.Listen()
		if !ok {
			return
		}
		fmt.Printf("button got event %#v\n", event)
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
		switch m := event.(type) {
		case ui.MouseLeave:
			return true
		case ui.MouseUpdate:
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
					println("click")
					//ctl.Dispatch(b.onclick, clicky{})
				}
				return true
			}
		}
	}
}
