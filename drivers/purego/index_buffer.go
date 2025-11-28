package purego

import (
	"fmt"
	"reflect"
)

type indexBuffer struct {
	fn       *Functions
	data     []byte
	primType primitiveType
}

func newIndexBuffer(fn *Functions, primType primitiveType, dataU16 []uint16) *indexBuffer {
	switch primType {
	case ptUbyte, ptUshort, ptUint:
		if !primType.isArrayOfType(dataU16) {
			panic(fmt.Errorf("index data is not of type %v", primType))
		}
	default:
		panic(fmt.Errorf("index type must be either UBYTE, USHORT or UINT. Got: %v", primType))
	}

	// HACK: Hardcode support for only ptUshort.
	data := make([]byte, len(dataU16)*2)
	for i, v := range dataU16 {
		data[2*i+0] = byte(v >> 0)
		data[2*i+1] = byte(v >> 8)
	}

	result := &indexBuffer{data: data, primType: primType, fn: fn}
	return result
}

func (b *indexBuffer) newContext() *indexBufferContext {
	dataVal := reflect.ValueOf(b.data)
	length := dataVal.Len() / 2 // HACK: Hardcode support for only ptUshort.

	buffer := b.fn.CreateBuffer()

	b.fn.BindBuffer(ELEMENT_ARRAY_BUFFER, buffer)
	b.fn.BufferData(ELEMENT_ARRAY_BUFFER, b.data, STATIC_DRAW)
	b.fn.BindBuffer(ELEMENT_ARRAY_BUFFER, uint32(0))
	checkError(b.fn)

	globalStats.indexBufferContextCount.inc()
	return &indexBufferContext{
		glBuffer: buffer,
		primType: b.primType,
		length:   int32(length),
	}
}

type indexBufferContext struct {
	contextResource
	glBuffer uint32
	primType primitiveType
	length   int32
}

func (c *indexBufferContext) destroy(fn *Functions) {
	globalStats.indexBufferContextCount.dec()

	fn.DeleteBuffer(c.glBuffer)
	c.glBuffer = 0
}

func (c *indexBufferContext) render(mode drawMode, fn *Functions) {
	fn.BindBuffer(ELEMENT_ARRAY_BUFFER, c.glBuffer)
	fn.DrawElements(Enum(mode), c.length, Enum(c.primType), 0)
	fn.BindBuffer(ELEMENT_ARRAY_BUFFER, uint32(0))
	checkError(fn)
}
