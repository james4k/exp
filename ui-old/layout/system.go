package layout

import "j4k.co/exp/ui/scene"

type System struct {
}

func (s *System) Run(ctx *scene.Context) {
	ctx.Stage(scene.LayoutStage)
	ctx.Declare(meshCol())
	var (
		meshes MeshColumn
	)
	ctx.Require(&meshes)

	for ctx.Step() {
	}
}
