package main

import (
	"fmt"

	"github.com/badu/gxui/drivers/purego"
)

func main() {
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
	window.SetCloseCallback(func(w *purego.Window) {
		fmt.Println("Window close requested!")
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

		// Get monitor work area (excludes taskbars, menu bars, etc.)
		wx, wy, ww, wh := monitor.GetWorkarea()
		fmt.Printf("Monitor work area: position=(%d, %d), size=%dx%d\n", wx, wy, ww, wh)
	}

	// Get window size
	width, height := window.GetWindowSize()
	fmt.Printf("Window size: %dx%d\n", width, height)

	// Get framebuffer size (in pixels, important for OpenGL viewport)
	fbWidth, fbHeight := window.GetFramebufferSize()
	fmt.Printf("Framebuffer size: %dx%d\n", fbWidth, fbHeight)

	// Calculate pixel ratio (important for high-DPI displays)
	pixelRatio := float64(fbWidth) / float64(width)
	fmt.Printf("Pixel ratio: %.2f\n", pixelRatio)

	// Example: Window manipulation functions

	// Get and set window position
	xw, yw := window.GetWindowPos()
	fmt.Printf("Window position: (%d, %d)\n", xw, yw)

	// Move window to center of screen (assuming 1920x1080 screen)
	window.SetWindowPos(960-400, 540-300)

	// Change window size
	window.SetWindowSize(1024, 768)
	fmt.Println("Window resized to 1024x768")

	// Update window title
	window.SetWindowTitle("Purego GLFW - Updated Title")

	// Example: Set input modes
	const (
		GLFW_CURSOR           = 0x00033001
		GLFW_CURSOR_DISABLED  = 0x00034003
		GLFW_STICKY_KEYS      = 0x00033002
		GLFW_RAW_MOUSE_MOTION = 0x00033005
	)

	// Disable cursor for FPS-style camera control
	window.SetInputMode(purego.InputMode(GLFW_CURSOR), GLFW_CURSOR_DISABLED)
	fmt.Println("Cursor disabled for FPS camera mode")

	// Enable sticky keys to avoid missing key presses
	window.SetInputMode(purego.InputMode(GLFW_STICKY_KEYS), purego.GLFW_TRUE)
	fmt.Println("Sticky keys enabled")

	// Example: Poll key state in event loop
	const (
		GLFW_KEY_ESCAPE          = 256
		GLFW_KEY_M               = 77
		GLFW_KEY_I               = 73
		GLFW_KEY_R               = 82
		GLFW_PRESS               = 1
		GLFW_MOUSE_BUTTON_LEFT   = 0
		GLFW_MOUSE_BUTTON_RIGHT  = 1
		GLFW_MOUSE_BUTTON_MIDDLE = 2
	)

	// Example: Use GetTime for delta time calculation
	lastTime := purego.GetTime()
	frameCount := 0

	// Event loop using ShouldClose
	for !window.ShouldClose() {
		purego.PollEvents()

		// Calculate delta time
		currentTime := purego.GetTime()
		deltaTime := currentTime - lastTime
		lastTime = currentTime

		// Check if Escape key is pressed
		if window.GetKey(purego.Key(GLFW_KEY_ESCAPE)) == purego.Action(GLFW_PRESS) {
			fmt.Println("Escape key pressed - closing window")
			window.SetShouldClose(true)
		}

		// Check for window state hotkeys
		if window.GetKey(purego.Key(GLFW_KEY_M)) == purego.Action(GLFW_PRESS) {
			window.Maximize()
			fmt.Println("Window maximized")
		}
		if window.GetKey(purego.Key(GLFW_KEY_I)) == purego.Action(GLFW_PRESS) {
			window.Iconify()
			fmt.Println("Window iconified")
		}
		if window.GetKey(purego.Key(GLFW_KEY_R)) == purego.Action(GLFW_PRESS) {
			window.Restore()
			fmt.Println("Window restored")
		}

		// Check if left mouse button is pressed
		if window.GetMouseButton(purego.MouseButton(GLFW_MOUSE_BUTTON_LEFT)) == purego.Action(GLFW_PRESS) {
			x, y := window.GetCursorPos()
			fmt.Printf("Left mouse button pressed at: %.2f, %.2f\n", x, y)
		}

		// Example rendering (swap buffers)
		// In a real OpenGL app, you'd render here
		window.SwapBuffers()

		frameCount++
		if frameCount%60 == 0 {
			fps := 1.0 / deltaTime
			window.SetWindowTitle(fmt.Sprintf("Purego GLFW - Frame %d (%.1f FPS)", frameCount, fps))
		}
	}

	// Cleanup
	window.Destroy()
	fmt.Println("Window destroyed")
	fmt.Println("Exiting...")
}
