package purego

import (
	"fmt"

	"github.com/badu/gxui"
)

func translateMouseButton(button MouseButton) gxui.MouseButton {
	switch button {
	case GLFW_MOUSE_BUTTON_LEFT:
		return gxui.MouseButtonLeft
	case GLFW_MOUSE_BUTTON_MIDDLE:
		return gxui.MouseButtonMiddle
	case GLFW_MOUSE_BUTTON_RIGHT:
		return gxui.MouseButtonRight
	default:
		panic(fmt.Errorf("unknown mouse button %v", button))
	}
}

func getMouseState(glfwWindow *Window) gxui.MouseState {
	var state gxui.MouseState
	for _, button := range []MouseButton{GLFW_MOUSE_BUTTON_LEFT, GLFW_MOUSE_BUTTON_MIDDLE, GLFW_MOUSE_BUTTON_RIGHT} {
		if glfwWindow.GetMouseButton(button) == Press {
			state |= 1 << uint(translateMouseButton(button))
		}
	}
	return state
}
