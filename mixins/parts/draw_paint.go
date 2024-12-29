// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"fmt"
	"github.com/badu/gxui/math"
	"runtime"

	"github.com/badu/gxui"
)

const debugVerifyDetachOnGC = false

type DrawPaintOuter interface {
	Attached() bool                                  // was outer.Attachable
	Attach()                                         // was outer.Attachable
	Detach()                                         // was outer.Attachable
	OnAttach(callback func()) gxui.EventSubscription // was outer.Attachable
	OnDetach(callback func()) gxui.EventSubscription // was outer.Attachable
	Paint(canvas gxui.Canvas)                        // was outer.Painter
	Parent() gxui.Parent                             // was outer.Parenter
	Size() math.Size                                 // was outer.Sized
	SetSize(newSize math.Size)                       // was outer.Sized
}

type DrawPaint struct {
	outer           DrawPaintOuter
	driver          gxui.Driver
	canvas          gxui.Canvas
	dirty           bool
	redrawRequested bool
}

func verifyDetach(outer DrawPaintOuter) {
	if outer.Attached() {
		panic(fmt.Errorf("%T garbage collected while still attached", outer))
	}
}

func (d *DrawPaint) Init(outer DrawPaintOuter, theme gxui.Theme) {
	d.outer = outer
	d.driver = theme.Driver()

	if debugVerifyDetachOnGC {
		runtime.SetFinalizer(d.outer, verifyDetach)
	}
}

func (d *DrawPaint) Redraw() {
	d.driver.AssertUIGoroutine()

	if !d.redrawRequested {
		if p := d.outer.Parent(); p != nil {
			d.redrawRequested = true
			p.Redraw()
		}
	}
}

func (d *DrawPaint) Draw() gxui.Canvas {
	if !d.outer.Attached() {
		panic(fmt.Errorf("attempting to draw a non-attached control %T", d.outer))
	}

	size := d.outer.Size()
	if size.Area() == 0 {
		return nil // No area to draw in
	}

	if d.canvas == nil || d.canvas.Size() != size || d.redrawRequested {
		d.canvas = d.driver.CreateCanvas(size)
		d.redrawRequested = false
		d.outer.Paint(d.canvas)
		d.canvas.Complete()
	}

	return d.canvas
}
