package main

import (
	"fmt"
	"image"
	"log"

	"j4k.co/exp/ui"
	"j4k.co/exp/ui/examples/internal/blendish"
	"j4k.co/exp/ui/examples/internal/widget"
	"j4k.co/exp/ui/glfwui"
)

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
	err = ui.Dispatch(wnd, &app{wnd: wnd})
	if err != nil {
		log.Fatalln(err)
	}
}

type app struct {
	ui.Box
	wnd *glfwui.Window

	drawc chan bool
	donec chan bool
}

func (a *app) Receive(ctl *ui.Controller, event interface{}) {
	switch event.(type) {
	case ui.Mount:
		a.mount(ctl)
		a.drawAndSwap()
	case ui.Unmount:
		break
	case ui.SizeUpdate:
		a.drawAndSwap()
	default:
		a.draw()
	}
}

func (a *app) mount(ctl *ui.Controller) {
	cc := clickCounter{}
	tf := widget.TextField{Text: "Hmm", Caret: [2]int{1, 3}}
	ctl.Mount(&cc, &tf)
	{
		x, y := 10, 10
		cc.SetBounds(image.Rect(x, y, x+150, y+blendish.WidgetHeight))
		y += blendish.WidgetHeight + 4
		tf.SetBounds(image.Rect(x, y, x+150, y+blendish.WidgetHeight))
	}
	a.drawc = make(chan bool, 1)
	a.donec = make(chan bool)
	go render(a.wnd, a)
}

func (a *app) draw() {
	a.drawc <- false
	<-a.donec
}

func (a *app) drawAndSwap() {
	a.drawc <- true
	<-a.donec
}

func render(wnd *glfwui.Window, app *app) {
	wnd.MakeContextCurrent()
	defer wnd.DetachContext()
	blendish.Init()
	for syncswap := range app.drawc {
		draw(wnd, app)
		if !syncswap {
			app.donec <- true
		}
		wnd.SwapBuffers()
		if syncswap {
			app.donec <- true
		}
	}
}

type clickCounter struct {
	ui.Box
	button widget.Button
	label  widget.Label
	clicks int
}

func (c *clickCounter) increment() {
	c.clicks++
	c.label.Text = fmt.Sprintf("%d click(s)", c.clicks)
}

func (c *clickCounter) Receive(ctl *ui.Controller, event interface{}) {
	switch event.(type) {
	case ui.Mount:
		c.mount(ctl)
	}
}

func (c *clickCounter) mount(ctl *ui.Controller) {
	// FIXME: callbacks for UI 'actions' break the future potential for
	// parallel execution of siblings, since siblings would be able to
	// call methods acting on the same state (eg. the parent's state).
	// The controller should schedule these callbacks... Might look like
	// ui.ActionFunc(c.increment), etc.. Must not be callable by user.
	c.button = widget.Button{Text: "Button", OnClick: c.increment}
	c.label = widget.Label{Text: "0 click(s)"}
	ctl.Mount(&c.button, &c.label)
	{
		x, y := 10, 10
		c.button.SetBounds(image.Rect(x, y, x+80, y+blendish.WidgetHeight))
		x += 80
		c.label.SetBounds(image.Rect(x, y, x+80, y+blendish.WidgetHeight))
	}
}
