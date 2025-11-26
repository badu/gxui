package purego

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

func newQuadShape(fn *Functions) *shape {
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
		fn,
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
		ctx.getOrCreateIndexBufferContext(s.indexBuffer).render(s.drawMode, ctx.fn)
	} else {
		ctx.fn.DrawArrays(Enum(s.drawMode), 0, s.vertexBuffer.count)
	}

	shader.unbind(ctx.fn)
	checkError(ctx.fn)
}
