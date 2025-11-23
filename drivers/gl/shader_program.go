// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/goxjs/gl"
)

type uniformBindings map[string]interface{}

type shaderProgram struct {
	uniforms   []shaderUniform
	attributes []shaderAttribute
	program    gl.Program
}

func compile(source string, shaderType int) gl.Shader {
	shader := gl.CreateShader(gl.Enum(shaderType))
	gl.ShaderSource(shader, source)

	gl.CompileShader(shader)
	if gl.GetShaderi(shader, gl.COMPILE_STATUS) != gl.TRUE {
		panic(gl.GetShaderInfoLog(shader))
	}

	checkError()

	return shader
}

func newShaderProgram(ctx *context, vertSource, fragSource string) *shaderProgram {
	vertex := compile(vertSource, gl.VERTEX_SHADER)
	fragment := compile(fragSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	if gl.GetProgrami(program, gl.LINK_STATUS) != gl.TRUE {
		panic(gl.GetProgramInfoLog(program))
	}

	gl.UseProgram(program)

	checkError()

	uniformCount := gl.GetProgrami(program, gl.ACTIVE_UNIFORMS)
	uniforms := make([]shaderUniform, uniformCount)
	textureUnit := 0
	for index := range uniforms {
		name, size, shaderType := gl.GetActiveUniform(program, uint32(index))
		location := gl.GetUniformLocation(program, name)
		uniforms[index] = shaderUniform{
			name:        name,
			size:        size,
			shaderType:  shaderDataType(shaderType),
			location:    location,
			textureUnit: textureUnit,
		}
		if shaderType == gl.SAMPLER_2D {
			textureUnit++
		}
	}

	attributeCount := gl.GetProgrami(program, gl.ACTIVE_ATTRIBUTES)
	attributes := make([]shaderAttribute, attributeCount)
	for index := range attributes {
		name, size, shaderType := gl.GetActiveAttrib(program, uint32(index))
		location := gl.GetAttribLocation(program, name)
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
	gl.DeleteProgram(s.program)
	s.program = gl.Program{}
	// TODO: Delete shaders.
	ctx.stats.shaderProgramCount--
}

func (s *shaderProgram) bind(ctx *context, buffer *vertexBuffer, uniforms uniformBindings) {
	gl.UseProgram(s.program)

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
	checkError()
}

func (s *shaderProgram) unbind(ctx *context) {
	for _, a := range s.attributes {
		a.disableArray()
	}
	checkError()
}
