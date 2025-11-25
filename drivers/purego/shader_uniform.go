package purego

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/pkg/math"
)

type shaderUniform struct {
	fn          *Functions
	name        string
	size        int
	shaderType  shaderDataType
	location    Uniform
	textureUnit int
}

func (u *shaderUniform) bind(target interface{}) {
	switch u.shaderType {
	case stFloatMat2:
		u.fn.UniformMatrix2fv(u.location, target.([]float32))

	case stFloatMat3:
		switch matrix := target.(type) {
		case math.Mat3:
			u.fn.UniformMatrix3fv(u.location, matrix[:])

		case []float32:
			u.fn.UniformMatrix3fv(u.location, matrix)
		}

	case stFloatMat4:
		u.fn.UniformMatrix4fv(u.location, target.([]float32))

	case stFloatVec1:
		switch vector := target.(type) {
		case float32:
			u.fn.Uniform1f(u.location, vector)

		case []float32:
			u.fn.Uniform1fv(u.location, vector)
		}

	case stFloatVec2:
		switch vector := target.(type) {
		case math.Vec2:
			u.fn.Uniform2fv(u.location, []float32{vector.X, vector.Y})

		case []float32:
			if len(vector)%2 != 0 {
				panic(fmt.Errorf("uniform '%s' of type vec2 should be an float32 array with a multiple of two length", u.name))
			}
			u.fn.Uniform2fv(u.location, vector)
		}

	case stFloatVec3:
		switch vector := target.(type) {
		case math.Vec3:
			u.fn.Uniform3fv(u.location, []float32{vector.X, vector.Y, vector.Z})

		case []float32:
			if len(vector)%3 != 0 {
				panic(fmt.Errorf("uniform '%s' of type vec3 should be an float32 array with a multiple of three length", u.name))
			}
			u.fn.Uniform3fv(u.location, vector)
		}

	case stFloatVec4:
		switch vector := target.(type) {
		case math.Vec4:
			u.fn.Uniform4fv(u.location, []float32{vector.X, vector.Y, vector.Z, vector.W})

		case gxui.Color:
			u.fn.Uniform4fv(u.location, []float32{vector.R, vector.G, vector.B, vector.A})

		case []float32:
			if len(vector)%4 != 0 {
				panic(fmt.Errorf("uniform '%s' of type vec4 should be an float32 array with a multiple of four length", u.name))
			}

			u.fn.Uniform4fv(u.location, vector)
		}

	case stSampler2d:
		textureCtx := target.(*textureContext)
		u.fn.ActiveTexture(Enum(TEXTURE0 + u.textureUnit))
		u.fn.BindTexture(TEXTURE_2D, textureCtx.texture)
		u.fn.Uniform1i(u.location, u.textureUnit)

	default:
		panic(fmt.Errorf("uniform of unsupported type %s", u.shaderType))
	}
}
