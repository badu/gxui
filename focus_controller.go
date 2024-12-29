// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type FocusController struct {
	window             Window
	focus              Focusable
	setFocusCount      int
	detachSubscription EventSubscription
}

func CreateFocusController(window Window) *FocusController {
	return &FocusController{window: window}
}

func (c *FocusController) SetFocus(target Focusable) {
	c.setFocusCount++
	if c.focus == target {
		return
	}

	if c.focus != nil {
		o := c.focus
		c.focus = nil
		c.detachSubscription.Unlisten()
		o.LostFocus()
		if c.focus != nil {
			return // Something in LostFocus() called SetFocus(). Respect their call.
		}
	}

	c.focus = target

	if c.focus != nil {
		c.detachSubscription = c.focus.OnDetach(func() { c.SetFocus(nil) })
		c.focus.GainedFocus()
	}
}

func (c *FocusController) SetFocusCount() int {
	return c.setFocusCount
}

func (c *FocusController) Focus() Focusable {
	return c.focus
}

func (c *FocusController) FocusNext() {
	c.SetFocus(c.NextFocusable(c.focus, true))
}

func (c *FocusController) FocusPrev() {
	c.SetFocus(c.NextFocusable(c.focus, false))
}

func (c *FocusController) NextFocusable(control Control, forwards bool) Focusable {
	container, _ := control.(Container)
	if container != nil {
		child := c.NextChildFocusable(container, nil, forwards)
		if child != nil {
			return child
		}
	}

	for control != nil {
		parent := control.Parent()
		if parent != nil {
			child := c.NextChildFocusable(parent, control, forwards)
			if child != nil {
				return child
			}
		}
		control, _ = parent.(Control)
	}

	return c.NextChildFocusable(c.window, nil, forwards)
}

func (c *FocusController) NextChildFocusable(parent Parent, control Control, forwards bool) Focusable {
	examineNext := control == nil
	children := parent.Children()

	index := 0
	numChildren := len(children)
	if !forwards {
		index = len(children) - 1
		numChildren = -1
	}

	for index != numChildren {
		child := children[index]
		if forwards {
			index++
		} else {
			index--
		}

		if !examineNext {
			if child.Control == control {
				examineNext = true
			}
			continue
		}

		if target := c.Focusable(child.Control); target != nil {
			return target
		}

		if container, ok := child.Control.(Container); ok {
			focusable := c.NextChildFocusable(container, nil, forwards)
			if focusable != nil {
				return focusable
			}
		}
	}
	return nil
}

func (c *FocusController) Focusable(control Control) Focusable {
	target, _ := control.(Focusable)
	if target != nil && target.IsFocusable() {
		return target
	}
	return nil
}
