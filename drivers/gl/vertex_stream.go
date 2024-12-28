// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"math"
	"reflect"

	"github.com/goxjs/gl"
)

type vertexStream struct {
	name       string
	data       []byte
	shaderType shaderDataType
	count      int
}

func newVertexStream(name string, shaderType shaderDataType, data32 []float32) *vertexStream {
	dataVal := reflect.ValueOf(data32)
	dataLen := dataVal.Len()

	if dataLen%shaderType.vectorElementCount() != 0 {
		panic(fmt.Errorf("incorrect multiple of elements. Got: %d, Requires multiple of %d", dataLen, shaderType.vectorElementCount()))
	}
	if !shaderType.vectorElementType().isArrayOfType(data32) {
		panic("Data is not of the specified type")
	}

	// HACK.
	data := float32Bytes(data32...)

	stream := &vertexStream{
		name:       name,
		data:       data,
		shaderType: shaderType,
		count:      dataLen / shaderType.vectorElementCount(),
	}
	return stream
}

func (s *vertexStream) newContext() *vertexStreamContext {
	buffer := gl.CreateBuffer()

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, s.data, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	checkError()

	globalStats.vertexStreamContextCount.inc()
	return &vertexStreamContext{glBuffer: buffer}
}

type vertexStreamContext struct {
	contextResource
	glBuffer gl.Buffer
}

func (c *vertexStreamContext) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.glBuffer)
}

func (c *vertexStreamContext) destroy() {
	globalStats.vertexStreamContextCount.dec()
	gl.DeleteBuffer(c.glBuffer)
	c.glBuffer = gl.Buffer{}
}

// float32Bytes returns the byte representation of float32 values in little endian byte order.
func float32Bytes(values ...float32) []byte {
	result := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		result[4*i+0] = byte(u >> 0)
		result[4*i+1] = byte(u >> 8)
		result[4*i+2] = byte(u >> 16)
		result[4*i+3] = byte(u >> 24)
	}
	return result
}
