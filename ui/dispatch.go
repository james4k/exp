package ui

import "image"

// master component is generally responsible for doing layout and drawing
// to screen. sized to the Environment.

func Dispatch(env Environment, master Component) error {
	var (
		root       *Box
		keyFocus   *Box
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
			root.unmount()
			return nil
		}
		switch e := event.(type) {
		case KeyDown, KeyUp, KeyRepeat, UnicodeTyped:
			if keyFocus != nil {
				keyFocus.send(e)
			}
		case MouseUpdate:
			target := mouseFocus
			if !e.Left {
				target = root.hitTest(e.Point)
			}
			if target != mouseFocus {
				if mouseFocus != nil {
					mouseFocus.send(e)
					mouseFocus.send(MouseLeave{})
				}
				if target != nil {
					target.send(MouseEnter{})
				}
				mouseFocus = target
			}
			if target != nil {
				target.send(e)
				if e.Left && target != keyFocus {
					if keyFocus != nil {
						keyFocus.send(FocusLost{})
					}
					keyFocus = target
					keyFocus.send(FocusGained{})
				}
			}
		case SizeUpdate:
			root.bounds = image.Rect(0, 0, e.Width, e.Height)
		}
		root.send(event)
	}
	return nil
}
