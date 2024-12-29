// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/badu/gxui/math"

	"github.com/goxjs/gl"
)

// contextResource is used as an anonymous field by types that are constructed per context.
type contextResource struct {
	lastContextUse int // used for mark-and-sweeping the resource.
}

type context struct {
	blitter              *blitter
	resolution           resolution
	stats                contextStats
	textureContexts      map[*TextureImpl]*textureContext
	vertexStreamContexts map[*vertexStream]*vertexStreamContext
	indexBufferContexts  map[*indexBuffer]*indexBufferContext
	sizeDips             math.Size
	sizePixels           math.Size
	clip                 math.Rect
	frame                int
}

func newContext() *context {
	result := &context{
		textureContexts:      make(map[*TextureImpl]*textureContext),
		vertexStreamContexts: make(map[*vertexStream]*vertexStreamContext),
		indexBufferContexts:  make(map[*indexBuffer]*indexBufferContext),
	}
	result.blitter = newBlitter(result, &result.stats)
	return result
}

func (c *context) destroy() {
	for textureCtx, tc := range c.textureContexts {
		delete(c.textureContexts, textureCtx)
		tc.destroy()
		c.stats.textureCount--
	}

	for stream, sc := range c.vertexStreamContexts {
		delete(c.vertexStreamContexts, stream)
		sc.destroy()
		c.stats.vertexStreamCount--
	}

	for buffer, ic := range c.indexBufferContexts {
		delete(c.indexBufferContexts, buffer)
		ic.destroy()
		c.stats.indexBufferCount--
	}

	c.blitter.destroy(c)
	c.blitter = nil
}

func (c *context) beginDraw(sizeDips, sizePixels math.Size) {
	dipsToPixels := float32(sizePixels.W) / float32(sizeDips.W)

	c.sizeDips = sizeDips
	c.sizePixels = sizePixels
	c.resolution = resolution(dipsToPixels*65536 + 0.5)

	c.stats.drawCallCount = 0
	c.stats.timer("Frame").start()
}

func (c *context) endDraw() {
	// Reap any unused resources
	for textureCtx, tc := range c.textureContexts {
		if tc.lastContextUse != c.frame {
			delete(c.textureContexts, textureCtx)
			tc.destroy()
			c.stats.textureCount--
		}
	}

	for stream, sc := range c.vertexStreamContexts {
		if sc.lastContextUse != c.frame {
			delete(c.vertexStreamContexts, stream)
			sc.destroy()
			c.stats.vertexStreamCount--
		}
	}

	for buffer, ic := range c.indexBufferContexts {
		if ic.lastContextUse != c.frame {
			delete(c.indexBufferContexts, buffer)
			ic.destroy()
			c.stats.indexBufferCount--
		}
	}

	c.stats.timer("Frame").stop()
	c.stats.frameCount++
	c.frame++
}

func (c *context) getOrCreateTextureContext(targetTexture *TextureImpl) *textureContext {
	textureCtx, found := c.textureContexts[targetTexture]
	if !found {
		textureCtx = targetTexture.newContext()
		c.textureContexts[targetTexture] = textureCtx
		c.stats.textureCount++
	}
	textureCtx.lastContextUse = c.frame
	return textureCtx
}

func (c *context) getOrCreateVertexStreamContext(targetStream *vertexStream) *vertexStreamContext {
	stream, found := c.vertexStreamContexts[targetStream]
	if !found {
		stream = targetStream.newContext()
		c.vertexStreamContexts[targetStream] = stream
		c.stats.vertexStreamCount++
	}
	stream.lastContextUse = c.frame
	return stream
}

func (c *context) getOrCreateIndexBufferContext(targetBuffer *indexBuffer) *indexBufferContext {
	buffer, found := c.indexBufferContexts[targetBuffer]
	if !found {
		buffer = targetBuffer.newContext()
		c.indexBufferContexts[targetBuffer] = buffer
		c.stats.indexBufferCount++
	}
	buffer.lastContextUse = c.frame
	return buffer
}

func (c *context) apply(state *drawState) {
	rect := state.ClipPixels
	cClip := c.clip
	if cClip != rect {
		c.clip = rect
		sizePixels := c.sizePixels
		rectSize := rect.Size()
		gl.Scissor(int32(rect.Min.X), int32(sizePixels.H)-int32(rect.Max.Y), int32(rectSize.W), int32(rectSize.H))
	}
}
