// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/goxjs/gl"
)

type drawStateStack []drawState

func (s *drawStateStack) head() *drawState {
	return &(*s)[len(*s)-1]
}
func (s *drawStateStack) push(ds drawState) {
	*s = append(*s, ds)
}
func (s *drawStateStack) pop() {
	*s = (*s)[:len(*s)-1]
}

type canvasOp func(ctx *context, stack *drawStateStack)

type drawState struct {
	// The below are all in window coordinates
	ClipPixels   math.Rect
	OriginPixels math.Point
}

type CanvasImpl struct {
	sizeDips          math.Size
	ops               []canvasOp
	built             bool
	buildingPushCount int
}

func NewCanvas(sizeDips math.Size) *CanvasImpl {
	if sizeDips.W <= 0 || sizeDips.H < 0 {
		panic(fmt.Errorf("canvas width and height must be positive. Size: %d", sizeDips))
	}

	result := &CanvasImpl{sizeDips: sizeDips}
	return result
}

func (c *CanvasImpl) draw(ctx *context, stack *drawStateStack) {
	head := stack.head()
	ctx.apply(head)

	for _, op := range c.ops {
		op(ctx, stack)
	}
}

func (c *CanvasImpl) appendOp(name string, op canvasOp) {
	if c.built {
		panic(fmt.Errorf("%s() called after Complete()", name))
	}
	c.ops = append(c.ops, op)
}

// gxui.Canvas compliance
func (c *CanvasImpl) Size() math.Size {
	return c.sizeDips
}

func (c *CanvasImpl) IsComplete() bool {
	return c.built
}

func (c *CanvasImpl) Complete() {
	if c.built {
		panic("complete() called twice")
	}

	if c.buildingPushCount != 0 {
		panic(fmt.Errorf("push() count was %d when calling Complete", c.buildingPushCount))
	}

	c.built = true
}

func (c *CanvasImpl) Push() {
	c.buildingPushCount++
	c.appendOp(
		"Push",
		func(ctx *context, stack *drawStateStack) {
			stack.push(*stack.head())
		},
	)
}

func (c *CanvasImpl) Pop() {
	c.buildingPushCount--
	c.appendOp(
		"Pop",
		func(ctx *context, stack *drawStateStack) {
			stack.pop()
			ctx.apply(stack.head())
		},
	)
}

func (c *CanvasImpl) AddClip(rect math.Rect) {
	c.appendOp(
		"AddClip",
		func(ctx *context, stack *drawStateStack) {
			head := stack.head()
			rectLocalPixels := ctx.resolution.rectDipsToPixels(rect)
			rectWindowPixels := rectLocalPixels.Offset(head.OriginPixels)
			head.ClipPixels = head.ClipPixels.Intersect(rectWindowPixels)
			ctx.apply(head)
		},
	)
}

func (c *CanvasImpl) Clear(color gxui.Color) {
	c.appendOp(
		"Clear",
		func(ctx *context, stack *drawStateStack) {
			gl.ClearColor(color.R, color.G, color.B, color.A)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		},
	)
}

func (c *CanvasImpl) DrawCanvas(targetCanvas gxui.Canvas, offsetDips math.Point) {
	if targetCanvas == nil {
		panic("target canvas cannot be nil")
	}

	childCanvas := targetCanvas.(*CanvasImpl)
	c.appendOp(
		"DrawCanvas",
		func(ctx *context, stack *drawStateStack) {
			offsetPixels := ctx.resolution.pointDipsToPixels(offsetDips)
			stack.push(*stack.head())
			head := stack.head()
			head.OriginPixels = head.OriginPixels.Add(offsetPixels)
			childCanvas.draw(ctx, stack)
			stack.pop()
			ctx.apply(stack.head())
		},
	)
}

func (c *CanvasImpl) DrawRunes(useFont gxui.Font, runes []rune, points []math.Point, color gxui.Color) {
	if useFont == nil {
		panic("font cannot be nil")
	}

	runesCopy := append([]rune{}, runes...)
	pointsCopy := append([]math.Point{}, points...)
	c.appendOp(
		"DrawRunes",
		func(ctx *context, stack *drawStateStack) {
			useFont.(*font).DrawRunes(ctx, runesCopy, pointsCopy, color, stack.head())
		},
	)
}

func (c *CanvasImpl) DrawLines(lines gxui.Polygon, pen gxui.Pen) {
	edge := openPolyToShape(lines, pen.Width)
	c.appendOp(
		"DrawLines",
		func(ctx *context, dss *drawStateStack) {
			head := dss.head()
			if edge != nil && pen.Color.A > 0 {
				ctx.blitter.blitShape(ctx, *edge, pen.Color, head)
			}
		},
	)
}

func (c *CanvasImpl) DrawPolygon(poly gxui.Polygon, pen gxui.Pen, brush gxui.Brush) {
	fill, edge := closedPolyToShape(poly, pen.Width)
	c.appendOp(
		"DrawPolygon",
		func(ctx *context, stack *drawStateStack) {
			head := stack.head()
			if fill != nil && brush.Color.A > 0 {
				ctx.blitter.blitShape(ctx, *fill, brush.Color, head)
			}
			if edge != nil && pen.Color.A > 0 {
				ctx.blitter.blitShape(ctx, *edge, pen.Color, head)
			}
		},
	)
}

func (c *CanvasImpl) DrawRect(rect math.Rect, brush gxui.Brush) {
	c.appendOp(
		"DrawRect",
		func(ctx *context, dss *drawStateStack) {
			ctx.blitter.blitRect(ctx, ctx.resolution.rectDipsToPixels(rect), brush.Color, dss.head())
		},
	)
}

func (c *CanvasImpl) DrawRoundedRect(rect math.Rect, tl, tr, bl, br float32, pen gxui.Pen, brush gxui.Brush) {
	if tl == 0 && tr == 0 && bl == 0 && br == 0 && pen.Color.A == 0 {
		c.DrawRect(rect, brush)
		return
	}

	polygon := gxui.Polygon{
		gxui.PolygonVertex{Position: rect.TL(), RoundedRadius: tl},
		gxui.PolygonVertex{Position: rect.TR(), RoundedRadius: tr},
		gxui.PolygonVertex{Position: rect.BR(), RoundedRadius: br},
		gxui.PolygonVertex{Position: rect.BL(), RoundedRadius: bl},
	}

	c.DrawPolygon(polygon, pen, brush)
}

func (c *CanvasImpl) DrawTexture(targetTexture gxui.Texture, r math.Rect) {
	if targetTexture == nil {
		panic("target texture cannot be nil")
	}

	c.appendOp(
		"DrawTexture",
		func(ctx *context, stack *drawStateStack) {
			textureCtx := ctx.getOrCreateTextureContext(targetTexture.(*TextureImpl))
			ctx.blitter.blit(ctx, textureCtx, textureCtx.sizePixels.Rect(), ctx.resolution.rectDipsToPixels(r), stack.head())
		},
	)
}
