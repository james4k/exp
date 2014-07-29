package glfwui

import (
	"log"
	"runtime"
	"time"

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

	t0 := time.Now()
	for {
		select {
		case fn, ok := <-mainc:
			if !ok {
				return nil
			}
			fn()
		default:
			// when under heavy activity, sleep and poll instead to
			// lessen sysmon() churn in Go's scheduler. with this
			// method, OS X's Activity Monitor reports over 1000 sleeps
			// a second, which is ridiculous, but better than 3000.
			// haven't been able to create a minimal reproduction of
			// this, but when tweaking values in runtime/proc.c's sysmon
			// func, it is obvious this is a scheduler issue.
			// TODO: best workaround is probably to write this entire
			// loop in C, which should keep the sysmon thread from
			// waking up (from these syscalls, anyways).
			dt := time.Now().Sub(t0)
			if dt < 150*time.Millisecond {
				time.Sleep(15 * time.Millisecond)
				glfw.PollEvents()
			} else {
				t0 = time.Now()
				glfw.WaitEvents()
			}
		}
	}
	return nil
}

func setupGlfw() {
	glfw.SetErrorCallback(onGlfwError)

	glfw.WindowHint(glfw.Resizable, 1)
	glfw.WindowHint(glfw.Visible, 1)
	glfw.WindowHint(glfw.Decorated, 1)
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
