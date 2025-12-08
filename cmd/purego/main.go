package main

import (
	"fmt"
	"runtime"

	"github.com/badu/gxui/drivers/purego"
)

func main() {
	// Ensure we're on the main thread for GLFW
	runtime.LockOSThread()

	// Load GLFW library
	if err := purego.LoadGLFW(); err != nil {
		panic(err)
	}

	// Initialize GLFW
	if purego.Init() != purego.GLFW_TRUE {
		panic("failed to initialize GLFW")
	}
	defer purego.Terminate()

	// Set window hints
	purego.DefaultWindowHints()
	purego.WindowHint(0x00022001, purego.GLFW_FALSE) // GLFW_RESIZABLE = false

	// Create window
	window := purego.CreateWindow(800, 600, "Purego GLFW Window", nil, nil)
	if window == nil {
		panic("failed to create window")
	}

	fmt.Println("GLFW window created successfully!")
	fmt.Println("Window handle:", window.Handle())
	fmt.Println("Press Ctrl+C to exit...")

	// Make context current
	window.MakeContextCurrent()

	// Example: Get cursor position
	x, y := window.GetCursorPos()
	fmt.Printf("Initial cursor position: %.2f, %.2f\n", x, y)

	// Set cursor position callback
	window.SetCursorPosCallback(func(w *purego.Window, xpos, ypos float64) {
		fmt.Printf("Cursor moved to: %.2f, %.2f\n", xpos, ypos)
	})

	// Set key callback
	window.SetKeyCallback(func(w *purego.Window, key purego.Key, scancode int32, action purego.Action, mods purego.ModifierKey) {
		fmt.Printf("Key event - key: %d, scancode: %d, action: %d, mods: %d\n", key, scancode, action, mods)
	})

	// Set character callback
	window.SetCharCallback(func(w *purego.Window, char rune) {
		fmt.Printf("Character input: %c (U+%04X)\n", char, char)
	})

	// Set scroll callback
	window.SetScrollCallback(func(w *purego.Window, xoff, yoff float64) {
		fmt.Printf("Scroll offset: X=%.2f, Y=%.2f\n", xoff, yoff)
	})

	// Set mouse button callback
	window.SetMouseButtonCallback(func(w *purego.Window, button purego.MouseButton, action purego.Action, mods purego.ModifierKey) {
		fmt.Printf("Mouse button event - button: %d, action: %d, mods: %d\n", button, action, mods)
	})

	// Set framebuffer size callback
	window.SetFramebufferSizeCallback(func(w *purego.Window, width, height int32) {
		fmt.Printf("Framebuffer resized to: %dx%d\n", width, height)
	})

	// Set close callback
	shouldClose := false
	window.SetCloseCallback(func(w *purego.Window) {
		fmt.Println("Window close requested!")
		shouldClose = true
	})

	// Set refresh callback
	window.SetRefreshCallback(func(w *purego.Window) {
		fmt.Println("Window needs to be redrawn")
		// In a real app, this is where you'd redraw the window content
	})

	// Set size callback
	window.SetSizeCallback(func(w *purego.Window, width, height int32) {
		fmt.Printf("Window resized to: %dx%d\n", width, height)
	})

	// Set cursor enter callback
	window.SetCursorEnterCallback(func(w *purego.Window, entered bool) {
		if entered {
			fmt.Println("Cursor entered window")
		} else {
			fmt.Println("Cursor left window")
		}
	})

	// Set character with modifiers callback
	window.SetCharModsCallback(func(w *purego.Window, char rune, mods purego.ModifierKey) {
		fmt.Printf("Character with mods: %c (U+%04X), mods: %d\n", char, char, mods)
	})

	// Set window position callback
	window.SetPosCallback(func(w *purego.Window, xpos, ypos int32) {
		fmt.Printf("Window moved to position: (%d, %d)\n", xpos, ypos)
	})

	// Set window focus callback
	window.SetFocusCallback(func(w *purego.Window, focused bool) {
		if focused {
			fmt.Println("Window gained focus")
		} else {
			fmt.Println("Window lost focus")
		}
	})

	// Set window iconify callback
	window.SetIconifyCallback(func(w *purego.Window, iconified bool) {
		if iconified {
			fmt.Println("Window minimized")
		} else {
			fmt.Println("Window restored")
		}
	})

	// Set file drop callback
	window.SetDropCallback(func(w *purego.Window, paths []string) {
		fmt.Printf("Files dropped: %d\n", len(paths))
		for i, path := range paths {
			fmt.Printf("  [%d] %s\n", i, path)
		}
	})

	// Example: Set clipboard string
	window.SetClipboardString("Hello from GLFW!")
	fmt.Println("Clipboard set to: Hello from GLFW!")

	// Get primary monitor video mode
	monitor := purego.GetPrimaryMonitor()
	if monitor != nil {
		videoMode := monitor.GetVideoMode()
		if videoMode != nil {
			fmt.Printf("Monitor resolution: %dx%d @ %dHz\n", videoMode.Width, videoMode.Height, videoMode.RefreshRate)
		}

		wx, wy, ww, wh := purego.GetMonitorWorkarea(monitor)
		fmt.Printf("Monitor work area: position=(%d, %d), size=%dx%d\n", wx, wy, ww, wh)
	}

	// Get window size
	width, height := purego.GetWindowSize(window)
	fmt.Printf("Window size: %dx%d\n", width, height)

	// Get framebuffer size (in pixels, important for OpenGL viewport)
	fbWidth, fbHeight := window.GetFramebufferSize()
	fmt.Printf("Framebuffer size: %dx%d\n", fbWidth, fbHeight)

	// Calculate pixel ratio (important for high-DPI displays)
	pixelRatio := float64(fbWidth) / float64(width)
	fmt.Printf("Pixel ratio: %.2f\n", pixelRatio)

	// Example: Window manipulation functions

	// Get and set window position
	xp, yp := purego.GetWindowPos(window)
	fmt.Printf("Window position: (%d, %d)\n", xp, yp)

	// Move window to center of screen (assuming 1920x1080 screen)
	purego.SetWindowPos(window, 960-400, 540-300)

	// Change window size
	purego.SetWindowSize(window, 1024, 768)
	fmt.Println("Window resized to 1024x768")

	// Update window title
	purego.SetWindowTitle(window, "Purego GLFW - Updated Title")

	// Example: Set input modes

	// Disable cursor for FPS-style camera control
	window.SetInputMode(purego.GLFW_CURSOR, purego.GLFW_CURSOR_DISABLED)
	fmt.Println("Cursor disabled for FPS camera mode")

	// Enable sticky keys to avoid missing key presses
	window.SetInputMode(purego.GLFW_STICKY_KEYS, purego.GLFW_TRUE)
	fmt.Println("Sticky keys enabled")

	// Example: Poll key state in event loop

	// Event loop
	for !shouldClose {
		purego.PollEvents()

		// Check if Escape key is pressed
		if window.GetKey(purego.KeyEscape) == purego.GLFW_PRESS {
			fmt.Println("Escape key is pressed!")
			shouldClose = true
		}

		// Check if left mouse button is pressed
		if window.GetMouseButton(purego.GLFW_MOUSE_BUTTON_LEFT) == purego.GLFW_PRESS {
			x, y := window.GetCursorPos()
			fmt.Printf("Left mouse button pressed at: %.2f, %.2f\n", x, y)
		}
	}

	fmt.Println("Exiting...")
}
