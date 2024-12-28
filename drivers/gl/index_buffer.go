// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/goxjs/gl"
)

type indexBuffer struct {
	data     []byte
	primType primitiveType
}

func newIndexBuffer(primType primitiveType, dataU16 []uint16) *indexBuffer {
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

	result := &indexBuffer{data: data, primType: primType}
	return result
}

func (b *indexBuffer) newContext() *indexBufferContext {
	dataVal := reflect.ValueOf(b.data)
	length := dataVal.Len() / 2 // HACK: Hardcode support for only ptUshort.

	buffer := gl.CreateBuffer()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, b.data, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{})
	checkError()

	globalStats.indexBufferContextCount.inc()
	return &indexBufferContext{
		glBuffer: buffer,
		primType: b.primType,
		length:   length,
	}
}

type indexBufferContext struct {
	contextResource
	glBuffer gl.Buffer
	primType primitiveType
	length   int
}

func (c *indexBufferContext) destroy() {
	globalStats.indexBufferContextCount.dec()

	gl.DeleteBuffer(c.glBuffer)

	c.glBuffer = gl.Buffer{}
}

func (c *indexBufferContext) render(mode drawMode) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.glBuffer)
	gl.DrawElements(gl.Enum(mode), c.length, gl.Enum(c.primType), 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{})
	checkError()
}
