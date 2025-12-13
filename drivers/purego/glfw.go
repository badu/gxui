package purego

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// GLFW constants
const (
	GLFW_TRUE  = 1
	GLFW_FALSE = 0
)

const (
	GLFW_CURSOR           = 0x00033001
	GLFW_CURSOR_DISABLED  = 0x00034003
	GLFW_STICKY_KEYS      = 0x00033002
	GLFW_RAW_MOUSE_MOTION = 0x00033005
)

const (
	GLFW_PRESS               = 1
	GLFW_MOUSE_BUTTON_LEFT   = 0
	GLFW_MOUSE_BUTTON_RIGHT  = 1
	GLFW_MOUSE_BUTTON_MIDDLE = 2
)

// Window represents a GLFW window handle

// GLFW function pointers
var (
	glfwInit               uintptr
	glfwTerminate          uintptr
	glfwWindowHint         uintptr
	glfwDefaultWindowHints uintptr
	glfwPollEvents         uintptr
	glfwWaitEvents         uintptr
	glfwPostEmptyEvent     uintptr
	glfwSwapInterval       uintptr

	glfwGetProcAddress uintptr
	glfwGetTime        uintptr
)

// Init initializes the GLFW library
func Init() int {
	ret, _, _ := purego.SyscallN(glfwInit)
	return int(ret)
}

// Terminate destroys all remaining windows and frees resources
func Terminate() {
	purego.SyscallN(glfwTerminate)
}

// WindowHint sets hints for the next window creation
func WindowHint(hint Hint, value int) {
	purego.SyscallN(glfwWindowHint, uintptr(hint), uintptr(value))
}

// DefaultWindowHints resets all window hints to their default values
func DefaultWindowHints() {
	purego.SyscallN(glfwDefaultWindowHints)
}

// GetPrimaryMonitor returns the primary monitor
func GetPrimaryMonitor() *Monitor {
	monitorProtoPtr := monitorProto.Load()
	if monitorProtoPtr == nil {
		fmt.Println("glfwGetPrimaryMonitor failed (monitorProtoPtr is nil)")
		return nil
	}
	clone := *monitorProtoPtr
	ret, _, _ := purego.SyscallN(clone.glfwGetPrimaryMonitor)
	if ret == 0 {
		return nil
	}

	clone.handle = ret

	return &clone
}

// PollEvents processes all pending events
func PollEvents() {
	purego.SyscallN(glfwPollEvents)
}

// WaitEvents waits until events are queued and processes them
func WaitEvents() {
	purego.SyscallN(glfwWaitEvents)
}

// PostEmptyEvent posts an empty event to the event queue
func PostEmptyEvent() {
	purego.SyscallN(glfwPostEmptyEvent)
}

// SwapInterval sets the swap interval for the current context
func SwapInterval(interval int) {
	purego.SyscallN(glfwSwapInterval, uintptr(interval))
}

// GetProcAddress returns the address of the specified OpenGL or OpenGL ES core or extension function
func GetProcAddress(procname string) uintptr {
	nameBytes := append([]byte(procname), 0)
	ret, _, _ := purego.SyscallN(glfwGetProcAddress, uintptr(unsafe.Pointer(&nameBytes[0])))
	return ret
}

// GetTime returns the current GLFW time in seconds
func GetTime() float64 {
	ret, _, _ := purego.SyscallN(glfwGetTime)
	return *(*float64)(unsafe.Pointer(&ret))
}

// cStringToGoString converts a C string (null-terminated) to a Go string
func cStringToGoString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	var length int
	for {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(length)))
		if b == 0 {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length))
}

// LoadGLFW loads the GLFW shared library and resolves function symbols
func LoadGLFW() error {
	// Try common GLFW library names and paths on Linux
	libPaths := []string{
		"libglfw.so.3",
		"libglfw.so",
		"/usr/lib/x86_64-linux-gnu/libglfw.so.3",
		"/usr/lib/x86_64-linux-gnu/libglfw.so",
		"/usr/local/lib/libglfw.so.3",
		"/usr/local/lib/libglfw.so",
		"/usr/lib/libglfw.so.3",
		"/usr/lib/libglfw.so",
	}

	var lib uintptr
	var err error

	for _, path := range libPaths {
		lib, err = purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil {
			fmt.Printf("Successfully loaded GLFW from: %s\n", path)
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to load GLFW library: %w\n\nTroubleshooting:\n"+
			"1. Install GLFW: sudo apt-get install libglfw3 libglfw3-dev\n"+
			"2. Find installed location: sudo find /usr -name 'libglfw.so*' 2>/dev/null\n"+
			"3. Check library path is in LD_LIBRARY_PATH or /etc/ld.so.conf\n"+
			"4. Run: sudo ldconfig", err)
	}

	prototypeWindow := Window{}
	prototypeMonitor := Monitor{}
	// Register function symbols
	funcs := map[string]*uintptr{
		"glfwInit":               &glfwInit,
		"glfwTerminate":          &glfwTerminate,
		"glfwWindowHint":         &glfwWindowHint,
		"glfwDefaultWindowHints": &glfwDefaultWindowHints,
		"glfwPollEvents":         &glfwPollEvents,
		"glfwWaitEvents":         &glfwWaitEvents,
		"glfwPostEmptyEvent":     &glfwPostEmptyEvent,
		"glfwSwapInterval":       &glfwSwapInterval,
		"glfwGetTime":            &glfwGetTime,
		"glfwGetProcAddress":     &glfwGetProcAddress,

		"glfwGetPrimaryMonitor":  &prototypeMonitor.glfwGetPrimaryMonitor,
		"glfwGetMonitorWorkarea": &prototypeMonitor.glfwGetMonitorWorkarea,
		"glfwGetVideoMode":       &prototypeMonitor.glfwGetVideoMode,

		"glfwCreateWindow":         &prototypeWindow.glfwCreateWindow,
		"glfwMakeContextCurrent":   &prototypeWindow.glfwMakeContextCurrent,
		"glfwGetClipboardString":   &prototypeWindow.glfwGetClipboardString,
		"glfwSetClipboardString":   &prototypeWindow.glfwSetClipboardString,
		"glfwGetCursorPos":         &prototypeWindow.glfwGetCursorPos,
		"glfwMaximizeWindow":       &prototypeWindow.glfwMaximizeWindow,
		"glfwIconifyWindow":        &prototypeWindow.glfwIconifyWindow,
		"glfwRestoreWindow":        &prototypeWindow.glfwRestoreWindow,
		"glfwGetInputMode":         &prototypeWindow.glfwGetInputMode,
		"glfwSetInputMode":         &prototypeWindow.glfwSetInputMode,
		"glfwGetWindowPos":         &prototypeWindow.glfwGetWindowPos,
		"glfwSetWindowPos":         &prototypeWindow.glfwSetWindowPos,
		"glfwGetWindowSize":        &prototypeWindow.glfwGetWindowSize,
		"glfwGetFramebufferSize":   &prototypeWindow.glfwGetFramebufferSize,
		"glfwGetKey":               &prototypeWindow.glfwGetKey,
		"glfwGetMouseButton":       &prototypeWindow.glfwGetMouseButton,
		"glfwSetWindowSize":        &prototypeWindow.glfwSetWindowSize,
		"glfwSetWindowTitle":       &prototypeWindow.glfwSetWindowTitle,
		"glfwSwapBuffers":          &prototypeWindow.glfwSwapBuffers,
		"glfwShowWindow":           &prototypeWindow.glfwShowWindow,
		"glfwHideWindow":           &prototypeWindow.glfwHideWindow,
		"glfwDestroyWindow":        &prototypeWindow.glfwDestroyWindow,
		"glfwWindowShouldClose":    &prototypeWindow.glfwWindowShouldClose,
		"glfwSetWindowShouldClose": &prototypeWindow.glfwSetWindowShouldClose,

		"glfwSetCursorPosCallback":       &prototypeWindow.glfwSetCursorPosCallback,
		"glfwSetKeyCallback":             &prototypeWindow.glfwSetKeyCallback,
		"glfwSetCharCallback":            &prototypeWindow.glfwSetCharCallback,
		"glfwSetScrollCallback":          &prototypeWindow.glfwSetScrollCallback,
		"glfwSetMouseButtonCallback":     &prototypeWindow.glfwSetMouseButtonCallback,
		"glfwSetFramebufferSizeCallback": &prototypeWindow.glfwSetFramebufferSizeCallback,
		"glfwSetWindowCloseCallback":     &prototypeWindow.glfwSetWindowCloseCallback,
		"glfwSetWindowRefreshCallback":   &prototypeWindow.glfwSetWindowRefreshCallback,
		"glfwSetWindowSizeCallback":      &prototypeWindow.glfwSetWindowSizeCallback,
		"glfwSetCursorEnterCallback":     &prototypeWindow.glfwSetCursorEnterCallback,
		"glfwSetCharModsCallback":        &prototypeWindow.glfwSetCharModsCallback,
		"glfwSetWindowPosCallback":       &prototypeWindow.glfwSetWindowPosCallback,
		"glfwSetWindowFocusCallback":     &prototypeWindow.glfwSetWindowFocusCallback,
		"glfwSetWindowIconifyCallback":   &prototypeWindow.glfwSetWindowIconifyCallback,
		"glfwSetDropCallback":            &prototypeWindow.glfwSetDropCallback,
	}

	for name, ptr := range funcs {
		sym, err := purego.Dlsym(lib, name)
		if err != nil {
			fmt.Printf("error loading symbol %q : %#v\n", name, err)
			return fmt.Errorf("failed to load symbol %s: %w", name, err)
		}

		if sym == 0 {
			fmt.Printf("no such %q : uintptr is zero\n", name)
			return fmt.Errorf("failed to load symbol %s: no such symbol", name)
		}

		*ptr = sym
	}

	windowProto.Store(&prototypeWindow)
	monitorProto.Store(&prototypeMonitor)

	return nil
}
