package purego

import (
	"runtime"
)

const (
	Press   = 1
	Release = 2
	Repeat  = 3

	KeySpace        Key = 32
	KeyApostrophe   Key = 39 /* ' */
	KeyComma        Key = 44 /* , */
	KeyMinus        Key = 45 /* - */
	KeyPeriod       Key = 46 /* . */
	KeySlash        Key = 47 /* / */
	Key0            Key = 48
	Key1            Key = 49
	Key2            Key = 50
	Key3            Key = 51
	Key4            Key = 52
	Key5            Key = 53
	Key6            Key = 54
	Key7            Key = 55
	Key8            Key = 56
	Key9            Key = 57
	KeySemicolon    Key = 59 /* ; */
	KeyEqual        Key = 61 /* = */
	KeyA            Key = 65
	KeyB            Key = 66
	KeyC            Key = 67
	KeyD            Key = 68
	KeyE            Key = 69
	KeyF            Key = 70
	KeyG            Key = 71
	KeyH            Key = 72
	KeyI            Key = 73
	KeyJ            Key = 74
	KeyK            Key = 75
	KeyL            Key = 76
	KeyM            Key = 77
	KeyN            Key = 78
	KeyO            Key = 79
	KeyP            Key = 80
	KeyQ            Key = 81
	KeyR            Key = 82
	KeyS            Key = 83
	KeyT            Key = 84
	KeyU            Key = 85
	KeyV            Key = 86
	KeyW            Key = 87
	KeyX            Key = 88
	KeyY            Key = 89
	KeyZ            Key = 90
	KeyLeftBracket  Key = 91  /* [ */
	KeyBackslash    Key = 92  /* \ */
	KeyRightBracket Key = 93  /* ] */
	KeyGraveAccent  Key = 96  /* ` */
	KeyWorld1       Key = 161 /* non-US #1 */
	KeyWorld2       Key = 162 /* non-US #2 */
	KeyEscape       Key = 256
	KeyEnter        Key = 257
	KeyTab          Key = 258
	KeyBackspace    Key = 259
	KeyInsert       Key = 260
	KeyDelete       Key = 261
	KeyRight        Key = 262
	KeyLeft         Key = 263
	KeyDown         Key = 264
	KeyUp           Key = 265
	KeyPageUp       Key = 266
	KeyPageDown     Key = 267
	KeyHome         Key = 268
	KeyEnd          Key = 269
	KeyCapsLock     Key = 280
	KeyScrollLock   Key = 281
	KeyNumLock      Key = 282
	KeyPrintScreen  Key = 283
	KeyPause        Key = 284
	KeyF1           Key = 290
	KeyF2           Key = 291
	KeyF3           Key = 292
	KeyF4           Key = 293
	KeyF5           Key = 294
	KeyF6           Key = 295
	KeyF7           Key = 296
	KeyF8           Key = 297
	KeyF9           Key = 298
	KeyF10          Key = 299
	KeyF11          Key = 300
	KeyF12          Key = 301
	KeyKP0          Key = 320
	KeyKP1          Key = 321
	KeyKP2          Key = 322
	KeyKP3          Key = 323
	KeyKP4          Key = 324
	KeyKP5          Key = 325
	KeyKP6          Key = 326
	KeyKP7          Key = 327
	KeyKP8          Key = 328
	KeyKP9          Key = 329
	KeyKPDecimal    Key = 330
	KeyKPDivide     Key = 331
	KeyKPMultiply   Key = 332
	KeyKPSubtract   Key = 333
	KeyKPAdd        Key = 334
	KeyKPEnter      Key = 335
	KeyKPEqual      Key = 336
	KeyLeftShift    Key = 340
	KeyLeftControl  Key = 341
	KeyLeftAlt      Key = 342
	KeyLeftSuper    Key = 343
	KeyRightShift   Key = 344
	KeyRightControl Key = 345
	KeyRightAlt     Key = 346
	KeyRightSuper   Key = 347
	KeyMenu         Key = 348

	ModShift   ModifierKey = 1
	ModControl ModifierKey = 1
	ModAlt     ModifierKey = 1
	ModSuper   ModifierKey = 1
)

type MouseButton = int32
type Action = int32
type ModifierKey = int32
type Key = int32
type InputMode = int32

func init() {
	runtime.LockOSThread()
}

type Window struct {
	handle uintptr
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
	MakeContextCurrent(w)
}

// GetCursorPos retrieves the cursor position for this window
func (w *Window) GetCursorPos() (xpos, ypos float64) {
	return GetCursorPos(w)
}

// SetCursorPosCallback sets the cursor position callback for this window
func (w *Window) SetCursorPosCallback(callback CursorPosCallback) {
	SetCursorPosCallback(w, callback)
}

// SetKeyCallback sets the key callback for this window
func (w *Window) SetKeyCallback(callback KeyCallback) {
	SetKeyCallback(w, callback)
}

// SetCharCallback sets the character callback for this window
func (w *Window) SetCharCallback(callback CharCallback) {
	SetCharCallback(w, callback)
}

// SetScrollCallback sets the scroll callback for this window
func (w *Window) SetScrollCallback(callback ScrollCallback) {
	SetScrollCallback(w, callback)
}

// SetMouseButtonCallback sets the mouse button callback for this window
func (w *Window) SetMouseButtonCallback(callback MouseButtonCallback) {
	SetMouseButtonCallback(w, callback)
}

// SetFramebufferSizeCallback sets the framebuffer size callback for this window
func (w *Window) SetFramebufferSizeCallback(callback FramebufferSizeCallback) {
	SetFramebufferSizeCallback(w, callback)
}

// SetCloseCallback sets the window close callback for this window
func (w *Window) SetCloseCallback(callback CloseCallback) {
	SetCloseCallback(w, callback)
}

// SetRefreshCallback sets the window refresh callback for this window
func (w *Window) SetRefreshCallback(callback RefreshCallback) {
	SetRefreshCallback(w, callback)
}

// SetSizeCallback sets the window size callback for this window
func (w *Window) SetSizeCallback(callback SizeCallback) {
	SetSizeCallback(w, callback)
}

// SetCursorEnterCallback sets the cursor enter/leave callback for this window
func (w *Window) SetCursorEnterCallback(callback CursorEnterCallback) {
	SetCursorEnterCallback(w, callback)
}

// SetCharModsCallback sets the character with modifiers callback for this window
func (w *Window) SetCharModsCallback(callback CharModsCallback) {
	SetCharModsCallback(w, callback)
}

// SetPosCallback sets the window position callback for this window
func (w *Window) SetPosCallback(callback PosCallback) {
	SetPosCallback(w, callback)
}

// SetFocusCallback sets the window focus callback for this window
func (w *Window) SetFocusCallback(callback FocusCallback) {
	SetFocusCallback(w, callback)
}

// SetIconifyCallback sets the window iconify callback for this window
func (w *Window) SetIconifyCallback(callback IconifyCallback) {
	SetIconifyCallback(w, callback)
}

// SetDropCallback sets the file drop callback for this window
func (w *Window) SetDropCallback(callback DropCallback) {
	SetDropCallback(w, callback)
}

// SetClipboardString sets the clipboard to the specified string for this window
func (w *Window) SetClipboardString(str string) {
	SetClipboardString(w, str)
}

func (w *Window) GetClipboardString() string {
	return GetClipboardString(w)
}

func (w *Window) GetKey(key Key) Action {
	return GetWindowKey(w, key)
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	return GetWindowMouseButton(w, button)
}

func (w *Window) GetInputMode(mode InputMode) int32 {
	return GetWindowInputMode(w, mode)
}

func (w *Window) SetInputMode(mode InputMode, value int32) {
	SetWindowInputMode(w, mode, value)
}

func (w *Window) GetFramebufferSize() (width, height int32) {
	return GetWindowFramebufferSize(w)
}

func (w *Window) GetPos() (xpos, ypos int) {
	return GetWindowPos(w)
}

func (w *Window) SetPos(xpos, ypos int) {
	SetWindowPos(w, xpos, ypos)
}

func (w *Window) GetSize() (width, height int) {
	return GetWindowSize(w)
}

func (w *Window) SetSize(width, height int) {
	SetWindowSize(w, width, height)
}

func (w *Window) SetTitle(title string) {
	SetWindowTitle(w, title)
}

func (w *Window) SwapBuffers() {
	SwapWindowBuffers(w)
}

func (w *Window) Show() {
	ShowWindow(w)
}

func (w *Window) Hide() {
	HideWindow(w)
}

func (w *Window) Destroy() {
	DestroyWindow(w)
}

// Monitor represents a GLFW monitor handle
type Monitor struct {
	handle uintptr
}

// Handle returns the raw uintptr handle
func (m *Monitor) Handle() uintptr {
	if m == nil {
		return 0
	}
	return m.handle
}

type VidMode struct {
	Width       int // The width, in pixels, of the video mode.
	Height      int // The height, in pixels, of the video mode.
	RedBits     int // The bit depth of the red channel of the video mode.
	GreenBits   int // The bit depth of the green channel of the video mode.
	BlueBits    int // The bit depth of the blue channel of the video mode.
	RefreshRate int // The refresh rate, in Hz, of the video mode.
}

func (m *Monitor) GetVideoMode() *VidMode {
	return GetVideoMode(m)
}

type Hint int

const (
	Samples = Hint(SAMPLES)
)
