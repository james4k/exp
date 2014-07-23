package ui

import "image"

func dispatch(target *Box, event interface{}) {
	if target.ctl == nil {
		return
	}
	done := make(chan message)
	target.ctl.recv <- message{
		event: event,
		done:  done,
	}
	<-done
}

// master component is generally responsible for doing layout and drawing
// to screen. sized to the Environment.

func DispatchEvents(env Environment, master Component) error {
	var (
		root       *Box
		mouseFocus *Box
	)
	{
		w, h, _ := env.Size()
		setup(nil, image.Rect(0, 0, w, h), master)
		root = master.box()
	}
	for {
		event, ok := env.Listen()
		if !ok {
			root.close()
			return nil
		}
		switch e := event.(type) {
		case MouseUpdate:
			// TODO: when left mouse button is held down, do not change
			// target/mouseFocus (also means no hitTest needed)
			target := mouseFocus
			if !e.Left {
				target = root.hitTest(e.Point)
			}
			if target != mouseFocus {
				if mouseFocus != nil {
					dispatch(mouseFocus, e)
					dispatch(mouseFocus, MouseLeave{})
				}
				if target != nil {
					dispatch(target, MouseEnter{})
				}
				mouseFocus = target
			}
			if target != nil {
				dispatch(target, e)
			}
			dispatch(root, e)
		case SizeUpdate:
			root.bounds = image.Rect(0, 0, e.Width, e.Height)
			dispatch(root, e)
		}
	}
	return nil
}
