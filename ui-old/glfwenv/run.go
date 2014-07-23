package glfwenv

import (
	"log"
	"runtime"

	glfw "github.com/go-gl/glfw3"
)

var mainc = make(chan func(), 4)

func init() {
	runtime.LockOSThread()
}

// Run should be called by your main function, and will block
// indefinitely. This is necessary because some platforms expect calls
// from the main thread, thus so does GLFW.
func Run() {
	// ensure we have at least 2 procs, due to the thread conditions we
	// have to work with.
	procs := runtime.GOMAXPROCS(0)
	if procs < 2 {
		runtime.GOMAXPROCS(2)
	}

	glfw.Init()
	defer glfw.Terminate()
	setupGlfw()

	for {
		select {
		case fn := <-mainc:
			fn()
		default:
			glfw.WaitEvents()
		}
	}
}

func setupGlfw() {
	glfw.SetErrorCallback(onGlfwError)

	glfw.WindowHint(glfw.Resizable, 1)
	glfw.WindowHint(glfw.Visible, 1)
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglApi)

	glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglDebugContext, 1)
}

func onGlfwError(code glfw.ErrorCode, desc string) {
	log.Println("glfw: " + desc)
}
