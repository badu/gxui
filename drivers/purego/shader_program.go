package purego

import (
	"fmt"
)

type uniformBindings map[string]interface{}

type shaderProgram struct {
	fn         *Functions
	uniforms   []shaderUniform
	attributes []shaderAttribute
	program    Program
}

func compile(fn *Functions, source string, shaderType int) Shader {
	shader := fn.CreateShader(Enum(shaderType))
	fn.ShaderSource(shader, source)

	fn.CompileShader(shader)
	if fn.GetShaderi(shader, COMPILE_STATUS) != TRUE {
		panic(fn.GetShaderInfoLog(shader))
	}

	checkError(fn)

	return shader
}

func newShaderProgram(fn *Functions, ctx *context, vertSource, fragSource string) *shaderProgram {
	vertex := compile(fn, vertSource, VERTEX_SHADER)
	fragment := compile(fn, fragSource, FRAGMENT_SHADER)

	program := fn.CreateProgram()
	fn.AttachShader(program, vertex)
	fn.AttachShader(program, fragment)
	fn.LinkProgram(program)

	if fn.GetProgrami(program, LINK_STATUS) != TRUE {
		panic(fn.GetProgramInfoLog(program))
	}

	fn.UseProgram(program)

	checkError(fn)

	uniformCount := fn.GetProgrami(program, ACTIVE_UNIFORMS)
	uniforms := make([]shaderUniform, uniformCount)
	textureUnit := 0
	for index := range uniforms {
		name, size, shaderType := fn.GetActiveUniform(program, uint32(index))
		location := fn.GetUniformLocation(program, name)
		uniforms[index] = shaderUniform{
			fn:          fn,
			name:        name,
			size:        size,
			shaderType:  shaderDataType(shaderType),
			location:    location,
			textureUnit: textureUnit,
		}
		if shaderType == SAMPLER_2D {
			textureUnit++
		}
	}

	attributeCount := fn.GetProgrami(program, ACTIVE_ATTRIBUTES)
	attributes := make([]shaderAttribute, attributeCount)
	for index := range attributes {
		name, size, shaderType := fn.GetActiveAttrib(program, uint32(index))
		location := fn.GetAttribLocation(program, name)
		attributes[index] = shaderAttribute{
			fn:         fn,
			name:       name,
			size:       size,
			shaderType: shaderDataType(shaderType),
			glAttr:     location,
		}
	}

	ctx.stats.shaderProgramCount++

	return &shaderProgram{
		fn:         fn,
		program:    program,
		uniforms:   uniforms,
		attributes: attributes,
	}
}

func (s *shaderProgram) destroy(ctx *context) {
	s.fn.DeleteProgram(s.program)
	s.program = Program{}
	// TODO: Delete shaders.
	ctx.stats.shaderProgramCount--
}

func (s *shaderProgram) bind(ctx *context, buffer *vertexBuffer, uniforms uniformBindings) {
	s.fn.UseProgram(s.program)

	for _, attr := range s.attributes {
		vertex, found := buffer.streams[attr.name]
		if !found {
			panic(fmt.Errorf("VertexBuffer missing required stream '%s'", attr.name))
		}

		if attr.shaderType != vertex.shaderType {
			panic(fmt.Errorf("attribute '%s' type '%s' does not match stream type '%s'", attr.name, attr.shaderType, vertex.shaderType))
		}

		elementCount := attr.shaderType.vectorElementCount()
		elementTy := attr.shaderType.vectorElementType()
		ctx.getOrCreateVertexStreamContext(vertex).bind()
		attr.enableArray()
		attr.attribPointer(int32(elementCount), uint32(elementTy), false, 0, 0)
	}

	for _, uni := range s.uniforms {
		uniform, found := uniforms[uni.name]
		if !found {
			panic(fmt.Errorf("uniforms missing '%s'", uni.name))
		}
		uni.bind(uniform)
	}
	checkError(s.fn)
}

func (s *shaderProgram) unbind(ctx *context) {
	for _, a := range s.attributes {
		a.disableArray()
	}
	checkError(s.fn)
}
