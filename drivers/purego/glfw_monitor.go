package purego

import (
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego"
)

var monitorProto atomic.Pointer[Monitor]

type VidMode struct {
	Width       int // The width, in pixels, of the video mode.
	Height      int // The height, in pixels, of the video mode.
	RedBits     int // The bit depth of the red channel of the video mode.
	GreenBits   int // The bit depth of the green channel of the video mode.
	BlueBits    int // The bit depth of the blue channel of the video mode.
	RefreshRate int // The refresh rate, in Hz, of the video mode.
}

// Monitor represents a GLFW monitor handle
type Monitor struct {
	handle uintptr

	glfwGetPrimaryMonitor  uintptr
	glfwGetVideoMode       uintptr
	glfwGetMonitorWorkarea uintptr
}

// Handle returns the raw uintptr handle
func (m *Monitor) Handle() uintptr {
	if m == nil {
		return 0
	}
	return m.handle
}

func (m *Monitor) GetVideoMode() *VidMode {
	vmPtr, _, _ := purego.SyscallN(m.glfwGetVideoMode, m.handle)
	if vmPtr == 0 {
		return nil
	}

	// The C struct GLFWvidmode has this layout:
	// typedef struct GLFWvidmode {
	//     int width;
	//     int height;
	//     int redBits;
	//     int greenBits;
	//     int blueBits;
	//     int refreshRate;
	// } GLFWvidmode;

	vm := &VidMode{
		Width:       int(*(*int32)(unsafe.Pointer(vmPtr + 0))),
		Height:      int(*(*int32)(unsafe.Pointer(vmPtr + 4))),
		RedBits:     int(*(*int32)(unsafe.Pointer(vmPtr + 8))),
		GreenBits:   int(*(*int32)(unsafe.Pointer(vmPtr + 12))),
		BlueBits:    int(*(*int32)(unsafe.Pointer(vmPtr + 16))),
		RefreshRate: int(*(*int32)(unsafe.Pointer(vmPtr + 20))),
	}

	return vm
}

func (m *Monitor) GetWorkarea() (int, int, int, int) {
	var xpos, ypos, width, height int
	purego.SyscallN(
		m.glfwGetMonitorWorkarea,
		m.handle,
		uintptr(unsafe.Pointer(&xpos)),
		uintptr(unsafe.Pointer(&ypos)),
		uintptr(unsafe.Pointer(&width)),
		uintptr(unsafe.Pointer(&height)),
	)
	return xpos, ypos, width, height
}
