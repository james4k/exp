package ui

import (
	glfw "github.com/go-gl/glfw3"
	"log"
	"runtime"
)

var mainc chan<- func()

// Run should be called by your main function, and will block
// indefinitely. This is necessary because some platforms expect calls
// from the main thread.
func Run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// ensure we have at least 2 procs, due to the thread conditions we
	// have to work with.
	procs := runtime.GOMAXPROCS(0)
	if procs < 2 {
		runtime.GOMAXPROCS(2)
	}

	glfw.Init()
	defer glfw.Terminate()
	setupGlfw()

	mainchan := make(chan func(), 32)
	mainc = mainchan

	for {
		select {
		case fn := <-mainchan:
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
