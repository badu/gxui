//go:build windows
// +build windows

// NOTE : vibe coded as a POC using Chat GPT
package purego

// Minimal GLFW-like "Window" implementation for Windows (user32/win32).
// Implements a small subset needed by gxui's purego driver:
// - CreateWindow
// - Destroy
// - SetTitle
// - GetFramebufferSize
// - ShouldClose
// - PollEvents
// - SetCursorPosCallback
// - SetKeyCallback
// - SetMouseButtonCallback
// - SetCloseCallback
// - SwapBuffers (uses GDI SwapBuffers on HDC)
// - MakeContextCurrent (no-op, placeholder)
//
// This is intentionally small and focuses on correctness and safety for callback
// dispatch. The message loop (PollEvents) must be called from the OS thread that
// created the window (see runtime.LockOSThread()).

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Windows API imports
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	gdi32                = syscall.NewLazyDLL("gdi32.dll")
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procPeekMessageW     = user32.NewProc("PeekMessageW")
	procGetMessageW      = user32.NewProc("GetMessageW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procUpdateWindow     = user32.NewProc("UpdateWindow")
	procSetWindowTextW   = user32.NewProc("SetWindowTextW")
	procGetClientRect    = user32.NewProc("GetClientRect")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procSwapBuffers      = gdi32.NewProc("SwapBuffers")
	procLoadCursorW      = user32.NewProc("LoadCursorW")
	procSetCursor        = user32.NewProc("SetCursor")

	// class/instance counters
	wndClassAtom uintptr
)

const (
	WM_DESTROY     = 0x0002
	WM_CLOSE       = 0x0010
	WM_SIZE        = 0x0005
	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	PM_REMOVE      = 0x0001
	SW_SHOW        = 5
)

// Window callbacks signatures (subset)
type CursorPosCallback func(w *Window, xpos, ypos float64)
type KeyCallback func(w *Window, key int, scancode int, action int, mods int)
type MouseButtonCallback func(w *Window, button int, action int, mods int)
type CloseCallback func(w *Window)

// Window represents a GLFW-like window
type Window struct {
	hwnd        syscall.Handle
	shouldClose bool

	cursorCb CursorPosCallback
	keyCb    KeyCallback
	mouseCb  MouseButtonCallback
	closeCb  CloseCallback

	// internal: event queue locking not included; callbacks are invoked on
	// thread that runs PollEvents. If you want callbacks on other goroutines,
	// make them push events into channels.
}

// helper: UTF-16 conversion
func utf16PtrFromString(s string) *uint16 {
	if s == "" {
		return &[]uint16{0}[0]
	}
	u := utf16.Encode([]rune(s + "\x00"))
	return &u[0]
}

// CreateWindow creates a top-level window. width/height are client area size (in pixels).
func CreateWindow(title string, width, height int) (*Window, error) {
	className, err := registerWindowClass()
	if err != nil {
		return nil, err
	}

	// CW_USEDEFAULT etc omitted; create window with given client size roughly
	tt := utf16PtrFromString(title)
	cn := utf16PtrFromString(className)

	hwndRaw, _, err := procCreateWindowExW.Call(
		0, // dwExStyle
		uintptr(unsafe.Pointer(cn)),
		uintptr(unsafe.Pointer(tt)),
		0xcf0000,   // WS_OVERLAPPEDWINDOW
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		uintptr(width),
		uintptr(height),
		0, // hWndParent
		0, // hMenu
		0, // hInstance
		0) // lpParam
	if hwndRaw == 0 {
		return nil, err
	}

	hwnd := syscall.Handle(hwndRaw)

	// show & update window
	procShowWindow.Call(uintptr(hwnd), SW_SHOW)
	procUpdateWindow.Call(uintptr(hwnd))

	w := &Window{hwnd: hwnd}
	// set the window long pointer to associate Go struct pointer with HWND
	setWindowUserData(hwnd, uintptr(unsafe.Pointer(w)))

	return w, nil
}

// registerWindowClass registers a single WNDCLASS and returns its name
func registerWindowClass() (string, error) {
	// only register once
	if wndClassAtom != 0 {
		return "go_glfw_class", nil
	}

	// minimal WNDCLASSEX struct
	// We'll call RegisterClassExW via syscall; build struct inline as bytes

	// Use our Go WndProc wrapper
	wndProc := syscall.NewCallback(wndProc)

	// prepare WNDCLASSEXW
	// typedef struct tagWNDCLASSEXW {
	//   UINT cbSize; WNDPROC lpfnWndProc; int cbClsExtra; int cbWndExtra;
	//   HINSTANCE hInstance; HICON hIcon; HCURSOR hCursor; HBRUSH hbrBackground;
	//   LPCWSTR lpszMenuName; LPCWSTR lpszClassName; HICON hIconSm;
	// } WNDCLASSEXW;

	// set up parameters by calling RegisterClassExW using global memory layout
	// Simpler: call RegisterClassExW via proc with pointer to stack struct

	// Build WNDCLASSEXW in memory
	type wndclassex struct {
		cbSize        uint32
		style         uint32
		lpfnWndProc   uintptr
		cbClsExtra    int32
		cbWndExtra    int32
		hInstance     uintptr
		hIcon         uintptr
		hCursor       uintptr
		hbrBackground uintptr
		lpszMenuName  uintptr
		lpszClassName uintptr
		hIconSm       uintptr
	}

	hinst := uintptr(0)
	cursor, _, _ := procLoadCursorW.Call(0, uintptr(32512)) // IDC_ARROW

	wc := wndclassex{
		cbSize:        uint32(unsafe.Sizeof(wndclassex{})),
		style:         0,
		lpfnWndProc:   wndProc,
		cbClsExtra:    0,
		cbWndExtra:    0,
		hInstance:     hinst,
		hIcon:         0,
		hCursor:       cursor,
		hbrBackground: 0,
		lpszMenuName:  0,
		lpszClassName: uintptr(unsafe.Pointer(utf16PtrFromString("go_glfw_class"))),
		hIconSm:       0,
	}

	r1, _, e1 := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if r1 == 0 {
		return "", e1
	}
	wndClassAtom = r1
	return "go_glfw_class", nil
}

// setWindowUserData stores a pointer-sized value in GWLP_USERDATA. We keep
// a pointer to the Go Window struct so wndProc can dispatch to Go methods.
func setWindowUserData(hwnd syscall.Handle, val uintptr) {
	procSetWindowLongPtr := user32.NewProc("SetWindowLongPtrW")
	procSetWindowLongPtr.Call(uintptr(hwnd), uintptr(-21), val) // GWLP_USERDATA = -21
}

func getWindowUserData(hwnd syscall.Handle) uintptr {
	procGetWindowLongPtr := user32.NewProc("GetWindowLongPtrW")
	r, _, _ := procGetWindowLongPtr.Call(uintptr(hwnd), uintptr(-21))
	return r
}

// wndProc receives Windows messages and dispatches them into our Go callbacks
func wndProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	// retrieve Go Window pointer
	u := getWindowUserData(hwnd)
	var w *Window
	if u != 0 {
		w = (*Window)(unsafe.Pointer(u))
	}

	switch msg {
	case WM_MOUSEMOVE:
		if w != nil && w.cursorCb != nil {
			x := float64(int32(lParam & 0xFFFF))
			y := float64(int32((lParam >> 16) & 0xFFFF))
			// invoke directly (on PollEvents thread)
			w.cursorCb(w, x, y)
		}
	case WM_LBUTTONDOWN:
		if w != nil && w.mouseCb != nil {
			w.mouseCb(w, 0, 1, 0) // button 0, action 1=press
		}
	case WM_LBUTTONUP:
		if w != nil && w.mouseCb != nil {
			w.mouseCb(w, 0, 0, 0) // button 0, action 0=release
		}
	case WM_RBUTTONDOWN:
		if w != nil && w.mouseCb != nil {
			w.mouseCb(w, 1, 1, 0)
		}
	case WM_RBUTTONUP:
		if w != nil && w.mouseCb != nil {
			w.mouseCb(w, 1, 0, 0)
		}
	case WM_KEYDOWN:
		if w != nil && w.keyCb != nil {
			w.keyCb(w, int(wParam), 0, 1, 0)
		}
	case WM_KEYUP:
		if w != nil && w.keyCb != nil {
			w.keyCb(w, int(wParam), 0, 0, 0)
		}
	case WM_CLOSE:
		if w != nil {
			w.shouldClose = true
			if w.closeCb != nil {
				w.closeCb(w)
			}
		}
		procPostQuitMessage.Call(0)
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
	default:
		// fallthrough to default
		ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
	return 0
}

// PollEvents pumps the Windows message queue and dispatches callbacks.
// It MUST be called on the OS thread that created the window (LockOSThread).
func PollEvents() {
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		ptx     int32
		pty     int32
	}

	for {
		r, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, PM_REMOVE)
		if r == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// Setters for callbacks
func (w *Window) SetCursorPosCallback(cb CursorPosCallback) {
	w.cursorCb = cb
}
func (w *Window) SetKeyCallback(cb KeyCallback) {
	w.keyCb = cb
}
func (w *Window) SetMouseButtonCallback(cb MouseButtonCallback) {
	w.mouseCb = cb
}
func (w *Window) SetCloseCallback(cb CloseCallback) {
	w.closeCb = cb
}

// ShouldClose reports whether a close was requested
func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

// SetTitle sets window title
func (w *Window) SetTitle(title string) {
	procSetWindowTextW.Call(uintptr(w.hwnd), uintptr(unsafe.Pointer(utf16PtrFromString(title))))
}

// GetFramebufferSize returns the client area size in pixels
func (w *Window) GetFramebufferSize() (int, int) {
	var rect struct{ left, top, right, bottom int32 }
	procGetClientRect.Call(uintptr(w.hwnd), uintptr(unsafe.Pointer(&rect)))
	widh := int(rect.right - rect.left)
	hei := int(rect.bottom - rect.top)
	return widh, hei
}

// SwapBuffers does a GDI SwapBuffers on the window's DC. This is a thin wrapper
// and assumes an OpenGL context is current on the calling thread for this window.
func (w *Window) SwapBuffers() {
	hdc, _, _ := procGetDC.Call(uintptr(w.hwnd))
	if hdc == 0 {
		return
	}
	procSwapBuffers.Call(hdc)
	procReleaseDC.Call(uintptr(w.hwnd), hdc)
}

// MakeContextCurrent is a placeholder — the user said OpenGL is handled elsewhere.
func (w *Window) MakeContextCurrent() {
	// no-op here; user manages GL context
}

// Destroy destroys the native window
func (w *Window) Destroy() {
	if w == nil || w.hwnd == 0 {
		return
	}
	user32.NewProc("DestroyWindow").Call(uintptr(w.hwnd))
	w.hwnd = 0
}

// --- End of file
