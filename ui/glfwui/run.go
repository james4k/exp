package glfwui

import (
	"log"
	"runtime"

	glfw "github.com/go-gl/glfw3"
)

var mainc = make(chan func(), 4)

func init() {
	// Lock the main goroutine to the main thread. See the comment in
	// runtime·main() at http://golang.org/src/pkg/runtime/proc.c#L221
	//
	// "Lock the main goroutine onto this, the main OS thread,
	// during initialization.  Most programs won't care, but a few
	// do require certain calls to be made by the main thread.
	// Those can arrange for main.main to run in the main thread
	// by calling runtime.LockOSThread during initialization
	// to preserve the lock."
	//
	// This is needed because on some platforms, we need the main
	// thread, or at least a consistent thread, to make OS or window
	// system calls.
	runtime.LockOSThread()
}

// ListenForEvents should be called by your main function, and will
// block indefinitely. This is necessary because some platforms expect
// calls from the main thread.
func ListenForEvents() error {
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
		case fn, ok := <-mainc:
			if !ok {
				return nil
			}
			fn()
		default:
			glfw.WaitEvents()
		}
	}
	return nil
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
