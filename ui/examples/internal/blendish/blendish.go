/*
Package blendish provides some Blender-like UI drawing calls for a GL
context. Meant only for demonstration purposes in some UI examples, so
would likely take some work for production use (say, to draw with
another API, or to integrate with your rendering pipeline).

All function calls expected to be on the GL thread.

This is just a thin wrapper around Leonard Ritter's and Mikko Mononen's
fine work:
https://bitbucket.org/duangle/blendish
https://github.com/memononen/nanovg
*/
package blendish

/*
#cgo CFLAGS: -Wno-typedef-redefinition
#cgo darwin LDFLAGS: -framework OpenGL -lGLEW
#cgo linux LDFLAGS: -lGLEW -lGL
#cgo windows LDFLAGS: -lglew32 -lopengl32
#include <GL/glew.h>
#include "nanovg.h"
#define NANOVG_GL3_IMPLEMENTATION
#include "nanovg_gl.h"
#define BLENDISH_IMPLEMENTATION
#include "blendish.h"
*/
import "C"
import (
	"log"
	"os"
	"path/filepath"
	"unsafe"
)

var vg *C.NVGcontext

const (
	WidgetHeight = 21
)

type WidgetState uint8

const (
	Default WidgetState = iota
	Hover
	Active
)

// CornerFlags indicate which corners are sharp (for grouping widgets)
type CornerFlags uint8

const (
	CornerNone CornerFlags = 0

	CornerTopLeft   = 0x1
	CornerTopRight  = 0x2
	CornerDownLeft  = 0x4
	CornerDownRight = 0x8

	CornerAll   = 0xf
	CornerTop   = 0x3
	CornerDown  = 0xc
	CornerLeft  = 0x9
	CornerRight = 0x6
)

var dirs = filepath.SplitList(os.Getenv("GOPATH"))

func init() {
	const assets = "src/j4k.co/exp/ui/examples/internal/blendish"
	for i := range dirs {
		dirs[i] = filepath.Join(dirs[i], assets)
		break // only care about first
	}
}

func Init() {
	if vg != nil {
		panic("blendish: already initialized")
	}
	C.glewExperimental = C.GL_TRUE
	if C.glewInit() != C.GLEW_OK {
		log.Fatalln("blendish: could not init glew (GL extensions)")
	}
	// GLEW generates GL error because it calls
	// glGetString(GL_EXTENSIONS), we'll consume it here.
	C.glGetError()
	vg = C.nvgCreateGL3(C.NVG_ANTIALIAS | C.NVG_STENCIL_STROKES)

	fontPath := C.CString(filepath.Join(dirs[0], "DejaVuSans.ttf"))
	iconsPath := C.CString(filepath.Join(dirs[0], "blender_icons16.png"))
	sys := C.CString("system")
	defer C.free(unsafe.Pointer(fontPath))
	defer C.free(unsafe.Pointer(iconsPath))
	defer C.free(unsafe.Pointer(sys))
	C.bndSetFont(C.nvgCreateFont(vg, sys, fontPath))
	C.bndSetIconImage(C.nvgCreateImage(vg, iconsPath, 0))
}

func BeginFrame(width, height int, devicePixelRatio float32) {
	w := C.GLsizei(float32(width) * devicePixelRatio)
	h := C.GLsizei(float32(height) * devicePixelRatio)
	C.glViewport(0, 0, w, h)
	C.glClearColor(0, 0, 0, 1)
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT | C.GL_STENCIL_BUFFER_BIT)
	C.nvgBeginFrame(vg, C.int(width), C.int(height), C.float(devicePixelRatio))
}

func EndFrame() {
	C.nvgEndFrame(vg)
}

func Background(x, y, w, h int) {
	C.bndBackground(vg, C.float(x), C.float(y), C.float(w), C.float(h))
}

func Label(x, y, w, h int, iconID int, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.bndLabel(vg, C.float(x), C.float(y), C.float(w), C.float(h),
		C.int(iconID), clabel)
}

func Button(x, y, w, h int, corners CornerFlags, state WidgetState, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.bndToolButton(vg, C.float(x), C.float(y), C.float(w), C.float(h),
		C.int(corners), C.BNDwidgetState(state), -1, clabel)
}
