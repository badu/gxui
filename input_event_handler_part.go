// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type InputEventHandlerPart struct {
	onClick       Event
	onDoubleClick Event
	onKeyPress    Event
	onKeyStroke   Event
	onMouseMove   Event
	onMouseEnter  Event
	onMouseExit   Event
	onMouseDown   Event
	onMouseUp     Event
	onMouseScroll Event
	onKeyDown     Event
	onKeyUp       Event
	onKeyRepeat   Event
	isMouseDown   map[MouseButton]bool
	isMouseOver   bool
}

func (m *InputEventHandlerPart) getOnClick() Event {
	if m.onClick == nil {
		m.onClick = NewListener(m.Click)
	}

	return m.onClick
}

func (m *InputEventHandlerPart) getOnDoubleClick() Event {
	if m.onDoubleClick == nil {
		m.onDoubleClick = NewListener(m.DoubleClick)
	}

	return m.onDoubleClick
}

func (m *InputEventHandlerPart) getOnKeyPress() Event {
	if m.onKeyPress == nil {
		m.onKeyPress = NewListener(m.KeyPress)
	}

	return m.onKeyPress
}

func (m *InputEventHandlerPart) getOnKeyStroke() Event {
	if m.onKeyStroke == nil {
		m.onKeyStroke = NewListener(m.KeyStroke)
	}

	return m.onKeyStroke
}

func (m *InputEventHandlerPart) getOnMouseMove() Event {
	if m.onMouseMove == nil {
		m.onMouseMove = NewListener(m.MouseMove)
	}

	return m.onMouseMove
}

func (m *InputEventHandlerPart) getOnMouseEnter() Event {
	if m.onMouseEnter == nil {
		m.onMouseEnter = NewListener(m.MouseEnter)
	}

	return m.onMouseEnter
}

func (m *InputEventHandlerPart) getOnMouseExit() Event {
	if m.onMouseExit == nil {
		m.onMouseExit = NewListener(m.MouseExit)
	}

	return m.onMouseExit
}

func (m *InputEventHandlerPart) getOnMouseDown() Event {
	if m.onMouseDown == nil {
		m.onMouseDown = NewListener(m.MouseDown)
	}

	return m.onMouseDown
}

func (m *InputEventHandlerPart) getOnMouseUp() Event {
	if m.onMouseUp == nil {
		m.onMouseUp = NewListener(m.MouseUp)
	}

	return m.onMouseUp
}

func (m *InputEventHandlerPart) getOnMouseScroll() Event {
	if m.onMouseScroll == nil {
		m.onMouseScroll = NewListener(m.MouseScroll)
	}

	return m.onMouseScroll
}

func (m *InputEventHandlerPart) getOnKeyDown() Event {
	if m.onKeyDown == nil {
		m.onKeyDown = NewListener(m.KeyDown)
	}

	return m.onKeyDown
}

func (m *InputEventHandlerPart) getOnKeyUp() Event {
	if m.onKeyUp == nil {
		m.onKeyUp = NewListener(m.KeyUp)
	}

	return m.onKeyUp
}

func (m *InputEventHandlerPart) getOnKeyRepeat() Event {
	if m.onKeyRepeat == nil {
		m.onKeyRepeat = NewListener(m.KeyRepeat)
	}

	return m.onKeyRepeat
}

func (m *InputEventHandlerPart) Init() {
	m.isMouseDown = make(map[MouseButton]bool)
}

func (m *InputEventHandlerPart) Click(ev MouseEvent) bool {
	m.getOnClick().Emit(ev)
	return false
}

func (m *InputEventHandlerPart) DoubleClick(ev MouseEvent) bool {
	m.getOnDoubleClick().Emit(ev)
	return false
}

func (m *InputEventHandlerPart) KeyPress(ev KeyboardEvent) bool {
	m.getOnKeyPress().Emit(ev)
	return false
}

func (m *InputEventHandlerPart) KeyStroke(ev KeyStrokeEvent) bool {
	m.getOnKeyStroke().Emit(ev)
	return false
}

func (m *InputEventHandlerPart) MouseScroll(ev MouseEvent) bool {
	m.getOnMouseScroll().Emit(ev)
	return false
}

func (m *InputEventHandlerPart) MouseMove(ev MouseEvent) {
	m.getOnMouseMove().Emit(ev)
}

func (m *InputEventHandlerPart) MouseEnter(ev MouseEvent) {
	m.isMouseOver = true
	m.getOnMouseEnter().Emit(ev)
}

func (m *InputEventHandlerPart) MouseExit(ev MouseEvent) {
	m.isMouseOver = false
	m.getOnMouseExit().Emit(ev)
}

func (m *InputEventHandlerPart) MouseDown(ev MouseEvent) {
	m.isMouseDown[ev.Button] = true
	m.getOnMouseDown().Emit(ev)
}

func (m *InputEventHandlerPart) MouseUp(ev MouseEvent) {
	m.isMouseDown[ev.Button] = false
	m.getOnMouseUp().Emit(ev)
}

func (m *InputEventHandlerPart) KeyDown(ev KeyboardEvent) {
	m.getOnKeyDown().Emit(ev)
}

func (m *InputEventHandlerPart) KeyUp(ev KeyboardEvent) {
	m.getOnKeyUp().Emit(ev)
}

func (m *InputEventHandlerPart) KeyRepeat(ev KeyboardEvent) {
	m.getOnKeyRepeat().Emit(ev)
}

func (m *InputEventHandlerPart) OnClick(callback func(MouseEvent)) EventSubscription {
	return m.getOnClick().Listen(callback)
}

func (m *InputEventHandlerPart) OnDoubleClick(callback func(MouseEvent)) EventSubscription {
	return m.getOnDoubleClick().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyPress(callback func(KeyboardEvent)) EventSubscription {
	return m.getOnKeyPress().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription {
	return m.getOnKeyStroke().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseMove(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseMove().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseEnter(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseEnter().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseExit(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseExit().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseDown(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseDown().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseUp(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseUp().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseScroll(callback func(MouseEvent)) EventSubscription {
	return m.getOnMouseScroll().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyDown(callback func(KeyboardEvent)) EventSubscription {
	return m.getOnKeyDown().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyUp(callback func(KeyboardEvent)) EventSubscription {
	return m.getOnKeyUp().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription {
	return m.getOnKeyRepeat().Listen(callback)
}

func (m *InputEventHandlerPart) IsMouseOver() bool {
	return m.isMouseOver
}

func (m *InputEventHandlerPart) IsMouseDown(button MouseButton) bool {
	return m.isMouseDown[button]
}
