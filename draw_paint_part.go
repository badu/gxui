// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/badu/gxui/math"
	"runtime"
)

const debugVerifyDetachOnGC = false

type DrawPaintOuter interface {
	Attached() bool                             // was outer.Attachable
	Attach()                                    // was outer.Attachable
	Detach()                                    // was outer.Attachable
	OnAttach(callback func()) EventSubscription // was outer.Attachable
	OnDetach(callback func()) EventSubscription // was outer.Attachable
	Paint(canvas Canvas)                        // was outer.Painter
	Parent() Parent                             // was outer.Parenter
	Size() math.Size                            // was outer.Sized
	SetSize(newSize math.Size)                  // was outer.Sized
}

type DrawPaintPart struct {
	outer           DrawPaintOuter
	driver          Driver
	canvas          Canvas
	dirty           bool
	redrawRequested bool
}

func verifyDetach(outer DrawPaintOuter) {
	if outer.Attached() {
		panic(fmt.Errorf("%T garbage collected while still attached", outer))
	}
}

func (d *DrawPaintPart) Init(outer DrawPaintOuter, theme App) {
	d.outer = outer
	d.driver = theme.Driver()

	if debugVerifyDetachOnGC {
		runtime.SetFinalizer(d.outer, verifyDetach)
	}
}

func (d *DrawPaintPart) Redraw() {
	d.driver.AssertUIGoroutine()

	if !d.redrawRequested {
		if p := d.outer.Parent(); p != nil {
			d.redrawRequested = true
			p.Redraw()
		}
	}
}

func (d *DrawPaintPart) Draw() Canvas {
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
