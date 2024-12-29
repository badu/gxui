// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/parts"
)

type ControlOuter interface {
	gxui.Control
	Paint(canvas gxui.Canvas) // was outer.Painter
	Redraw()                  // was outer.Redrawer
	Relayout()                // was outer.Relayouter
}

type Control struct {
	parts.Attachable
	parts.DrawPaint
	parts.InputEventHandler
	parts.Layoutable
	parts.Parentable
	parts.Visible
}

func (c *Control) Init(outer ControlOuter, theme gxui.Theme) {
	c.Attachable.Init(outer)
	c.DrawPaint.Init(outer, theme)
	c.Layoutable.Init(outer, theme)
	c.InputEventHandler.Init(outer)
	c.Parentable.Init(outer)
	c.Visible.Init(outer)
}

func (c *Control) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (c *Control) ContainsPoint(point math.Point) bool {
	return c.IsVisible() && c.Size().Rect().Contains(point)
}
