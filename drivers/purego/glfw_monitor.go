package purego

type VidMode struct {
	Width       int // The width, in pixels, of the video mode.
	Height      int // The height, in pixels, of the video mode.
	RedBits     int // The bit depth of the red channel of the video mode.
	GreenBits   int // The bit depth of the green channel of the video mode.
	BlueBits    int // The bit depth of the blue channel of the video mode.
	RefreshRate int // The refresh rate, in Hz, of the video mode.
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

func (m *Monitor) GetVideoMode() *VidMode {
	return GetVideoMode(m)
}

func (m *Monitor) GetWorkarea() (xpos, ypos, width, height int) {
	return GetMonitorWorkarea(m)
}
