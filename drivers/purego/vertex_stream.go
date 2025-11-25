package purego

import (
	"fmt"
	"math"
	"reflect"
)

type vertexStream struct {
	fn         *Functions // TODO : fill me up
	name       string
	data       []byte
	shaderType shaderDataType
	count      int
}

func newVertexStream(fn *Functions, name string, shaderType shaderDataType, data32 []float32) *vertexStream {
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
		fn:         fn,
		name:       name,
		data:       data,
		shaderType: shaderType,
		count:      dataLen / shaderType.vectorElementCount(),
	}
	return stream
}

func (s *vertexStream) newContext() *vertexStreamContext {
	buffer := s.fn.CreateBuffer()

	s.fn.BindBuffer(ARRAY_BUFFER, buffer)
	s.fn.BufferData(ARRAY_BUFFER, s.data, STATIC_DRAW)
	s.fn.BindBuffer(ARRAY_BUFFER, Buffer{})
	checkError(s.fn)

	globalStats.vertexStreamContextCount.inc()
	return &vertexStreamContext{glBuffer: buffer}
}

type vertexStreamContext struct {
	contextResource
	fn       *Functions // TODO : fill me up
	glBuffer Buffer
}

func (c *vertexStreamContext) bind() {
	c.fn.BindBuffer(ARRAY_BUFFER, c.glBuffer)
}

func (c *vertexStreamContext) destroy() {
	globalStats.vertexStreamContextCount.dec()
	c.fn.DeleteBuffer(c.glBuffer)
	c.glBuffer = Buffer{}
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
