package purego

import "unsafe"

// CursorPosCallback is the function signature for cursor position callbacks
type CursorPosCallback func(window *Window, xpos, ypos float64)

// cursorPosCallbacks stores the Go callbacks for each window
var cursorPosCallbacks = make(map[uintptr]CursorPosCallback)

// cursorPosCallbackWrapper is the C callback that gets called by GLFW
func cursorPosCallbackWrapper(windowHandle uintptr, xpos, ypos float64) {
	if callback, ok := cursorPosCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, xpos, ypos)
	}
}

// KeyCallback is the function signature for key callbacks
type KeyCallback func(window *Window, key Key, scancode int32, action Action, mods ModifierKey)

// keyCallbacks stores the Go callbacks for each window
var keyCallbacks = make(map[uintptr]KeyCallback)

// keyCallbackWrapper is the C callback that gets called by GLFW
func keyCallbackWrapper(windowHandle uintptr, key Key, scancode int32, action Action, mods ModifierKey) {
	if callback, ok := keyCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, key, scancode, action, mods)
	}
}

// CharCallback is the function signature for character callbacks
type CharCallback func(window *Window, char rune)

// charCallbacks stores the Go callbacks for each window
var charCallbacks = make(map[uintptr]CharCallback)

// charCallbackWrapper is the C callback that gets called by GLFW
func charCallbackWrapper(windowHandle uintptr, char rune) {
	if callback, ok := charCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, char)
	}
}

// ScrollCallback is the function signature for scroll callbacks
type ScrollCallback func(window *Window, xoff, yoff float64)

// scrollCallbacks stores the Go callbacks for each window
var scrollCallbacks = make(map[uintptr]ScrollCallback)

// scrollCallbackWrapper is the C callback that gets called by GLFW
func scrollCallbackWrapper(windowHandle uintptr, xoff, yoff float64) {
	if callback, ok := scrollCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, xoff, yoff)
	}
}

// MouseButtonCallback is the function signature for mouse button callbacks
type MouseButtonCallback func(window *Window, button MouseButton, action Action, mods ModifierKey)

// mouseButtonCallbacks stores the Go callbacks for each window
var mouseButtonCallbacks = make(map[uintptr]MouseButtonCallback)

// mouseButtonCallbackWrapper is the C callback that gets called by GLFW
func mouseButtonCallbackWrapper(windowHandle uintptr, button MouseButton, action Action, mods ModifierKey) {
	if callback, ok := mouseButtonCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, button, action, mods)
	}
}

// FramebufferSizeCallback is the function signature for framebuffer size callbacks
type FramebufferSizeCallback func(window *Window, width, height int32)

// framebufferSizeCallbacks stores the Go callbacks for each window
var framebufferSizeCallbacks = make(map[uintptr]FramebufferSizeCallback)

// framebufferSizeCallbackWrapper is the C callback that gets called by GLFW
func framebufferSizeCallbackWrapper(windowHandle uintptr, width, height int32) {
	if callback, ok := framebufferSizeCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, width, height)
	}
}

// CloseCallback is the function signature for window close callbacks
type CloseCallback func(window *Window)

// closeCallbacks stores the Go callbacks for each window
var closeCallbacks = make(map[uintptr]CloseCallback)

// closeCallbackWrapper is the C callback that gets called by GLFW
func closeCallbackWrapper(windowHandle uintptr) {
	if callback, ok := closeCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle})
	}
}

// RefreshCallback is the function signature for window refresh callbacks
type RefreshCallback func(window *Window)

// refreshCallbacks stores the Go callbacks for each window
var refreshCallbacks = make(map[uintptr]RefreshCallback)

// refreshCallbackWrapper is the C callback that gets called by GLFW
func refreshCallbackWrapper(windowHandle uintptr) {
	if callback, ok := refreshCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle})
	}
}

// SizeCallback is the function signature for window size callbacks
type SizeCallback func(window *Window, width, height int32)

// sizeCallbacks stores the Go callbacks for each window
var sizeCallbacks = make(map[uintptr]SizeCallback)

// sizeCallbackWrapper is the C callback that gets called by GLFW
func sizeCallbackWrapper(windowHandle uintptr, width, height int32) {
	if callback, ok := sizeCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, width, height)
	}
}

// CursorEnterCallback is the function signature for cursor enter/leave callbacks
type CursorEnterCallback func(window *Window, entered bool)

// cursorEnterCallbacks stores the Go callbacks for each window
var cursorEnterCallbacks = make(map[uintptr]CursorEnterCallback)

// cursorEnterCallbackWrapper is the C callback that gets called by GLFW
func cursorEnterCallbackWrapper(windowHandle uintptr, entered int) {
	if callback, ok := cursorEnterCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, entered == GLFW_TRUE)
	}
}

// CharModsCallback is the function signature for character with modifiers callbacks
type CharModsCallback func(window *Window, char rune, mods ModifierKey)

// charModsCallbacks stores the Go callbacks for each window
var charModsCallbacks = make(map[uintptr]CharModsCallback)

// charModsCallbackWrapper is the C callback that gets called by GLFW
func charModsCallbackWrapper(windowHandle uintptr, char rune, mods ModifierKey) {
	if callback, ok := charModsCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, char, mods)
	}
}

// PosCallback is the function signature for window position callbacks
type PosCallback func(window *Window, xpos, ypos int32)

// posCallbacks stores the Go callbacks for each window
var posCallbacks = make(map[uintptr]PosCallback)

// posCallbackWrapper is the C callback that gets called by GLFW
func posCallbackWrapper(windowHandle uintptr, xpos, ypos int32) {
	if callback, ok := posCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, xpos, ypos)
	}
}

// FocusCallback is the function signature for window focus callbacks
type FocusCallback func(window *Window, focused bool)

// focusCallbacks stores the Go callbacks for each window
var focusCallbacks = make(map[uintptr]FocusCallback)

// focusCallbackWrapper is the C callback that gets called by GLFW
func focusCallbackWrapper(windowHandle uintptr, focused int) {
	if callback, ok := focusCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, focused == GLFW_TRUE)
	}
}

// IconifyCallback is the function signature for window iconify callbacks
type IconifyCallback func(window *Window, iconified bool)

// iconifyCallbacks stores the Go callbacks for each window
var iconifyCallbacks = make(map[uintptr]IconifyCallback)

// iconifyCallbackWrapper is the C callback that gets called by GLFW
func iconifyCallbackWrapper(windowHandle uintptr, iconified int) {
	if callback, ok := iconifyCallbacks[windowHandle]; ok {
		callback(&Window{handle: windowHandle}, iconified == GLFW_TRUE)
	}
}

// DropCallback is the function signature for file drop callbacks
type DropCallback func(window *Window, paths []string)

// dropCallbacks stores the Go callbacks for each window
var dropCallbacks = make(map[uintptr]DropCallback)

// dropCallbackWrapper is the C callback that gets called by GLFW
func dropCallbackWrapper(windowHandle uintptr, count int, pathsPtr uintptr) {
	if callback, ok := dropCallbacks[windowHandle]; ok {
		// Convert C string array to Go string slice
		paths := make([]string, count)
		for i := 0; i < count; i++ {
			// Get pointer to the i-th string pointer
			strPtr := *(*uintptr)(unsafe.Pointer(pathsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
			// Convert C string to Go string
			paths[i] = cStringToGoString(strPtr)
		}
		callback(&Window{handle: windowHandle}, paths)
	}
}
