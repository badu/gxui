// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"

	"github.com/badu/gxui/pkg/math"
)

type Parent interface {
	Children() Children
	ReLayout()
	Redraw()
}

type PainterAndLayouter interface {
	PaintChild(canvas Canvas, child *Child, idx int)
	Paint(canvas Canvas)
	LayoutChildren()
}

type ContainerBaseNoControl interface {
	Container
	PainterAndLayouter
}

type BaseContainerParent interface {
	Container
	Control
	PainterAndLayouter
}

type ContainerPartParent interface {
	Parent
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

type ContainerBase struct {
	InputEventHandlerPart
	PaintChildrenPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	ContainerPart
	PaddablePart
	LayoutablePart
}

func (c *ContainerBase) Init(parent BaseContainerParent, driver Driver) {
	c.ContainerPart.Init(parent)
	c.DrawPaintPart.Init(parent, driver)
	c.InputEventHandlerPart.Init()
	c.LayoutablePart.Init(parent)
	c.PaddablePart.Init(parent)
	c.PaintChildrenPart.Init(parent)
	c.ParentablePart.Init()
	c.VisiblePart.Init(parent)
}
