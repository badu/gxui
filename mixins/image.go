// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/base"
)

type Image struct {
	base.ControlBase
	base.BackgroundBorderPainter
	outer        base.ControlBaseOuter
	texture      gxui.Texture
	canvas       gxui.Canvas
	scalingMode  gxui.ScalingMode
	aspectMode   gxui.AspectMode
	explicitSize math.Size
}

func (i *Image) calculateDrawRect() math.Rect {
	rect := i.outer.Size().Rect()
	texW, texH := i.texture.Size().WH()
	aspectSrc := float32(texH) / float32(texW)
	aspectDst := float32(rect.H()) / float32(rect.W())
	switch i.aspectMode {
	case gxui.AspectCorrectLetterbox, gxui.AspectCorrectCrop:
		if (aspectDst < aspectSrc) != (i.aspectMode == gxui.AspectCorrectLetterbox) {
			contract := rect.H() - int(float32(rect.W())*aspectSrc)
			rect = rect.Contract(math.Spacing{T: contract / 2, B: contract / 2})
		} else {
			contract := rect.W() - int(float32(rect.H())/aspectSrc)
			rect = rect.Contract(math.Spacing{L: contract / 2, R: contract / 2})
		}
	default:
		//
	}
	return rect
}

func (i *Image) Init(outer base.ControlBaseOuter, theme gxui.Theme) {
	i.outer = outer
	i.ControlBase.Init(outer, theme)
	i.BackgroundBorderPainter.Init(outer)
	i.SetBorderPen(gxui.TransparentPen)
	i.SetBackgroundBrush(gxui.TransparentBrush)
}

func (i *Image) Texture() gxui.Texture {
	return i.texture
}

func (i *Image) SetTexture(texture gxui.Texture) {
	if i.texture != texture {
		i.texture = texture
		i.canvas = nil
		i.outer.Relayout()
	}
}

func (i *Image) Canvas() gxui.Canvas {
	return i.canvas
}

func (i *Image) SetCanvas(canvas gxui.Canvas) {
	if !canvas.IsComplete() {
		panic("SetCanvas() called with an incomplete canvas")
	}

	if i.canvas != canvas {
		i.canvas = canvas
		i.texture = nil
		i.outer.Relayout()
	}
}

func (i *Image) ScalingMode() gxui.ScalingMode {
	return i.scalingMode
}

func (i *Image) SetScalingMode(mode gxui.ScalingMode) {
	if i.scalingMode != mode {
		i.scalingMode = mode
		i.outer.Relayout()
	}
}

func (i *Image) AspectMode() gxui.AspectMode {
	return i.aspectMode
}

func (i *Image) SetAspectMode(mode gxui.AspectMode) {
	if i.aspectMode != mode {
		i.aspectMode = mode
		i.outer.Redraw()
	}
}

func (i *Image) SetExplicitSize(explicitSize math.Size) {
	if i.explicitSize != explicitSize {
		i.explicitSize = explicitSize
		i.outer.Relayout()
	}
	i.SetScalingMode(gxui.ScalingExplicitSize)
}

func (i *Image) PixelAt(point math.Point) (math.Point, bool) {
	rect := i.calculateDrawRect()
	if tex := i.Texture(); tex != nil {
		size := tex.SizePixels()
		point = point.Sub(rect.Min).
			ScaleX(float32(size.W) / float32(rect.W())).
			ScaleY(float32(size.H) / float32(rect.H()))
		if size.Rect().Contains(point) {
			return point, true
		}
	}
	return math.Point{X: -1, Y: -1}, false
}

func (i *Image) DesiredSize(min, max math.Size) math.Size {
	size := max
	switch i.scalingMode {
	case gxui.ScalingExplicitSize:
		size = i.explicitSize
	case gxui.Scaling1to1:
		switch {
		case i.texture != nil:
			size = i.texture.Size()
		case i.canvas != nil:
			size = i.canvas.Size()
		}
	}
	return size.Expand(math.CreateSpacing(int(i.BorderPen().Width))).Clamp(min, max)
}

func (i *Image) Paint(canvas gxui.Canvas) {
	rect := i.outer.Size().Rect()
	i.PaintBackground(canvas, rect)
	switch {
	case i.texture != nil:
		canvas.DrawTexture(i.texture, i.calculateDrawRect())
	case i.canvas != nil:
		canvas.DrawCanvas(i.canvas, math.ZeroPoint)
	}
	i.PaintBorder(canvas, rect)
}
