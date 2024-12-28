// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/badu/gxui"

	"github.com/goxjs/glfw"
)

func translateMouseButton(button glfw.MouseButton) gxui.MouseButton {
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

func getMouseState(glfwWindow *glfw.Window) gxui.MouseState {
	var state gxui.MouseState
	for _, button := range []glfw.MouseButton{glfw.MouseButtonLeft, glfw.MouseButtonMiddle, glfw.MouseButtonRight} {
		if glfwWindow.GetMouseButton(button) == glfw.Press {
			state |= 1 << uint(translateMouseButton(button))
		}
	}
	return state
}
