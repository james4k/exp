package main

import (
	"image"
	"log"

	"j4k.co/exp/ui"
	"j4k.co/exp/ui/glfwenv"
)

type dragbox struct {
}

func (d *dragbox) Mouse(ctl *ui.Control, m, prevm ui.Mouse) {
	if m.Left {
		// TODO: update state.. have something kind of like :hover and
		// such for the style system
		delta := m.Point.Sub(prevm.Point)
		rect := ctl.Rect().Add(delta)
		ctl.SetRect(rect)
	}
}

// Visual - Core visual element in a scene. Delegates events to next
// responder/listener. Much like responder chains in iOS. Maybe. Doesn't
// feel right, tho. :(

// Controller - Handles events, communicates with other controllers and
// its visual and subvisuals. Responsible for creation of its
// visuals.

// Controls.
// So, we only want to manage enough stuff outside the scene graph to
// handle input properly and communicate with other controls...nothing
// visual, including animation, ideally. Only defining of visual style
// indirectly by names or whatever.

func setup() {
	env, err := glfwenv.New(500, 500, "dragbox")
	w, err := ui.Open(env)
	if err != nil {
		log.Fatalln(err)
	}

	w.Register(&dragbox{})

	/*
		w.DefineVisuals(`{
			"body": [
				"quad(#dd3030)",
			],
			"dragbox": {
				"mesh": "dragbox_quad",
			}
		}`)

		w.DefineShapes(`{
			"body_quad": {
				"material": "flat_unlit",
				"shape": "quad",
				"albedo": "rgb(200, 40, 40)"
			},
			"dragbox_quad": {
				"material": "flat_unlit",
				"shape": "quad",
				"albedo": "rgb(255, 255, 255)"
			}
		}`)
	*/

	const doc = `
<box>
{{loop 4}}
<dragbox />
{{end}}
</box>
<trashbin />
`
	//box := w.Insert(nil)
	for i := 0; i < 4; i++ {
		drag := w.Insert(w.Body(), &dragbox{})
		min := image.Pt(20+i*120, 20)
		size := image.Pt(100, 100)
		rect := image.Rectangle{Min: min, Max: min.Add(size)}
		drag.SetRect(rect)
		//layout.Link(w, drag, box)
	}
	//w.Insert(&trashbin{})

	/*
		// note: could be serialized in from disk, network, stdin, etc.
		// supporting quick workflows and inspection/debugging opportunites.
		err = json.Unmarshal([]byte(`{
			"logic": "dragbox",
			"rect": "0 0 100 100",
		}`), w)
		if err != nil {
			log.Fatalln(err)
		}
	*/
}

func main() {
	go setup()
	glfwenv.Run()
}
