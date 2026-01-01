// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"strings"

	"github.com/badu/gxui/pkg/math"
)

type Child struct {
	Control Control
	Offset  math.Point
}

type Parent interface {
	Children() Children
	ReLayout()
	Redraw()
}

type Container interface {
	// Parent interface
	Children() Children
	ReLayout()
	Redraw()
	// Container interface
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
}

type ContainerBaseNoControlOuter interface {
	// Parent interface
	Children() Children
	ReLayout()
	Redraw()
	// Container interface
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
	// ContainerBaseNoControlOuter interface
	PaintChild(canvas Canvas, child *Child, idx int) // was outer.PaintChilder
	Paint(canvas Canvas)                             // was outer.Painter
	LayoutChildren()                                 // was outer.LayoutChildren
}

type BaseContainerParent interface {
	// Parent interface
	Children() Children
	ReLayout()
	Redraw()
	// Container interface
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
	// Control interface
	Size() math.Size
	SetSize(newSize math.Size)
	Draw() Canvas
	Parent() Parent
	SetParent(newParent Parent)
	Attached() bool
	Attach()
	Detach()
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
	KeyPress(event KeyboardEvent) (consume bool)
	KeyStroke(event KeyStrokeEvent) (consume bool)
	MouseScroll(event MouseEvent) (consume bool)
	MouseMove(event MouseEvent)
	MouseEnter(event MouseEvent)
	MouseExit(event MouseEvent)
	MouseDown(event MouseEvent)
	MouseUp(event MouseEvent)
	KeyDown(event KeyboardEvent)
	KeyUp(event KeyboardEvent)
	KeyRepeat(event KeyboardEvent)
	OnAttach(callback func()) EventSubscription
	OnDetach(callback func()) EventSubscription
	OnKeyPress(callback func(KeyboardEvent)) EventSubscription
	OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription
	OnClick(callback func(MouseEvent)) EventSubscription
	OnDoubleClick(callback func(MouseEvent)) EventSubscription
	OnMouseMove(callback func(MouseEvent)) EventSubscription
	OnMouseEnter(callback func(MouseEvent)) EventSubscription
	OnMouseExit(callback func(MouseEvent)) EventSubscription
	OnMouseDown(callback func(MouseEvent)) EventSubscription
	OnMouseUp(callback func(MouseEvent)) EventSubscription
	OnMouseScroll(callback func(MouseEvent)) EventSubscription
	OnKeyDown(callback func(KeyboardEvent)) EventSubscription
	OnKeyUp(callback func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription
	// BaseContainerParent
	PaintChild(canvas Canvas, child *Child, idx int) // was outer.PaintChilder
	Paint(canvas Canvas)                             // was outer.Painter
	LayoutChildren()                                 // was outer.LayoutChildren
}

// String returns a string describing the child type and bounds.
func (c *Child) String() string {
	return fmt.Sprintf("Type: %T, Bounds: %v", c.Control, c.Bounds())
}

// Bounds returns the Child bounds relative to the parent.
func (c *Child) Bounds() math.Rect {
	return c.Control.Size().Rect().Offset(c.Offset)
}

// Layout sets the Child size and offset relative to the parent.
// Layout should only be called by the Child's parent.
func (c *Child) Layout(rect math.Rect) {
	c.Offset = rect.Min
	c.Control.SetSize(rect.Size())
}

// Children is a list of Child pointers.
type Children []*Child

// String returns a string describing the child type and bounds.
func (c Children) String() string {
	s := make([]string, len(c))
	for i, c := range c {
		s[i] = fmt.Sprintf("%d: %s", i, c.String())
	}
	return strings.Join(s, "\n")
}

// IndexOf returns and returns the index of the child control, or -1 if the
// child is not in this Children list.
func (c Children) IndexOf(control Control) int {
	for i, child := range c {
		if child.Control == control {
			return i
		}
	}
	return -1
}

// Find returns and returns the Child pointer for the given Control, or nil
// if the child is not in this Children list.
func (c Children) Find(control Control) *Child {
	for _, child := range c {
		if child.Control == control {
			return child
		}
	}
	return nil
}

type ContainerPartParent interface {
	// Parent interface
	Children() Children
	ReLayout()
	Redraw()

	AddChildAt(index int, child Control) *Child
	RemoveChildAt(index int)

	Attached() bool
	OnAttach(callback func()) EventSubscription
	OnDetach(callback func()) EventSubscription
	IsVisible() bool
	Size() math.Size
}

type ContainerPart struct {
	parent             ContainerPartParent
	children           Children
	isMouseEventTarget bool
	reLayoutSuspended  bool
}

func (c *ContainerPart) Init(parent ContainerPartParent) {
	c.parent = parent
	c.children = Children{}
	parent.OnAttach(
		func() {
			for _, child := range c.children {
				child.Control.Attach()
			}
		},
	)
	parent.OnDetach(
		func() {
			for _, child := range c.children {
				child.Control.Detach()
			}
		},
	)
}

func (c *ContainerPart) SetMouseEventTarget(mouseEventTarget bool) {
	c.isMouseEventTarget = mouseEventTarget
}

func (c *ContainerPart) IsMouseEventTarget() bool {
	return c.isMouseEventTarget
}

// RelayoutSuspended returns true if adding or removing a child Control to this
// ContainerPart will not trigger a relayout of this ContainerPart. The default is false
// where any mutation will trigger a relayout.
func (c *ContainerPart) RelayoutSuspended() bool {
	return c.reLayoutSuspended
}

// SetRelayoutSuspended enables or disables relayout of the ContainerPart on
// adding or removing a child Control to this ContainerPart.
func (c *ContainerPart) SetRelayoutSuspended(enable bool) {
	c.reLayoutSuspended = true
}

// gxui.Parent compliance
func (c *ContainerPart) Children() Children {
	return c.children
}

// gxui.Container compliance
func (c *ContainerPart) AddChild(control Control) *Child {
	return c.parent.AddChildAt(len(c.children), control)
}

func (c *ContainerPart) AddChildAt(index int, control Control) *Child {
	if control.Parent() != nil {
		panic("child already has a parent")
	}
	if index < 0 || index > len(c.children) {
		panic(fmt.Errorf("index %d is out of bounds. Acceptable range: [%d - %d]", index, 0, len(c.children)))
	}

	child := &Child{Control: control}

	c.children = append(c.children, nil)
	copy(c.children[index+1:], c.children[index:])
	c.children[index] = child

	control.SetParent(c.parent)
	if c.parent.Attached() {
		control.Attach()
	}

	if !c.reLayoutSuspended {
		c.parent.ReLayout()
	}

	return child
}

func (c *ContainerPart) RemoveChild(control Control) {
	for i := range c.children {
		if c.children[i].Control == control {
			c.parent.RemoveChildAt(i)
			return
		}
	}

	panic("child not part of container")
}

func (c *ContainerPart) RemoveChildAt(index int) {
	child := c.children[index]
	c.children = append(c.children[:index], c.children[index+1:]...)
	child.Control.SetParent(nil)
	if c.parent.Attached() {
		child.Control.Detach()
	}

	if !c.reLayoutSuspended {
		c.parent.ReLayout()
	}
}

func (c *ContainerPart) RemoveAll() {
	for i := len(c.children) - 1; i >= 0; i-- {
		c.parent.RemoveChildAt(i)
	}
}

func (c *ContainerPart) ContainsPoint(point math.Point) bool {
	if !c.parent.IsVisible() || !c.parent.Size().Rect().Contains(point) {
		return false
	}

	for _, v := range c.children {
		if v.Control.ContainsPoint(point.Sub(v.Offset)) {
			return true
		}
	}

	if c.IsMouseEventTarget() {
		return true
	}

	return false
}
