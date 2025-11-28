package purego

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func translateMouseButton(button MouseButton) gxui.MouseButton {
	switch button {
	case glfw.MouseButtonLeft:
		return gxui.MouseButtonLeft
	case glfw.MouseButtonMiddle:
		return gxui.MouseButtonMiddle
	case glfw.MouseButtonRight:
		return gxui.MouseButtonRight
	default:
		panic(fmt.Errorf("unknown mouse button %v", button))
	}
}

func getMouseState(glfwWindow *Window) gxui.MouseState {
	var state gxui.MouseState
	for _, button := range []MouseButton{glfw.MouseButtonLeft, glfw.MouseButtonMiddle, glfw.MouseButtonRight} {
		if glfwWindow.GetMouseButton(button) == glfw.Press {
			state |= 1 << uint(translateMouseButton(button))
		}
	}
	return state
}
