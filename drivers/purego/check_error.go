package purego

import (
	"fmt"
)

func checkError(fn *Functions) {
	if v := fn.GetError(); v != 0 {
		switch v {
		case INVALID_ENUM:
			panic("GL returned error GL_INVALID_ENUM")
		case INVALID_FRAMEBUFFER_OPERATION:
			panic("GL returned error GL_INVALID_FRAMEBUFFER_OPERATION")
		case INVALID_OPERATION:
			panic("GL returned error GL_INVALID_OPERATION")
		case INVALID_VALUE:
			panic("GL returned error GL_INVALID_VALUE")
		default:
			panic(fmt.Errorf("GL returned error 0x%.4x", v))
		}
	}
}
