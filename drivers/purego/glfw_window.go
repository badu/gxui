package purego

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

func (w *Window) GetPos() (xpos, ypos int32) {
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

// ShouldClose returns the value of the close flag for this window
func (w *Window) ShouldClose() bool {
	return WindowShouldClose(w)
}

// SetShouldClose sets the value of the close flag for this window
func (w *Window) SetShouldClose(value bool) {
	SetWindowShouldClose(w, value)
}

// Maximize maximizes this window
func (w *Window) Maximize() {
	MaximizeWindow(w)
}

// Iconify iconifies (minimizes) this window
func (w *Window) Iconify() {
	IconifyWindow(w)
}

// Restore restores this window
func (w *Window) Restore() {
	RestoreWindow(w)
}

// GetWindowSize retrieves the size of the content area for this window
func (w *Window) GetWindowSize() (width, height int) {
	return GetWindowSize(w)
}

// GetWindowPos retrieves the position of the content area for this window
func (w *Window) GetWindowPos() (xpos, ypos int32) {
	return GetWindowPos(w)
}

// SetWindowPos sets the position of the content area for this window
func (w *Window) SetWindowPos(xpos, ypos int) {
	SetWindowPos(w, xpos, ypos)
}

// SetWindowSize sets the size of the content area for this window
func (w *Window) SetWindowSize(width, height int) {
	SetWindowSize(w, width, height)
}

// SetWindowTitle sets the title for this window
func (w *Window) SetWindowTitle(title string) {
	SetWindowTitle(w, title)
}
