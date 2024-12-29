// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type KeyboardController struct {
	window Window
}

func CreateKeyboardController(window Window) *KeyboardController {
	result := &KeyboardController{
		window: window,
	}
	window.OnKeyDown(result.keyDown)
	window.OnKeyUp(result.keyUp)
	window.OnKeyRepeat(result.keyPress)
	window.OnKeyStroke(result.keyStroke)
	return result
}

func (c *KeyboardController) keyDown(event KeyboardEvent) {
	target := Control(c.window.Focus())
	for target != nil {
		target.KeyDown(event)
		target, _ = target.Parent().(Control)
	}
	c.keyPress(event)
}

func (c *KeyboardController) keyUp(event KeyboardEvent) {
	target := Control(c.window.Focus())
	for target != nil {
		target.KeyUp(event)
		target, _ = target.Parent().(Control)
	}
}

func (c *KeyboardController) keyPress(event KeyboardEvent) {
	target := Control(c.window.Focus())
	for target != nil {
		if target.KeyPress(event) {
			return
		}
		target, _ = target.Parent().(Control)
	}
	c.window.KeyPress(event)
}

func (c *KeyboardController) keyStroke(event KeyStrokeEvent) {
	target := Control(c.window.Focus())
	for target != nil {
		if target.KeyStroke(event) {
			return
		}
		target, _ = target.Parent().(Control)
	}
	c.window.KeyStroke(event)
}
