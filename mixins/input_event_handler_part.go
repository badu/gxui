// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type InputEventHandlerOuter interface{}

type InputEventHandlerPart struct {
	outer         InputEventHandlerOuter
	isMouseOver   bool
	isMouseDown   map[gxui.MouseButton]bool
	onClick       gxui.Event
	onDoubleClick gxui.Event
	onKeyPress    gxui.Event
	onKeyStroke   gxui.Event
	onMouseMove   gxui.Event
	onMouseEnter  gxui.Event
	onMouseExit   gxui.Event
	onMouseDown   gxui.Event
	onMouseUp     gxui.Event
	onMouseScroll gxui.Event
	onKeyDown     gxui.Event
	onKeyUp       gxui.Event
	onKeyRepeat   gxui.Event
}

func (m *InputEventHandlerPart) getOnClick() gxui.Event {
	if m.onClick == nil {
		m.onClick = gxui.CreateEvent(m.Click)
	}
	return m.onClick
}

func (m *InputEventHandlerPart) getOnDoubleClick() gxui.Event {
	if m.onDoubleClick == nil {
		m.onDoubleClick = gxui.CreateEvent(m.DoubleClick)
	}
	return m.onDoubleClick
}

func (m *InputEventHandlerPart) getOnKeyPress() gxui.Event {
	if m.onKeyPress == nil {
		m.onKeyPress = gxui.CreateEvent(m.KeyPress)
	}
	return m.onKeyPress
}

func (m *InputEventHandlerPart) getOnKeyStroke() gxui.Event {
	if m.onKeyStroke == nil {
		m.onKeyStroke = gxui.CreateEvent(m.KeyStroke)
	}
	return m.onKeyStroke
}

func (m *InputEventHandlerPart) getOnMouseMove() gxui.Event {
	if m.onMouseMove == nil {
		m.onMouseMove = gxui.CreateEvent(m.MouseMove)
	}
	return m.onMouseMove
}

func (m *InputEventHandlerPart) getOnMouseEnter() gxui.Event {
	if m.onMouseEnter == nil {
		m.onMouseEnter = gxui.CreateEvent(m.MouseEnter)
	}
	return m.onMouseEnter
}

func (m *InputEventHandlerPart) getOnMouseExit() gxui.Event {
	if m.onMouseExit == nil {
		m.onMouseExit = gxui.CreateEvent(m.MouseExit)
	}
	return m.onMouseExit
}

func (m *InputEventHandlerPart) getOnMouseDown() gxui.Event {
	if m.onMouseDown == nil {
		m.onMouseDown = gxui.CreateEvent(m.MouseDown)
	}
	return m.onMouseDown
}

func (m *InputEventHandlerPart) getOnMouseUp() gxui.Event {
	if m.onMouseUp == nil {
		m.onMouseUp = gxui.CreateEvent(m.MouseUp)
	}
	return m.onMouseUp
}

func (m *InputEventHandlerPart) getOnMouseScroll() gxui.Event {
	if m.onMouseScroll == nil {
		m.onMouseScroll = gxui.CreateEvent(m.MouseScroll)
	}
	return m.onMouseScroll
}

func (m *InputEventHandlerPart) getOnKeyDown() gxui.Event {
	if m.onKeyDown == nil {
		m.onKeyDown = gxui.CreateEvent(m.KeyDown)
	}
	return m.onKeyDown
}

func (m *InputEventHandlerPart) getOnKeyUp() gxui.Event {
	if m.onKeyUp == nil {
		m.onKeyUp = gxui.CreateEvent(m.KeyUp)
	}
	return m.onKeyUp
}

func (m *InputEventHandlerPart) getOnKeyRepeat() gxui.Event {
	if m.onKeyRepeat == nil {
		m.onKeyRepeat = gxui.CreateEvent(m.KeyRepeat)
	}
	return m.onKeyRepeat
}

func (m *InputEventHandlerPart) Init(outer InputEventHandlerOuter) {
	m.outer = outer
	m.isMouseDown = make(map[gxui.MouseButton]bool)
}

func (m *InputEventHandlerPart) Click(ev gxui.MouseEvent) (consume bool) {
	m.getOnClick().Fire(ev)
	return false
}

func (m *InputEventHandlerPart) DoubleClick(ev gxui.MouseEvent) (consume bool) {
	m.getOnDoubleClick().Fire(ev)
	return false
}

func (m *InputEventHandlerPart) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	m.getOnKeyPress().Fire(ev)
	return false
}

func (m *InputEventHandlerPart) KeyStroke(ev gxui.KeyStrokeEvent) (consume bool) {
	m.getOnKeyStroke().Fire(ev)
	return false
}

func (m *InputEventHandlerPart) MouseScroll(ev gxui.MouseEvent) (consume bool) {
	m.getOnMouseScroll().Fire(ev)
	return false
}

func (m *InputEventHandlerPart) MouseMove(ev gxui.MouseEvent) {
	m.getOnMouseMove().Fire(ev)
}

func (m *InputEventHandlerPart) MouseEnter(ev gxui.MouseEvent) {
	m.isMouseOver = true
	m.getOnMouseEnter().Fire(ev)
}

func (m *InputEventHandlerPart) MouseExit(ev gxui.MouseEvent) {
	m.isMouseOver = false
	m.getOnMouseExit().Fire(ev)
}

func (m *InputEventHandlerPart) MouseDown(ev gxui.MouseEvent) {
	m.isMouseDown[ev.Button] = true
	m.getOnMouseDown().Fire(ev)
}

func (m *InputEventHandlerPart) MouseUp(ev gxui.MouseEvent) {
	m.isMouseDown[ev.Button] = false
	m.getOnMouseUp().Fire(ev)
}

func (m *InputEventHandlerPart) KeyDown(ev gxui.KeyboardEvent) {
	m.getOnKeyDown().Fire(ev)
}

func (m *InputEventHandlerPart) KeyUp(ev gxui.KeyboardEvent) {
	m.getOnKeyUp().Fire(ev)
}

func (m *InputEventHandlerPart) KeyRepeat(ev gxui.KeyboardEvent) {
	m.getOnKeyRepeat().Fire(ev)
}

func (m *InputEventHandlerPart) OnClick(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnClick().Listen(callback)
}

func (m *InputEventHandlerPart) OnDoubleClick(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnDoubleClick().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyPress(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyPress().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyStroke(callback func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return m.getOnKeyStroke().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseMove(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseMove().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseEnter(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseEnter().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseExit(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseExit().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseDown(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseDown().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseUp(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseUp().Listen(callback)
}

func (m *InputEventHandlerPart) OnMouseScroll(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseScroll().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyDown(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyDown().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyUp(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyUp().Listen(callback)
}

func (m *InputEventHandlerPart) OnKeyRepeat(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyRepeat().Listen(callback)
}

func (m *InputEventHandlerPart) IsMouseOver() bool {
	return m.isMouseOver
}

func (m *InputEventHandlerPart) IsMouseDown(button gxui.MouseButton) bool {
	return m.isMouseDown[button]
}
