// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type ScalingMode int

const (
	Scaling1to1 ScalingMode = iota
	ScalingExpandGreedy
	ScalingExplicitSize
)

type AspectMode int

const (
	AspectStretch = iota
	AspectCorrectLetterbox
	AspectCorrectCrop
)

type Image interface {
	Control
	Texture() Texture
	SetTexture(Texture)
	Canvas() Canvas
	SetCanvas(Canvas)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	ScalingMode() ScalingMode
	SetScalingMode(ScalingMode)
	SetExplicitSize(size math.Size)
	AspectMode() AspectMode
	SetAspectMode(AspectMode)
	PixelAt(point math.Point) (math.Point, bool) // TODO: Remove
}

type ImageImpl struct {
	ControlBase
	BackgroundBorderPainter
	parent       ControlBaseParent
	texture      Texture
	canvas       Canvas
	scalingMode  ScalingMode
	aspectMode   AspectMode
	explicitSize math.Size
}

func (i *ImageImpl) calculateDrawRect() math.Rect {
	rect := i.parent.Size().Rect()
	texW, texH := i.texture.Size().WH()
	aspectSrc := float32(texH) / float32(texW)
	aspectDst := float32(rect.H()) / float32(rect.W())
	switch i.aspectMode {
	case AspectCorrectLetterbox, AspectCorrectCrop:
		if (aspectDst < aspectSrc) != (i.aspectMode == AspectCorrectLetterbox) {
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

func (i *ImageImpl) Init(parent ControlBaseParent, driver Driver) {
	i.parent = parent
	i.ControlBase.Init(parent, driver)
	i.BackgroundBorderPainter.Init(parent)
	i.SetBorderPen(TransparentPen)
	i.SetBackgroundBrush(TransparentBrush)
}

func (i *ImageImpl) Texture() Texture {
	return i.texture
}

func (i *ImageImpl) SetTexture(texture Texture) {
	if i.texture != texture {
		i.texture = texture
		i.canvas = nil
		i.parent.Relayout()
	}
}

func (i *ImageImpl) Canvas() Canvas {
	return i.canvas
}

func (i *ImageImpl) SetCanvas(canvas Canvas) {
	if !canvas.IsComplete() {
		panic("SetCanvas() called with an incomplete canvas")
	}

	if i.canvas != canvas {
		i.canvas = canvas
		i.texture = nil
		i.parent.Relayout()
	}
}

func (i *ImageImpl) ScalingMode() ScalingMode {
	return i.scalingMode
}

func (i *ImageImpl) SetScalingMode(mode ScalingMode) {
	if i.scalingMode != mode {
		i.scalingMode = mode
		i.parent.Relayout()
	}
}

func (i *ImageImpl) AspectMode() AspectMode {
	return i.aspectMode
}

func (i *ImageImpl) SetAspectMode(mode AspectMode) {
	if i.aspectMode != mode {
		i.aspectMode = mode
		i.parent.Redraw()
	}
}

func (i *ImageImpl) SetExplicitSize(explicitSize math.Size) {
	if i.explicitSize != explicitSize {
		i.explicitSize = explicitSize
		i.parent.Relayout()
	}
	i.SetScalingMode(ScalingExplicitSize)
}

func (i *ImageImpl) PixelAt(point math.Point) (math.Point, bool) {
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

func (i *ImageImpl) DesiredSize(min, max math.Size) math.Size {
	size := max
	switch i.scalingMode {
	case ScalingExplicitSize:
		size = i.explicitSize
	case Scaling1to1:
		switch {
		case i.texture != nil:
			size = i.texture.Size()
		case i.canvas != nil:
			size = i.canvas.Size()
		}
	}
	return size.Expand(math.CreateSpacing(int(i.BorderPen().Width))).Clamp(min, max)
}

func (i *ImageImpl) Paint(canvas Canvas) {
	rect := i.parent.Size().Rect()
	i.PaintBackground(canvas, rect)
	switch {
	case i.texture != nil:
		canvas.DrawTexture(i.texture, i.calculateDrawRect())
	case i.canvas != nil:
		canvas.DrawCanvas(i.canvas, math.ZeroPoint)
	}
	i.PaintBorder(canvas, rect)
}
