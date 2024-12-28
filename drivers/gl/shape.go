// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/goxjs/gl"

type shape struct {
	vertexBuffer *vertexBuffer
	indexBuffer  *indexBuffer
	drawMode     drawMode
}

func newShape(vertexBuffer *vertexBuffer, indexBuffer *indexBuffer, mode drawMode) *shape {
	if vertexBuffer == nil {
		panic("VertexBuffer cannot be nil")
	}

	result := &shape{
		vertexBuffer: vertexBuffer,
		indexBuffer:  indexBuffer,
		drawMode:     mode,
	}
	return result
}

func newQuadShape() *shape {
	pos := newVertexStream(
		"aPosition",
		stFloatVec2,
		[]float32{
			0.0, 0.0,
			1.0, 0.0,
			0.0, 1.0,
			1.0, 1.0,
		},
	)
	vBuffer := newVertexBuffer(pos)
	iBuffer := newIndexBuffer(
		ptUshort,
		[]uint16{
			0, 1, 2,
			2, 1, 3,
		},
	)

	return newShape(vBuffer, iBuffer, dmTriangles)
}

func (s shape) draw(ctx *context, shader *shaderProgram, bindings uniformBindings) {
	shader.bind(ctx, s.vertexBuffer, bindings)

	if s.indexBuffer != nil {
		ctx.getOrCreateIndexBufferContext(s.indexBuffer).render(s.drawMode)
	} else {
		gl.DrawArrays(gl.Enum(s.drawMode), 0, s.vertexBuffer.count)
	}

	shader.unbind(ctx)
	checkError()
}
