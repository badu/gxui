//go:build !js
// +build !js

package cgo

import (
	"io"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

// Init initializes the library.
//
// A valid ContextWatcher must be provided. It gets notified when context becomes current or detached.
// It should be provided by the GL bindings you are using, so you can do glfw.Init(gl.ContextWatcher).
func Init() error {
	return glfw.Init()
}

func Terminate() {
	glfw.Terminate()
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var m *glfw.Monitor
	if monitor != nil {
		m = monitor.Monitor
	}

	var s *glfw.Window
	if share != nil {
		s = share.Window
	}

	w, err := glfw.CreateWindow(width, height, title, m, s)
	if err != nil {
		return nil, err
	}

	window := &Window{Window: w}

	return window, err
}

func SwapInterval(interval int) {
	glfw.SwapInterval(interval)
}

func (w *Window) MakeContextCurrent() {
	w.Window.MakeContextCurrent()
}

func DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

type Window struct {
	*glfw.Window
}

type Monitor struct {
	*glfw.Monitor
}

func GetPrimaryMonitor() *Monitor {
	m := glfw.GetPrimaryMonitor()
	return &Monitor{Monitor: m}
}

func PollEvents() {
	glfw.PollEvents()
}

type CursorPosCallback func(w *Window, xpos float64, ypos float64)

func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos float64, ypos float64) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetCursorPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type MouseMovementCallback func(w *Window, xpos, ypos, xdelta, ydelta float64)

var lastMousePos [2]float64 // HACK.

// TODO: For now, this overrides SetCursorPosCallback; should support both.
func (w *Window) SetMouseMovementCallback(cbfun MouseMovementCallback) (previous MouseMovementCallback) {
	lastMousePos[0], lastMousePos[1] = w.Window.GetCursorPos()
	wrappedCbfun := func(_ *glfw.Window, xpos, ypos float64) {
		xdelta, ydelta := xpos-lastMousePos[0], ypos-lastMousePos[1]
		lastMousePos[0], lastMousePos[1] = xpos, ypos
		cbfun(w, xpos, ypos, xdelta, ydelta)
	}

	p := w.Window.SetCursorPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	wrappedCbfun := func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, key, scancode, action, mods)
	}

	p := w.Window.SetKeyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CharCallback func(w *Window, char rune)

func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	wrappedCbfun := func(_ *glfw.Window, char rune) {
		cbfun(w, char)
	}

	p := w.Window.SetCharCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type ScrollCallback func(w *Window, xoff float64, yoff float64)

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	wrappedCbfun := func(_ *glfw.Window, xoff float64, yoff float64) {
		cbfun(w, xoff, yoff)
	}

	p := w.Window.SetScrollCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	wrappedCbfun := func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, button, action, mods)
	}

	p := w.Window.SetMouseButtonCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FramebufferSizeCallback func(w *Window, width int, height int)

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetFramebufferSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

func (w *Window) GetKey(key Key) Action {
	return w.Window.GetKey(key)
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	return w.Window.GetMouseButton(button)
}

func (w *Window) GetInputMode(mode InputMode) int {
	return w.Window.GetInputMode(mode)
}

func (w *Window) SetInputMode(mode InputMode, value int) {
	w.Window.SetInputMode(mode, value)
}

type Key = glfw.Key

const (
	KeyF13 = glfw.KeyF13
	KeyF14 = glfw.KeyF14
	KeyF15 = glfw.KeyF15
	KeyF16 = glfw.KeyF16
	KeyF17 = glfw.KeyF17
	KeyF18 = glfw.KeyF18
	KeyF19 = glfw.KeyF19
	KeyF20 = glfw.KeyF20
	KeyF21 = glfw.KeyF21
	KeyF22 = glfw.KeyF22
	KeyF23 = glfw.KeyF23
	KeyF24 = glfw.KeyF24
	KeyF25 = glfw.KeyF25
)

type MouseButton = glfw.MouseButton

const (
	MouseButton1 = glfw.MouseButton1
	MouseButton2 = glfw.MouseButton2
	MouseButton3 = glfw.MouseButton3
)

type Action = glfw.Action

type InputMode = glfw.InputMode

const (
	CursorMode             = glfw.CursorMode
	StickyKeysMode         = glfw.StickyKeysMode
	StickyMouseButtonsMode = glfw.StickyMouseButtonsMode
)

const (
	CursorNormal   = glfw.CursorNormal
	CursorHidden   = glfw.CursorHidden
	CursorDisabled = glfw.CursorDisabled
)

type ModifierKey = glfw.ModifierKey

// Open opens a named asset. It's the caller's responsibility to close it when done.
//
// For now, assets are read directly from the current working directory.
func Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// ---

func WaitEvents() {
	glfw.WaitEvents()
}

func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func DefaultWindowHints() {
	glfw.DefaultWindowHints()
}

type CloseCallback func(w *Window)

func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetCloseCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type RefreshCallback func(w *Window)

func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetRefreshCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type SizeCallback func(w *Window, width int, height int)

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CursorEnterCallback func(w *Window, entered bool)

func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {
	wrappedCbfun := func(_ *glfw.Window, entered bool) {
		cbfun(w, entered)
	}

	p := w.Window.SetCursorEnterCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CharModsCallback func(w *Window, char rune, mods ModifierKey)

func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	wrappedCbfun := func(_ *glfw.Window, char rune, mods glfw.ModifierKey) {
		cbfun(w, char, mods)
	}

	p := w.Window.SetCharModsCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type PosCallback func(w *Window, xpos int, ypos int)

func (w *Window) SetPosCallback(cbfun PosCallback) (previous PosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos int, ypos int) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FocusCallback func(w *Window, focused bool)

func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	wrappedCbfun := func(_ *glfw.Window, focused bool) {
		cbfun(w, focused)
	}

	p := w.Window.SetFocusCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type IconifyCallback func(w *Window, iconified bool)

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	wrappedCbfun := func(_ *glfw.Window, iconified bool) {
		cbfun(w, iconified)
	}

	p := w.Window.SetIconifyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	wrappedCbfun := func(_ *glfw.Window, names []string) {
		cbfun(w, names)
	}

	p := w.Window.SetDropCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type Hint int

const (
	ClientAPI = Hint(glfw.ClientAPI)

	AlphaBits   = Hint(glfw.AlphaBits)
	DepthBits   = Hint(glfw.DepthBits)
	StencilBits = Hint(glfw.StencilBits)
	Samples     = Hint(glfw.Samples)
	Resizable   = Hint(glfw.Resizable)

	// These hints used for WebGL contexts, ignored on desktop.
	PremultipliedAlpha = noopHint
	PreserveDrawingBuffer
	PreferLowPowerToHighPerformance
	FailIfMajorPerformanceCaveat
)

const (
	NoAPI = glfw.NoAPI
)

// noopHint is ignored.
const noopHint Hint = -1

func WindowHint(target Hint, hint int) {
	if target == noopHint {
		return
	}

	glfw.WindowHint(glfw.Hint(target), hint)
}
