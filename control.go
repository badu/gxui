// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/pkg/math"

type Control interface {
	Size() math.Size
	SetSize(newSize math.Size)

	Draw() Canvas

	Parent() Parent
	SetParent(newParent Parent)

	Attached() bool
	Attach()
	Detach()
	OnAttach(callback func()) EventSubscription
	OnDetach(callback func()) EventSubscription

	DesiredSize(min, max math.Size) math.Size

	Margin() math.Spacing
	SetMargin(math.Spacing)

	IsVisible() bool
	SetVisible(isVisible bool)

	ContainsPoint(point math.Point) bool

	IsMouseOver() bool
	IsMouseDown(button MouseButton) bool
	Click(event MouseEvent) (consume bool)
	DoubleClick(event MouseEvent) (consume bool)
	OnClick(callback func(MouseEvent)) EventSubscription
	OnDoubleClick(callback func(MouseEvent)) EventSubscription
	OnMouseMove(callback func(MouseEvent)) EventSubscription
	OnMouseEnter(callback func(MouseEvent)) EventSubscription
	OnMouseExit(callback func(MouseEvent)) EventSubscription
	OnMouseDown(callback func(MouseEvent)) EventSubscription
	OnMouseUp(callback func(MouseEvent)) EventSubscription
	OnMouseScroll(callback func(MouseEvent)) EventSubscription
	MouseScroll(event MouseEvent) (consume bool)
	MouseMove(event MouseEvent)
	MouseEnter(event MouseEvent)
	MouseExit(event MouseEvent)
	MouseDown(event MouseEvent)
	MouseUp(event MouseEvent)

	KeyPress(event KeyboardEvent) (consume bool)
	KeyStroke(event KeyStrokeEvent) (consume bool)
	KeyDown(event KeyboardEvent)
	KeyUp(event KeyboardEvent)
	KeyRepeat(event KeyboardEvent)
	OnKeyPress(callback func(KeyboardEvent)) EventSubscription
	OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription
	OnKeyDown(callback func(KeyboardEvent)) EventSubscription
	OnKeyUp(callback func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription
}

type ControlBaseParent interface {
	Control
	Paint(canvas Canvas)
	Redraw()
	ReLayout()
}

type ControlBase struct {
	InputEventHandlerPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	LayoutablePart
}

func (c *ControlBase) Init(parent ControlBaseParent, driver Driver) {
	c.DrawPaintPart.Init(parent, driver)
	c.LayoutablePart.Init(parent)
	c.InputEventHandlerPart.Init()
	c.ParentablePart.Init()
	c.VisiblePart.Init(parent)
}

func (c *ControlBase) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (c *ControlBase) ContainsPoint(point math.Point) bool {
	return c.IsVisible() && c.Size().Rect().Contains(point)
}
