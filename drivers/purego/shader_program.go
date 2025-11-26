package purego

import (
	"fmt"
)

type uniformBindings map[string]interface{}

type shaderProgram struct {
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

func newShaderProgram(ctx *context, vertSource, fragSource string) *shaderProgram {
	vertex := compile(ctx.fn, vertSource, VERTEX_SHADER)
	fragment := compile(ctx.fn, fragSource, FRAGMENT_SHADER)

	program := ctx.fn.CreateProgram()
	ctx.fn.AttachShader(program, vertex)
	ctx.fn.AttachShader(program, fragment)
	ctx.fn.LinkProgram(program)

	if ctx.fn.GetProgrami(program, LINK_STATUS) != TRUE {
		panic(ctx.fn.GetProgramInfoLog(program))
	}

	ctx.fn.UseProgram(program)

	checkError(ctx.fn)

	uniformCount := ctx.fn.GetProgrami(program, ACTIVE_UNIFORMS)
	uniforms := make([]shaderUniform, uniformCount)
	textureUnit := 0
	for index := range uniforms {
		name, size, shaderType := ctx.fn.GetActiveUniform(program, uint32(index))
		location := ctx.fn.GetUniformLocation(program, name)
		uniforms[index] = shaderUniform{
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

	attributeCount := ctx.fn.GetProgrami(program, ACTIVE_ATTRIBUTES)
	attributes := make([]shaderAttribute, attributeCount)
	for index := range attributes {
		name, size, shaderType := ctx.fn.GetActiveAttrib(program, uint32(index))
		location := ctx.fn.GetAttribLocation(program, name)
		attributes[index] = shaderAttribute{
			name:       name,
			size:       size,
			shaderType: shaderDataType(shaderType),
			glAttr:     location,
		}
	}

	ctx.stats.shaderProgramCount++

	return &shaderProgram{
		program:    program,
		uniforms:   uniforms,
		attributes: attributes,
	}
}

func (s *shaderProgram) destroy(ctx *context) {
	ctx.fn.DeleteProgram(s.program)
	s.program = Program{}
	// TODO: Delete shaders.
	ctx.stats.shaderProgramCount--
}

func (s *shaderProgram) bind(ctx *context, buffer *vertexBuffer, uniforms uniformBindings) {
	ctx.fn.UseProgram(s.program)

	for _, attr := range s.attributes {
		if attr.name == "" {
			panic("empty attribute name")
		}

		vertex, found := buffer.streams[attr.name]
		if !found {
			panic(fmt.Errorf("VertexBuffer missing required stream %q", attr.name))
		}

		if attr.shaderType != vertex.shaderType {
			panic(fmt.Errorf("attribute %q type %q does not match stream type %q", attr.name, attr.shaderType, vertex.shaderType))
		}

		elementCount := attr.shaderType.vectorElementCount()
		elementTy := attr.shaderType.vectorElementType()
		ctx.getOrCreateVertexStreamContext(vertex).bind(ctx.fn)
		attr.enableArray(ctx.fn)
		attr.attribPointer(ctx.fn, int32(elementCount), uint32(elementTy), false, 0, 0)
	}

	for _, uni := range s.uniforms {
		if uni.name == "" {
			panic("empty uniform name")
		}

		uniform, found := uniforms[uni.name]
		if !found {
			panic(fmt.Errorf("uniforms missing %q", uni.name))
		}

		uni.bind(ctx.fn, uniform)
	}

	checkError(ctx.fn)
}

func (s *shaderProgram) unbind(fn *Functions) {
	for _, a := range s.attributes {
		if a.name == "" {
			panic("empty attribute name on unbind")
		}

		a.disableArray(fn)
	}

	checkError(fn)
}
