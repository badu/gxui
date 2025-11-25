package purego

import (
	"fmt"

	"github.com/badu/gxui"
)

func translateMouseButton(button MouseButton) gxui.MouseButton {
	switch button {
	case MouseButtonLeft:
		return gxui.MouseButtonLeft
	case MouseButtonMiddle:
		return gxui.MouseButtonMiddle
	case MouseButtonRight:
		return gxui.MouseButtonRight
	default:
		panic(fmt.Errorf("unknown mouse button %v", button))
	}
}

func getMouseState(glfwWindow *Window) gxui.MouseState {
	var state gxui.MouseState
	for _, button := range []MouseButton{MouseButtonLeft, MouseButtonMiddle, MouseButtonRight} {
		if glfwWindow.GetMouseButton(button) == Press {
			state |= 1 << uint(translateMouseButton(button))
		}
	}
	return state
}
