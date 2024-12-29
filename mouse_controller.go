// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"time"
)

var doubleClickTime = time.Millisecond * 300

type MouseController struct {
	window          Window
	focusController *FocusController
	lastOver        ControlPointList
	lastDown        map[MouseButton]ControlPointList
	lastUpTime      map[MouseButton]time.Time
}

func CreateMouseController(window Window, focusCtrl *FocusController) *MouseController {
	result := &MouseController{
		window:          window,
		focusController: focusCtrl,
		lastDown:        make(map[MouseButton]ControlPointList),
		lastUpTime:      make(map[MouseButton]time.Time),
	}
	window.OnMouseMove(result.mouseMove)
	window.OnMouseEnter(result.mouseMove)
	window.OnMouseExit(result.mouseMove)
	window.OnMouseDown(result.mouseDown)
	window.OnMouseUp(result.mouseUp)
	window.OnMouseScroll(result.mouseScroll)
	return result
}

func (m *MouseController) updatePosition(event MouseEvent) {
	ValidateHierarchy(m.window)

	nowOver := TopControlsUnder(event.Point, m.window)

	for _, point := range m.lastOver {
		if !nowOver.Contains(point.Control) {
			e := event
			e.Point = point.Point
			point.Control.MouseExit(e)
		}
	}

	for _, point := range nowOver {
		if !m.lastOver.Contains(point.Control) {
			e := event
			e.Point = point.Point
			point.Control.MouseEnter(e)
		}
	}

	m.lastOver = nowOver
}

func (m *MouseController) mouseMove(event MouseEvent) {
	m.updatePosition(event)
	for _, point := range m.lastOver {
		e := event
		e.Point = point.Point
		point.Control.MouseMove(e)
	}
}

func (m *MouseController) mouseDown(event MouseEvent) {
	m.updatePosition(event)

	for _, point := range m.lastOver {
		e := event
		e.Point = point.Point
		point.Control.MouseDown(e)
	}

	m.lastDown[event.Button] = m.lastOver
}

func (m *MouseController) mouseUp(event MouseEvent) {
	m.updatePosition(event)

	for _, point := range m.lastDown[event.Button] {
		e := event
		e.Point = point.Point
		point.Control.MouseUp(e)
	}

	setFocusCount := m.focusController.SetFocusCount()

	dblClick := time.Since(m.lastUpTime[event.Button]) < doubleClickTime
	clickConsumed := false
	for i := len(m.lastDown[event.Button]) - 1; i >= 0; i-- {
		point := m.lastDown[event.Button][i]
		if p, found := m.lastOver.Find(point.Control); found {
			event.Point = p
			if (dblClick && point.Control.DoubleClick(event)) || (!dblClick && point.Control.Click(event)) {
				clickConsumed = true
				break
			}
		}
	}

	if !clickConsumed {
		event.Point = event.WindowPoint
		if dblClick {
			m.window.DoubleClick(event)
		} else {
			m.window.Click(event)
		}
	}

	focusSet := setFocusCount != m.focusController.SetFocusCount()
	if !focusSet {
		for i := len(m.lastDown[event.Button]) - 1; i >= 0; i-- {
			point := m.lastDown[event.Button][i]
			if m.lastOver.Contains(point.Control) && m.window.SetFocus(point.Control) {
				focusSet = true
				break
			}
		}

		if !focusSet {
			m.window.SetFocus(nil)
		}
	}

	delete(m.lastDown, event.Button)
	m.lastUpTime[event.Button] = time.Now()
}

func (m *MouseController) mouseScroll(event MouseEvent) {
	m.updatePosition(event)

	for i := len(m.lastOver) - 1; i >= 0; i-- {
		point := m.lastOver[i]
		e := event
		e.Point = point.Point
		point.Control.MouseScroll(e)
	}
}
