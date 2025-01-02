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

type CanvasCreatorDriver interface {
	CreateCanvas(size math.Size) Canvas
}

type DrawPaintParent interface {
	Attached() bool      // was outer.Attachable
	Parent() Parent      // was outer.Parenter
	Size() math.Size     // was outer.Sized
	Paint(canvas Canvas) // was outer.Painter

	// Attach()                                    // was outer.Attachable
	// Detach()                                    // was outer.Attachable
	// OnAttach(callback func()) EventSubscription // was outer.Attachable
	// OnDetach(callback func()) EventSubscription // was outer.Attachable
	// SetSize(newSize math.Size)                  // was outer.Sized
}

type DrawPaintPart struct {
	parent          DrawPaintParent
	driver          CanvasCreatorDriver
	canvas          Canvas
	dirty           bool
	redrawRequested bool
}

func verifyDetach(parent DrawPaintParent) {
	if parent.Attached() {
		panic(fmt.Errorf("%T garbage collected while still attached", parent))
	}
}

func (d *DrawPaintPart) Init(parent DrawPaintParent, driver CanvasCreatorDriver) {
	d.parent = parent
	d.driver = driver

	if debugVerifyDetachOnGC {
		runtime.SetFinalizer(d.parent, verifyDetach)
	}
}

func (d *DrawPaintPart) Redraw() {
	// TODO : @Badu - on desktop, why?
	//d.driver.AssertUIGoroutine()

	if !d.redrawRequested {
		if p := d.parent.Parent(); p != nil {
			d.redrawRequested = true
			p.Redraw()
		}
	}
}

func (d *DrawPaintPart) Draw() Canvas {
	if !d.parent.Attached() {
		panic(fmt.Errorf("attempting to draw a non-attached control %T", d.parent))
	}

	size := d.parent.Size()
	if size.Area() == 0 {
		return nil // No area to draw in
	}

	if d.canvas == nil || d.canvas.Size() != size || d.redrawRequested {
		d.canvas = d.driver.CreateCanvas(size)
		d.redrawRequested = false
		d.parent.Paint(d.canvas)
		d.canvas.Complete()
	}

	return d.canvas
}
