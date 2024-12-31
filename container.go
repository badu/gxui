// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"strings"

	"github.com/badu/gxui/math"
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
	Parent
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
}

type ContainerBaseNoControlOuter interface {
	Container
	PaintChild(canvas Canvas, child *Child, idx int) // was outer.PaintChilder
	Paint(canvas Canvas)                             // was outer.Painter
	LayoutChildren()                                 // was outer.LayoutChildren
}

type BaseContainerParent interface {
	Container
	Control
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
	Container
	Attached() bool                             // was outer.Attachable
	Attach()                                    // was outer.Attachable
	Detach()                                    // was outer.Attachable
	OnAttach(callback func()) EventSubscription // was outer.Attachable
	OnDetach(callback func()) EventSubscription // was outer.Attachable
	IsVisible() bool                            // was outer.IsVisibler
	LayoutChildren()                            // was outer.LayoutChildren
	Parent() Parent                             // was outer.Parenter
	Size() math.Size                            // was outer.Sized
	SetSize(newSize math.Size)                  // was outer.Sized
}

type ContainerPart struct {
	parent             ContainerPartParent
	children           Children
	isMouseEventTarget bool
	relayoutSuspended  bool
}

func (c *ContainerPart) Init(parent ContainerPartParent) {
	c.parent = parent
	c.children = Children{}
	parent.OnAttach(
		func() {
			for _, v := range c.children {
				v.Control.Attach()
			}
		},
	)
	parent.OnDetach(
		func() {
			for _, v := range c.children {
				v.Control.Detach()
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
	return c.relayoutSuspended
}

// SetRelayoutSuspended enables or disables relayout of the ContainerPart on
// adding or removing a child Control to this ContainerPart.
func (c *ContainerPart) SetRelayoutSuspended(enable bool) {
	c.relayoutSuspended = true
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

	if !c.relayoutSuspended {
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

	if !c.relayoutSuspended {
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

type ContainerBase struct {
	AttachablePart
	ContainerPart
	DrawPaintPart
	InputEventHandlerPart
	LayoutablePart
	PaddablePart
	PaintChildrenPart
	ParentablePart
	VisiblePart
}

func (c *ContainerBase) Init(parent BaseContainerParent, driver Driver) {
	c.ContainerPart.Init(parent)
	c.DrawPaintPart.Init(parent, driver)
	c.InputEventHandlerPart.Init()
	c.LayoutablePart.Init(parent, driver)
	c.PaddablePart.Init(parent)
	c.PaintChildrenPart.Init(parent)
	c.ParentablePart.Init()
	c.VisiblePart.Init(parent)
}
