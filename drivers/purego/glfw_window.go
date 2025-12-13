package purego

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego"
)

var windowProto atomic.Pointer[Window]

type Window struct {
	handle                   uintptr
	glfwCreateWindow         uintptr
	glfwMakeContextCurrent   uintptr
	glfwGetClipboardString   uintptr
	glfwSetClipboardString   uintptr
	glfwGetCursorPos         uintptr
	glfwMaximizeWindow       uintptr
	glfwIconifyWindow        uintptr
	glfwRestoreWindow        uintptr
	glfwGetKey               uintptr
	glfwGetMouseButton       uintptr
	glfwGetWindowSize        uintptr
	glfwGetInputMode         uintptr
	glfwSetInputMode         uintptr
	glfwGetWindowPos         uintptr
	glfwSetWindowPos         uintptr
	glfwSetWindowSize        uintptr
	glfwSetWindowTitle       uintptr
	glfwHideWindow           uintptr
	glfwDestroyWindow        uintptr
	glfwWindowShouldClose    uintptr
	glfwSetWindowShouldClose uintptr
	glfwShowWindow           uintptr
	glfwSwapBuffers          uintptr
	glfwGetFramebufferSize   uintptr

	glfwSetCursorPosCallback       uintptr
	glfwSetKeyCallback             uintptr
	glfwSetCharCallback            uintptr
	glfwSetScrollCallback          uintptr
	glfwSetMouseButtonCallback     uintptr
	glfwSetFramebufferSizeCallback uintptr
	glfwSetWindowCloseCallback     uintptr
	glfwSetWindowRefreshCallback   uintptr
	glfwSetWindowSizeCallback      uintptr
	glfwSetCursorEnterCallback     uintptr
	glfwSetCharModsCallback        uintptr
	glfwSetWindowPosCallback       uintptr
	glfwSetWindowFocusCallback     uintptr
	glfwSetWindowIconifyCallback   uintptr
	glfwSetDropCallback            uintptr

	cursorPosCallback       CursorPosCallback
	keyCallback             KeyCallback
	charCallback            CharCallback
	scrollCallback          ScrollCallback
	mouseButtonCallback     MouseButtonCallback
	framebufferSizeCallback FramebufferSizeCallback
	closeCallback           CloseCallback
	refreshCallback         RefreshCallback
	sizeCallback            SizeCallback
	cursorEnterCallback     CursorEnterCallback
	charModsCallback        CharModsCallback
	posCallback             PosCallback
	focusCallback           FocusCallback
	iconifyCallback         IconifyCallback
	dropCallback            DropCallback
}

// Handle returns the raw uintptr handle
func (w *Window) Handle() uintptr {
	if w == nil {
		return 0
	}
	return w.handle
}

// MakeContextCurrent makes this window's context current
func (w *Window) MakeContextCurrent() {
	purego.SyscallN(w.glfwMakeContextCurrent, w.handle)
}

// TODO : DetachCurrentContext was not tested

// DetachCurrentContext detaches the current context from the current thread
// This is done by calling MakeContextCurrent with NULL (0)
func (w *Window) DetachCurrentContext() {
	purego.SyscallN(w.glfwMakeContextCurrent, uintptr(0))
}

// GetCursorPos retrieves the cursor position for this window
func (w *Window) GetCursorPos() (float64, float64) {
	var xpos, ypos float64
	if w.glfwGetCursorPos == 0 {
		fmt.Println("glfwGetCursorPos not found")
		return 0, 0
	}
	purego.SyscallN(w.glfwGetCursorPos, w.handle, uintptr(unsafe.Pointer(&xpos)), uintptr(unsafe.Pointer(&ypos)))
	return xpos, ypos
}

// SetCursorPosCallback sets the cursor position callback for this window
func (w *Window) SetCursorPosCallback(callback CursorPosCallback) {
	if callback != nil {
		w.cursorPosCallback = callback
		callbackPtr := purego.NewCallback(
			func(handler uintptr, xpos float64, ypos float64) {
				if w.handle == handler {
					callback(w, xpos, ypos)
				}
			},
		)
		purego.SyscallN(w.glfwSetCursorPosCallback, w.handle, callbackPtr)
	} else {
		w.cursorPosCallback = nil
		purego.SyscallN(w.glfwSetCursorPosCallback, w.handle, 0)
	}
}

// SetKeyCallback sets the key callback for this window
func (w *Window) SetKeyCallback(callback KeyCallback) {
	if callback != nil {
		w.keyCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, key Key, scancode int32, action Action, mods ModifierKey) {
			if w.handle == handler {
				callback(w, key, scancode, action, mods)
			}
		})
		purego.SyscallN(w.glfwSetKeyCallback, w.handle, callbackPtr)
	} else {
		w.keyCallback = nil
		purego.SyscallN(w.glfwSetKeyCallback, w.handle, 0)
	}
}

// SetCharCallback sets the character callback for this window
func (w *Window) SetCharCallback(callback CharCallback) {
	if callback != nil {
		w.charCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, char rune) {
			if w.handle != handler {
				callback(w, char)
			}
		})
		purego.SyscallN(w.glfwSetCharCallback, w.handle, callbackPtr)
	} else {
		w.charCallback = nil
		purego.SyscallN(w.glfwSetCharCallback, w.handle, 0)
	}
}

// SetScrollCallback sets the scroll callback for this window
func (w *Window) SetScrollCallback(callback ScrollCallback) {
	if callback != nil {
		w.scrollCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, xoffset float64, yoffset float64) {
			if w.handle == handler {
				callback(w, xoffset, yoffset)
			}
		})
		purego.SyscallN(w.glfwSetScrollCallback, w.handle, callbackPtr)
	} else {
		w.scrollCallback = nil
		purego.SyscallN(w.glfwSetScrollCallback, w.handle, 0)
	}
}

// SetMouseButtonCallback sets the mouse button callback for this window
func (w *Window) SetMouseButtonCallback(callback MouseButtonCallback) {
	if callback != nil {
		w.mouseButtonCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, button MouseButton, action Action, mods ModifierKey) {
			if w.handle == handler {
				callback(w, button, action, mods)
			}
		})
		purego.SyscallN(w.glfwSetMouseButtonCallback, w.handle, callbackPtr)
	} else {
		w.mouseButtonCallback = nil
		purego.SyscallN(w.glfwSetMouseButtonCallback, w.handle, 0)
	}
}

// SetFramebufferSizeCallback sets the framebuffer size callback for this window
func (w *Window) SetFramebufferSizeCallback(callback FramebufferSizeCallback) {
	if callback != nil {
		w.framebufferSizeCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, width int32, height int32) {
			if w.handle == handler {
				callback(w, width, height)
			}
		})
		purego.SyscallN(w.glfwSetFramebufferSizeCallback, w.handle, callbackPtr)
	} else {
		w.framebufferSizeCallback = nil
		purego.SyscallN(w.glfwSetFramebufferSizeCallback, w.handle, 0)
	}
}

// SetCloseCallback sets the window close callback for this window
func (w *Window) SetCloseCallback(callback CloseCallback) {
	if callback != nil {
		w.closeCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr) {
			if w.handle == handler {
				callback(w)
			}
		})
		purego.SyscallN(w.glfwSetWindowCloseCallback, w.handle, callbackPtr)
	} else {
		w.closeCallback = nil
		purego.SyscallN(w.glfwSetWindowCloseCallback, w.handle, 0)
	}
}

// SetRefreshCallback sets the window refresh callback for this window
func (w *Window) SetRefreshCallback(callback RefreshCallback) {
	if callback != nil {
		w.refreshCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr) {
			if w.handle == handler {
				callback(w)
			}
		})
		purego.SyscallN(w.glfwSetWindowRefreshCallback, w.handle, callbackPtr)
	} else {
		w.refreshCallback = nil
		purego.SyscallN(w.glfwSetWindowRefreshCallback, w.handle, 0)
	}
}

// SetSizeCallback sets the window size callback for this window
func (w *Window) SetSizeCallback(callback SizeCallback) {
	if callback != nil {
		w.sizeCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, width int32, height int32) {
			if w.handle == handler {
				callback(w, width, height)
			}
		})
		purego.SyscallN(w.glfwSetWindowSizeCallback, w.handle, callbackPtr)
	} else {
		w.sizeCallback = nil
		purego.SyscallN(w.glfwSetWindowSizeCallback, w.handle, 0)
	}
}

// SetCursorEnterCallback sets the cursor enter/leave callback for this window
func (w *Window) SetCursorEnterCallback(callback CursorEnterCallback) {
	if callback != nil {
		w.cursorEnterCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, entered int) {
			if w.handle == handler {
				callback(w, entered == GLFW_TRUE)
			}
		})
		purego.SyscallN(w.glfwSetCursorEnterCallback, w.handle, callbackPtr)
	} else {
		w.cursorEnterCallback = nil
		purego.SyscallN(w.glfwSetCursorEnterCallback, w.handle, 0)
	}
}

// SetCharModsCallback sets the character with modifiers callback for this window
func (w *Window) SetCharModsCallback(callback CharModsCallback) {
	if callback != nil {
		w.charModsCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, char rune, mods ModifierKey) {
			if w.handle == handler {
				callback(w, char, mods)
			}
		})
		purego.SyscallN(w.glfwSetCharModsCallback, w.handle, callbackPtr)
	} else {
		w.charModsCallback = nil
		purego.SyscallN(w.glfwSetCharModsCallback, w.handle, 0)
	}
}

// SetPosCallback sets the window position callback for this window
func (w *Window) SetPosCallback(callback PosCallback) {
	if callback != nil {
		w.posCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, xpos int32, ypos int32) {
			if w.handle == handler {
				callback(w, xpos, ypos)
			}
		})
		purego.SyscallN(w.glfwSetWindowPosCallback, w.handle, callbackPtr)
	} else {
		w.posCallback = nil
		purego.SyscallN(w.glfwSetWindowPosCallback, w.handle, 0)
	}
}

// SetFocusCallback sets the window focus callback for this window
func (w *Window) SetFocusCallback(callback FocusCallback) {
	if callback != nil {
		w.focusCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, focused bool) {
			if w.handle == handler {
				callback(w, focused)
			}
		})
		purego.SyscallN(w.glfwSetWindowFocusCallback, w.handle, callbackPtr)
	} else {
		w.focusCallback = nil
		purego.SyscallN(w.glfwSetWindowFocusCallback, w.handle, 0)
	}
}

// SetIconifyCallback sets the window iconify callback for this window
func (w *Window) SetIconifyCallback(callback IconifyCallback) {
	if callback != nil {
		w.iconifyCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, iconified bool) {
			if w.handle == handler {
				callback(w, iconified)
			}
		})
		purego.SyscallN(w.glfwSetWindowIconifyCallback, w.handle, callbackPtr)
	} else {
		w.iconifyCallback = nil
		purego.SyscallN(w.glfwSetWindowIconifyCallback, w.handle, 0)
	}
}

// SetDropCallback sets the file drop callback for this window
func (w *Window) SetDropCallback(callback DropCallback) {
	if callback != nil {
		w.dropCallback = callback
		callbackPtr := purego.NewCallback(func(handler uintptr, count int, pathsPtr uintptr) {
			if w.handle == handler {
				// Convert C string array to Go string slice
				paths := make([]string, count)
				for i := 0; i < count; i++ {
					// Get pointer to the i-th string pointer
					strPtr := *(*uintptr)(unsafe.Pointer(pathsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
					// Convert C string to Go string
					paths[i] = cStringToGoString(strPtr)
				}
				callback(w, paths)
			}
		})
		purego.SyscallN(w.glfwSetDropCallback, w.handle, callbackPtr)
	} else {
		w.dropCallback = nil
		purego.SyscallN(w.glfwSetDropCallback, w.handle, 0)
	}
}

// SetClipboardString sets the clipboard to the specified string for this window
func (w *Window) SetClipboardString(str string) {
	strBytes := append([]byte(str), 0)
	purego.SyscallN(w.glfwSetClipboardString, w.handle, uintptr(unsafe.Pointer(&strBytes[0])))
}

func (w *Window) GetClipboardString() string {
	ret, _, _ := purego.SyscallN(w.glfwGetClipboardString, w.handle)
	return cStringToGoString(ret)
}

func (w *Window) GetKey(key Key) Action {
	ret, _, _ := purego.SyscallN(w.glfwGetKey, w.handle, uintptr(key))
	return Action(ret)
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	if w.glfwGetMouseButton == 0 {
		fmt.Println("glfwGetMouseButton not found")
		return 0
	}
	ret, _, _ := purego.SyscallN(w.glfwGetMouseButton, w.handle, uintptr(button))
	return Action(ret)
}

func (w *Window) GetInputMode(mode InputMode) int32 {
	ret, _, _ := purego.SyscallN(w.glfwGetInputMode, w.handle, uintptr(mode))
	return int32(ret)
}

func (w *Window) SetInputMode(mode InputMode, value int32) {
	purego.SyscallN(w.glfwSetInputMode, w.handle, uintptr(mode), uintptr(value))
}

func (w *Window) GetFramebufferSize() (int32, int32) {
	var width, height int32
	purego.SyscallN(w.glfwGetFramebufferSize, w.handle, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return width, height
}

func (w *Window) GetPos() (int32, int32) {
	var xpos, ypos int32
	purego.SyscallN(w.glfwGetWindowPos, w.handle, uintptr(unsafe.Pointer(&xpos)), uintptr(unsafe.Pointer(&ypos)))
	return xpos, ypos
}

func (w *Window) SetPos(xpos, ypos int) {
	purego.SyscallN(w.glfwSetWindowPos, w.handle, uintptr(xpos), uintptr(ypos))
}

func (w *Window) GetSize() (int, int) {
	var width, height int
	purego.SyscallN(w.glfwGetWindowSize, w.handle, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return width, height
}

func (w *Window) SetSize(width, height int) {
	purego.SyscallN(w.glfwSetWindowSize, w.handle, uintptr(width), uintptr(height))
}

func (w *Window) SetTitle(title string) {
	titleBytes := append([]byte(title), 0)
	purego.SyscallN(w.glfwSetWindowTitle, w.handle, uintptr(unsafe.Pointer(&titleBytes[0])))
}

func (w *Window) SwapBuffers() {
	purego.SyscallN(w.glfwSwapBuffers, w.handle)
}

func (w *Window) Show() {
	purego.SyscallN(w.glfwShowWindow, w.handle)
}

func (w *Window) Hide() {
	purego.SyscallN(w.glfwHideWindow, w.handle)
}

func (w *Window) Destroy() {
	purego.SyscallN(w.glfwDestroyWindow, w.handle)
}

// ShouldClose returns the value of the close flag for this window
func (w *Window) ShouldClose() bool {
	ret, _, _ := purego.SyscallN(w.glfwWindowShouldClose, w.handle)
	return ret == GLFW_TRUE
}

// SetShouldClose sets the value of the close flag for this window
func (w *Window) SetShouldClose(value bool) {
	val := GLFW_FALSE
	if value {
		val = GLFW_TRUE
	}
	purego.SyscallN(w.glfwSetWindowShouldClose, w.handle, uintptr(val))
}

// Maximize maximizes this window
func (w *Window) Maximize() {
	purego.SyscallN(w.glfwMaximizeWindow, w.handle)
}

// Iconify iconifies (minimizes) this window
func (w *Window) Iconify() {
	purego.SyscallN(w.glfwIconifyWindow, w.handle)
}

// Restore restores this window
func (w *Window) Restore() {
	purego.SyscallN(w.glfwRestoreWindow, w.handle)
}

// GetWindowSize retrieves the size of the content area for this window
func (w *Window) GetWindowSize() (int, int) {
	var width, height int
	purego.SyscallN(w.glfwGetWindowSize, w.handle, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return width, height
}

// GetWindowPos retrieves the position of the content area for this window
func (w *Window) GetWindowPos() (int32, int32) {
	var xpos, ypos int32
	purego.SyscallN(w.glfwGetWindowPos, w.handle, uintptr(unsafe.Pointer(&xpos)), uintptr(unsafe.Pointer(&ypos)))
	return xpos, ypos
}

// SetWindowPos sets the position of the content area for this window
func (w *Window) SetWindowPos(xpos, ypos int) {
	purego.SyscallN(w.glfwSetWindowPos, w.handle, uintptr(xpos), uintptr(ypos))
}

// SetWindowSize sets the size of the content area for this window
func (w *Window) SetWindowSize(width, height int) {
	purego.SyscallN(w.glfwSetWindowSize, w.handle, uintptr(width), uintptr(height))
}

// SetWindowTitle sets the title for this window
func (w *Window) SetWindowTitle(title string) {
	titleBytes := append([]byte(title), 0)
	purego.SyscallN(w.glfwSetWindowTitle, w.handle, uintptr(unsafe.Pointer(&titleBytes[0])))
}

// CreateWindow creates a window and its associated context
func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) *Window {
	windowProtoPtr := windowProto.Load()
	if windowProtoPtr == nil {
		fmt.Println("glfwCreateWindow failed (windowProtoPtr is nil)")
		return nil
	}

	titlePtr := append([]byte(title), 0)

	var monitorHandle uintptr
	if monitor != nil {
		monitorHandle = monitor.handle
	}

	var shareHandle uintptr
	if share != nil {
		shareHandle = share.handle
	}

	clone := *windowProtoPtr

	newHandle, _, _ := purego.SyscallN(
		clone.glfwCreateWindow,
		uintptr(width),
		uintptr(height),
		uintptr(unsafe.Pointer(&titlePtr[0])),
		monitorHandle,
		shareHandle,
	)

	if newHandle == 0 {
		fmt.Println("glfwCreateWindow failed (new handle is zero)")
		return nil
	}

	clone.handle = newHandle

	return &clone
}
