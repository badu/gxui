package purego

import (
	"fmt"
)

type drawMode int

const (
	dmPoints        drawMode = POINTS
	dmLineStrip     drawMode = LINE_STRIP
	dmLineLoop      drawMode = LINE_LOOP
	dmLines         drawMode = LINES
	dmTriangleStrip drawMode = TRIANGLE_STRIP
	dmTriangleFan   drawMode = TRIANGLE_FAN
	dmTriangles     drawMode = TRIANGLES
)

func (d drawMode) primitiveCount(vertexCount int) int {
	switch d {
	case dmPoints:
		return vertexCount
	case dmLineStrip:
		return vertexCount - 1
	case dmLineLoop:
		return vertexCount
	case dmLines:
		return vertexCount / 2
	case dmTriangleStrip:
		return vertexCount - 2
	case dmTriangleFan:
		return vertexCount - 2
	case dmTriangles:
		return vertexCount / 3
	default:
		panic(fmt.Errorf("unknown drawMode 0x%.4x", d))
	}
}
