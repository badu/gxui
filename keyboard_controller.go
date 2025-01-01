// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type KeyboardKey int

const (
	KeyUnknown KeyboardKey = iota
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyKp0
	KeyKp1
	KeyKp2
	KeyKp3
	KeyKp4
	KeyKp5
	KeyKp6
	KeyKp7
	KeyKp8
	KeyKp9
	KeyKpDecimal
	KeyKpDivide
	KeyKpMultiply
	KeyKpSubtract
	KeyKpAdd
	KeyKpEnter
	KeyKpEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
	KeyLast
)

type KeyboardModifier int

const (
	ModNone    KeyboardModifier = 0
	ModShift   KeyboardModifier = 1
	ModControl KeyboardModifier = 2
	ModAlt     KeyboardModifier = 4
	ModSuper   KeyboardModifier = 8
)

func (m KeyboardModifier) Shift() bool {
	return m&ModShift != 0
}

func (m KeyboardModifier) Control() bool {
	return m&ModControl != 0
}

func (m KeyboardModifier) Alt() bool {
	return m&ModAlt != 0
}

func (m KeyboardModifier) Super() bool {
	return m&ModSuper != 0
}

type KeyStrokeEvent struct {
	Character rune
	Modifier  KeyboardModifier
}

type KeyboardEvent struct {
	Key      KeyboardKey
	Modifier KeyboardModifier
}

type KeyboardController struct {
	window *WindowImpl
}

func CreateKeyboardController(window *WindowImpl) *KeyboardController {
	result := &KeyboardController{window: window}
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
