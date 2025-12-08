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
	glfwInit                       uintptr
	glfwTerminate                  uintptr
	glfwWindowHint                 uintptr
	glfwDefaultWindowHints         uintptr
	glfwCreateWindow               uintptr
	glfwGetPrimaryMonitor          uintptr
	glfwPollEvents                 uintptr
	glfwWaitEvents                 uintptr
	glfwPostEmptyEvent             uintptr
	glfwSwapInterval               uintptr
	glfwMakeContextCurrent         uintptr
	glfwGetCursorPos               uintptr
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
	glfwSetClipboardString         uintptr
	glfwGetClipboardString         uintptr
	glfwGetVideoMode               uintptr
	glfwGetWindowSize              uintptr
	glfwGetFramebufferSize         uintptr
	glfwGetKey                     uintptr
	glfwGetMouseButton             uintptr
	glfwGetInputMode               uintptr
	glfwSetInputMode               uintptr
	glfwGetWindowPos               uintptr
	glfwSetWindowPos               uintptr
	glfwSetWindowSize              uintptr
	glfwSetWindowTitle             uintptr
	glfwSwapBuffers                uintptr
	glfwShowWindow                 uintptr
	glfwHideWindow                 uintptr
	glfwDestroyWindow              uintptr
	glfwGetMonitorWorkarea         uintptr
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

// CreateWindow creates a window and its associated context
func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) *Window {
	titlePtr := append([]byte(title), 0)

	var monitorHandle uintptr
	if monitor != nil {
		monitorHandle = monitor.handle
	}

	var shareHandle uintptr
	if share != nil {
		shareHandle = share.handle
	}

	ret, _, _ := purego.SyscallN(
		glfwCreateWindow,
		uintptr(width),
		uintptr(height),
		uintptr(unsafe.Pointer(&titlePtr[0])),
		monitorHandle,
		shareHandle,
	)

	if ret == 0 {
		return nil
	}

	return &Window{handle: ret}
}

// GetPrimaryMonitor returns the primary monitor
func GetPrimaryMonitor() *Monitor {
	ret, _, _ := purego.SyscallN(glfwGetPrimaryMonitor)
	if ret == 0 {
		return nil
	}
	return &Monitor{handle: ret}
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

// MakeContextCurrent makes the OpenGL context of the specified window current
func MakeContextCurrent(window *Window) {
	var handle uintptr
	if window != nil {
		handle = window.handle
	}
	purego.SyscallN(glfwMakeContextCurrent, handle)
}

// DetachCurrentContext detaches the current context from the current thread
// This is done by calling MakeContextCurrent with NULL (0)
func DetachCurrentContext() {
	MakeContextCurrent(nil)
}

// GetCursorPos retrieves the position of the cursor relative to the content area of the window
func GetCursorPos(window *Window) (xpos, ypos float64) {
	if window == nil {
		return 0, 0
	}

	purego.SyscallN(glfwGetCursorPos, window.handle, uintptr(unsafe.Pointer(&xpos)), uintptr(unsafe.Pointer(&ypos)))
	return
}

// SetClipboardString sets the clipboard to the specified string
func SetClipboardString(window *Window, str string) {
	if window == nil {
		return
	}
	strBytes := append([]byte(str), 0)
	purego.SyscallN(glfwSetClipboardString, window.handle, uintptr(unsafe.Pointer(&strBytes[0])))
}

// GetClipboardString sets the clipboard to the specified string
func GetClipboardString(window *Window) string {
	if window == nil {
		return ""
	}
	ret, _, _ := purego.SyscallN(glfwGetClipboardString, window.handle)
	return cStringToGoString(ret)
}

func GetVideoMode(monitor *Monitor) *VidMode {
	if monitor == nil {
		return nil
	}

	vmPtr, _, _ := purego.SyscallN(glfwGetVideoMode, monitor.handle)
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

func GetWindowSize(window *Window) (width, height int) {
	if window == nil {
		return 0, 0
	}

	purego.SyscallN(glfwGetWindowSize, window.handle, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return width, height
}

func GetWindowKey(window *Window, key Key) Action {
	if window == nil {
		return 0
	}

	ret, _, _ := purego.SyscallN(glfwGetKey, window.handle, uintptr(key))
	return Action(ret)
}

func GetWindowMouseButton(window *Window, button MouseButton) Action {
	if window == nil {
		return 0
	}
	ret, _, _ := purego.SyscallN(glfwGetMouseButton, window.handle, uintptr(button))
	return Action(ret)
}

func GetWindowInputMode(window *Window, mode InputMode) int32 {
	if window == nil {
		return 0
	}

	ret, _, _ := purego.SyscallN(glfwGetInputMode, window.handle, uintptr(mode))
	return int32(ret)
}

func SetWindowInputMode(window *Window, mode InputMode, value int32) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwSetInputMode, window.handle, uintptr(mode), uintptr(value))
}

func GetWindowFramebufferSize(window *Window) (int32, int32) {
	if window == nil {
		return 0, 0
	}

	var width, height int32
	purego.SyscallN(glfwGetFramebufferSize, window.handle, uintptr(unsafe.Pointer(&width)), uintptr(unsafe.Pointer(&height)))
	return width, height
}

// GetWindowPos retrieves the position of the content area of the specified window
func GetWindowPos(window *Window) (xpos, ypos int) {
	if window == nil {
		return 0, 0
	}
	purego.SyscallN(glfwGetWindowPos, window.handle, uintptr(unsafe.Pointer(&xpos)), uintptr(unsafe.Pointer(&ypos)))
	return
}

// SetWindowPos sets the position of the content area of the specified window
func SetWindowPos(window *Window, xpos, ypos int) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwSetWindowPos, window.handle, uintptr(xpos), uintptr(ypos))
}

// SetWindowSize sets the size of the content area of the specified window
func SetWindowSize(window *Window, width, height int) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwSetWindowSize, window.handle, uintptr(width), uintptr(height))
}

// SetWindowTitle sets the title of the specified window
func SetWindowTitle(window *Window, title string) {
	if window == nil {
		return
	}
	titleBytes := append([]byte(title), 0)
	purego.SyscallN(glfwSetWindowTitle, window.handle, uintptr(unsafe.Pointer(&titleBytes[0])))
}

// SwapBuffers swaps the front and back buffers of the specified window
func SwapWindowBuffers(window *Window) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwSwapBuffers, window.handle)
}

// ShowWindow makes the specified window visible
func ShowWindow(window *Window) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwShowWindow, window.handle)
}

// HideWindow hides the specified window
func HideWindow(window *Window) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwHideWindow, window.handle)
}

// DestroyWindow destroys the specified window and its context
func DestroyWindow(window *Window) {
	if window == nil {
		return
	}
	purego.SyscallN(glfwDestroyWindow, window.handle)
}

func GetMonitorWorkarea(monitor *Monitor) (xpos, ypos, width, height int) {
	if monitor == nil {
		return 0, 0, 0, 0
	}
	purego.SyscallN(glfwGetMonitorWorkarea, monitor.handle,
		uintptr(unsafe.Pointer(&xpos)),
		uintptr(unsafe.Pointer(&ypos)),
		uintptr(unsafe.Pointer(&width)),
		uintptr(unsafe.Pointer(&height)))
	return
}

// SetCursorPosCallback sets the cursor position callback for the specified window
func SetCursorPosCallback(window *Window, callback CursorPosCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		cursorPosCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(cursorPosCallbackWrapper)
		purego.SyscallN(glfwSetCursorPosCallback, window.handle, callbackPtr)
	} else {
		delete(cursorPosCallbacks, window.handle)
		purego.SyscallN(glfwSetCursorPosCallback, window.handle, 0)
	}
}

// SetKeyCallback sets the key callback for the specified window
func SetKeyCallback(window *Window, callback KeyCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		keyCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(keyCallbackWrapper)
		purego.SyscallN(glfwSetKeyCallback, window.handle, callbackPtr)
	} else {
		delete(keyCallbacks, window.handle)
		purego.SyscallN(glfwSetKeyCallback, window.handle, 0)
	}
}

// SetCharCallback sets the character callback for the specified window
func SetCharCallback(window *Window, callback CharCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		charCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(charCallbackWrapper)
		purego.SyscallN(glfwSetCharCallback, window.handle, callbackPtr)
	} else {
		delete(charCallbacks, window.handle)
		purego.SyscallN(glfwSetCharCallback, window.handle, 0)
	}
}

// SetScrollCallback sets the scroll callback for the specified window
func SetScrollCallback(window *Window, callback ScrollCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		scrollCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(scrollCallbackWrapper)
		purego.SyscallN(glfwSetScrollCallback, window.handle, callbackPtr)
	} else {
		delete(scrollCallbacks, window.handle)
		purego.SyscallN(glfwSetScrollCallback, window.handle, 0)
	}
}

// SetMouseButtonCallback sets the mouse button callback for the specified window
func SetMouseButtonCallback(window *Window, callback MouseButtonCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		mouseButtonCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(mouseButtonCallbackWrapper)
		purego.SyscallN(glfwSetMouseButtonCallback, window.handle, callbackPtr)
	} else {
		delete(mouseButtonCallbacks, window.handle)
		purego.SyscallN(glfwSetMouseButtonCallback, window.handle, 0)
	}
}

// SetFramebufferSizeCallback sets the framebuffer size callback for the specified window
func SetFramebufferSizeCallback(window *Window, callback FramebufferSizeCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		framebufferSizeCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(framebufferSizeCallbackWrapper)
		purego.SyscallN(glfwSetFramebufferSizeCallback, window.handle, callbackPtr)
	} else {
		delete(framebufferSizeCallbacks, window.handle)
		purego.SyscallN(glfwSetFramebufferSizeCallback, window.handle, 0)
	}
}

// SetCloseCallback sets the window close callback for the specified window
func SetCloseCallback(window *Window, callback CloseCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		closeCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(closeCallbackWrapper)
		purego.SyscallN(glfwSetWindowCloseCallback, window.handle, callbackPtr)
	} else {
		delete(closeCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowCloseCallback, window.handle, 0)
	}
}

// SetRefreshCallback sets the window refresh callback for the specified window
func SetRefreshCallback(window *Window, callback RefreshCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		refreshCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(refreshCallbackWrapper)
		purego.SyscallN(glfwSetWindowRefreshCallback, window.handle, callbackPtr)
	} else {
		delete(refreshCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowRefreshCallback, window.handle, 0)
	}
}

// SetSizeCallback sets the window size callback for the specified window
func SetSizeCallback(window *Window, callback SizeCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		sizeCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(sizeCallbackWrapper)
		purego.SyscallN(glfwSetWindowSizeCallback, window.handle, callbackPtr)
	} else {
		delete(sizeCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowSizeCallback, window.handle, 0)
	}
}

// SetCursorEnterCallback sets the cursor enter/leave callback for the specified window
func SetCursorEnterCallback(window *Window, callback CursorEnterCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		cursorEnterCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(cursorEnterCallbackWrapper)
		purego.SyscallN(glfwSetCursorEnterCallback, window.handle, callbackPtr)
	} else {
		delete(cursorEnterCallbacks, window.handle)
		purego.SyscallN(glfwSetCursorEnterCallback, window.handle, 0)
	}
}

// SetCharModsCallback sets the character with modifiers callback for the specified window
func SetCharModsCallback(window *Window, callback CharModsCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		charModsCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(charModsCallbackWrapper)
		purego.SyscallN(glfwSetCharModsCallback, window.handle, callbackPtr)
	} else {
		delete(charModsCallbacks, window.handle)
		purego.SyscallN(glfwSetCharModsCallback, window.handle, 0)
	}
}

// SetPosCallback sets the window position callback for the specified window
func SetPosCallback(window *Window, callback PosCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		posCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(posCallbackWrapper)
		purego.SyscallN(glfwSetWindowPosCallback, window.handle, callbackPtr)
	} else {
		delete(posCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowPosCallback, window.handle, 0)
	}
}

// SetFocusCallback sets the window focus callback for the specified window
func SetFocusCallback(window *Window, callback FocusCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		focusCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(focusCallbackWrapper)
		purego.SyscallN(glfwSetWindowFocusCallback, window.handle, callbackPtr)
	} else {
		delete(focusCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowFocusCallback, window.handle, 0)
	}
}

// SetIconifyCallback sets the window iconify callback for the specified window
func SetIconifyCallback(window *Window, callback IconifyCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		iconifyCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(iconifyCallbackWrapper)
		purego.SyscallN(glfwSetWindowIconifyCallback, window.handle, callbackPtr)
	} else {
		delete(iconifyCallbacks, window.handle)
		purego.SyscallN(glfwSetWindowIconifyCallback, window.handle, 0)
	}
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

// SetDropCallback sets the file drop callback for the specified window
func SetDropCallback(window *Window, callback DropCallback) {
	if window == nil {
		return
	}

	if callback != nil {
		dropCallbacks[window.handle] = callback
		// Register the wrapper function as a callback
		callbackPtr := purego.NewCallback(dropCallbackWrapper)
		purego.SyscallN(glfwSetDropCallback, window.handle, callbackPtr)
	} else {
		delete(dropCallbacks, window.handle)
		purego.SyscallN(glfwSetDropCallback, window.handle, 0)
	}
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

	// Register function symbols
	funcs := map[string]*uintptr{
		"glfwInit":                       &glfwInit,
		"glfwTerminate":                  &glfwTerminate,
		"glfwWindowHint":                 &glfwWindowHint,
		"glfwDefaultWindowHints":         &glfwDefaultWindowHints,
		"glfwCreateWindow":               &glfwCreateWindow,
		"glfwGetPrimaryMonitor":          &glfwGetPrimaryMonitor,
		"glfwPollEvents":                 &glfwPollEvents,
		"glfwWaitEvents":                 &glfwWaitEvents,
		"glfwPostEmptyEvent":             &glfwPostEmptyEvent,
		"glfwSwapInterval":               &glfwSwapInterval,
		"glfwMakeContextCurrent":         &glfwMakeContextCurrent,
		"glfwGetCursorPos":               &glfwGetCursorPos,
		"glfwSetCursorPosCallback":       &glfwSetCursorPosCallback,
		"glfwSetKeyCallback":             &glfwSetKeyCallback,
		"glfwSetCharCallback":            &glfwSetCharCallback,
		"glfwSetScrollCallback":          &glfwSetScrollCallback,
		"glfwSetMouseButtonCallback":     &glfwSetMouseButtonCallback,
		"glfwSetFramebufferSizeCallback": &glfwSetFramebufferSizeCallback,
		"glfwSetWindowCloseCallback":     &glfwSetWindowCloseCallback,
		"glfwSetWindowRefreshCallback":   &glfwSetWindowRefreshCallback,
		"glfwSetWindowSizeCallback":      &glfwSetWindowSizeCallback,
		"glfwSetCursorEnterCallback":     &glfwSetCursorEnterCallback,
		"glfwSetCharModsCallback":        &glfwSetCharModsCallback,
		"glfwSetWindowPosCallback":       &glfwSetWindowPosCallback,
		"glfwSetWindowFocusCallback":     &glfwSetWindowFocusCallback,
		"glfwSetWindowIconifyCallback":   &glfwSetWindowIconifyCallback,
		"glfwSetDropCallback":            &glfwSetDropCallback,
		"glfwSetClipboardString":         &glfwSetClipboardString,
		"glfwGetClipboardString":         &glfwGetClipboardString,
		"glfwGetVideoMode":               &glfwGetVideoMode,
		"glfwGetWindowSize":              &glfwGetWindowSize,
		"glfwGetFramebufferSize":         &glfwGetFramebufferSize,
		"glfwGetKey":                     &glfwGetKey,
		"glfwGetMouseButton":             &glfwGetMouseButton,
		"glfwGetInputMode":               &glfwGetInputMode,
		"glfwSetInputMode":               &glfwSetInputMode,
		"glfwGetWindowPos":               &glfwGetWindowPos,
		"glfwSetWindowPos":               &glfwSetWindowPos,
		"glfwSetWindowSize":              &glfwSetWindowSize,
		"glfwSetWindowTitle":             &glfwSetWindowTitle,
		"glfwSwapBuffers":                &glfwSwapBuffers,
		"glfwShowWindow":                 &glfwShowWindow,
		"glfwHideWindow":                 &glfwHideWindow,
		"glfwDestroyWindow":              &glfwDestroyWindow,
		"glfwGetMonitorWorkarea":         &glfwGetMonitorWorkarea,
	}

	for name, ptr := range funcs {
		sym, err := purego.Dlsym(lib, name)
		if err != nil {
			fmt.Printf("error loading symbol %q : %#v\n", name, err)
			return fmt.Errorf("failed to load symbol %s: %w", name, err)
		}
		*ptr = sym
	}

	return nil
}
