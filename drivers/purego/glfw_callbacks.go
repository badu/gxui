package purego

// CursorPosCallback is the function signature for cursor position callbacks
type CursorPosCallback func(window *Window, xpos, ypos float64)

// KeyCallback is the function signature for key callbacks
type KeyCallback func(window *Window, key Key, scancode int32, action Action, mods ModifierKey)

// CharCallback is the function signature for character callbacks
type CharCallback func(window *Window, char rune)

// ScrollCallback is the function signature for scroll callbacks
type ScrollCallback func(window *Window, xoff, yoff float64)

// MouseButtonCallback is the function signature for mouse button callbacks
type MouseButtonCallback func(window *Window, button MouseButton, action Action, mods ModifierKey)

// FramebufferSizeCallback is the function signature for framebuffer size callbacks
type FramebufferSizeCallback func(window *Window, width, height int32)

// CloseCallback is the function signature for window close callbacks
type CloseCallback func(window *Window)

// RefreshCallback is the function signature for window refresh callbacks
type RefreshCallback func(window *Window)

// SizeCallback is the function signature for window size callbacks
type SizeCallback func(window *Window, width, height int32)

// CursorEnterCallback is the function signature for cursor enter/leave callbacks
type CursorEnterCallback func(window *Window, entered bool)

// CharModsCallback is the function signature for character with modifiers callbacks
type CharModsCallback func(window *Window, char rune, mods ModifierKey)

// PosCallback is the function signature for window position callbacks
type PosCallback func(window *Window, xpos, ypos int32)

// FocusCallback is the function signature for window focus callbacks
type FocusCallback func(window *Window, focused bool)

// IconifyCallback is the function signature for window iconify callbacks
type IconifyCallback func(window *Window, iconified bool)

// DropCallback is the function signature for file drop callbacks
type DropCallback func(window *Window, paths []string)
