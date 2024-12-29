// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type ContainerPartOuter interface {
	gxui.Container
	Attached() bool                         // was outer.Attachable
	Attach()                                // was outer.Attachable
	Detach()                                // was outer.Attachable
	OnAttach(func()) gxui.EventSubscription // was outer.Attachable
	OnDetach(func()) gxui.EventSubscription // was outer.Attachable
	IsVisible() bool                        // was outer.IsVisibler
	LayoutChildren()                        // was outer.LayoutChildren
	Parent() gxui.Parent                    // was outer.Parenter
	Size() math.Size                        // was outer.Sized
	SetSize(newSize math.Size)              // was outer.Sized
}

type ContainerPart struct {
	outer              ContainerPartOuter
	children           gxui.Children
	isMouseEventTarget bool
	relayoutSuspended  bool
}

func (c *ContainerPart) Init(outer ContainerPartOuter) {
	c.outer = outer
	c.children = gxui.Children{}
	outer.OnAttach(
		func() {
			for _, v := range c.children {
				v.Control.Attach()
			}
		},
	)
	outer.OnDetach(
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
func (c *ContainerPart) Children() gxui.Children {
	return c.children
}

// gxui.Container compliance
func (c *ContainerPart) AddChild(control gxui.Control) *gxui.Child {
	return c.outer.AddChildAt(len(c.children), control)
}

func (c *ContainerPart) AddChildAt(index int, control gxui.Control) *gxui.Child {
	if control.Parent() != nil {
		panic("child already has a parent")
	}
	if index < 0 || index > len(c.children) {
		panic(fmt.Errorf("index %d is out of bounds. Acceptable range: [%d - %d]", index, 0, len(c.children)))
	}

	child := &gxui.Child{Control: control}

	c.children = append(c.children, nil)
	copy(c.children[index+1:], c.children[index:])
	c.children[index] = child

	control.SetParent(c.outer)
	if c.outer.Attached() {
		control.Attach()
	}

	if !c.relayoutSuspended {
		c.outer.Relayout()
	}

	return child
}

func (c *ContainerPart) RemoveChild(control gxui.Control) {
	for i := range c.children {
		if c.children[i].Control == control {
			c.outer.RemoveChildAt(i)
			return
		}
	}

	panic("child not part of container")
}

func (c *ContainerPart) RemoveChildAt(index int) {
	child := c.children[index]
	c.children = append(c.children[:index], c.children[index+1:]...)
	child.Control.SetParent(nil)
	if c.outer.Attached() {
		child.Control.Detach()
	}

	if !c.relayoutSuspended {
		c.outer.Relayout()
	}
}

func (c *ContainerPart) RemoveAll() {
	for i := len(c.children) - 1; i >= 0; i-- {
		c.outer.RemoveChildAt(i)
	}
}

func (c *ContainerPart) ContainsPoint(point math.Point) bool {
	if !c.outer.IsVisible() || !c.outer.Size().Rect().Contains(point) {
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
