package ui_test

import (
	"fmt"
	"runtime"
	"testing"

	"j4k.co/ui"
)

type env struct {
	list []interface{}
}

func testEnv(events []interface{}) ui.Environment {
	return &env{
		list: events,
	}
}

func (e *env) Size() (w, h int, pixelRatio float32) {
	return 100, 100, 1
}

func (e *env) Listen() (event interface{}, ok bool) {
	if len(e.list) == 0 {
		return nil, false
	}
	event = e.list[0]
	ok = true
	e.list = e.list[1:]
	return
}

type eventChecker struct {
	ui.Box
	t      *testing.T
	events []interface{}
	done   chan<- error
}

func (e *eventChecker) Run(ctl *ui.Controller) {
	for {
		event, ok := ctl.Listen()
		if ok && len(e.events) == 0 {
			e.fatalf("expected end of listener")
		}
		if !ok && len(e.events) > 0 {
			e.fatalf("unexpected end of listener")
		}
		if !ok {
			break
		}
		if event != e.events[0] {
			e.fatalf("expected event %#v, got %#v\n", e.events[0], event)
		}
		e.events = e.events[1:]
	}
	if len(e.events) > 0 {
		e.fatalf("expeced more events")
	}
	e.done <- nil
}

func (e *eventChecker) fatalf(format string, args ...interface{}) {
	e.done <- fmt.Errorf(format, args...)
	runtime.Goexit()
}

func TestDispatchEvents(t *testing.T) {
	events := []interface{}{
		ui.SizeUpdate{
			Width: 110, Height: 110,
		},
		ui.SizeUpdate{
			Width: 120, Height: 120,
		},
	}
	env := testEnv(events)
	done := make(chan error, 1)
	go func() {
		done <- ui.DispatchEvents(env, &eventChecker{
			t:      t,
			events: events,
			done:   done,
		})
	}()
	err := <-done
	if err != nil {
		t.Fatal(err)
	}
	err = <-done
	if err != nil {
		t.Fatal(err)
	}
}
