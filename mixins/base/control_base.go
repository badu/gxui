// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type ControlBaseOuter interface {
	gxui.Control
	Paint(canvas gxui.Canvas) // was outer.Painter
	Redraw()                  // was outer.Redrawer
	Relayout()                // was outer.Relayouter
}

type ControlBase struct {
	AttachablePart
	DrawPaintPart
	InputEventHandlerPart
	LayoutablePart
	ParentablePart
	VisiblePart
}

func (c *ControlBase) Init(outer ControlBaseOuter, theme gxui.Theme) {
	c.AttachablePart.Init(outer)
	c.DrawPaintPart.Init(outer, theme)
	c.LayoutablePart.Init(outer, theme)
	c.InputEventHandlerPart.Init(outer)
	c.ParentablePart.Init(outer)
	c.VisiblePart.Init(outer)
}

func (c *ControlBase) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (c *ControlBase) ContainsPoint(point math.Point) bool {
	return c.IsVisible() && c.Size().Rect().Contains(point)
}
