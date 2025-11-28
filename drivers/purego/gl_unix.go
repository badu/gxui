//go:build darwin || linux || freebsd || openbsd

package purego

import (
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type Functions struct {
	gpActiveTexture            uintptr
	gpAttachShader             uintptr
	gpBindAttribLocation       uintptr
	gpBindBuffer               uintptr
	gpBindFramebuffer          uintptr
	gpBindRenderbuffer         uintptr
	gpBindTexture              uintptr
	gpBindVertexArray          uintptr
	gpBlendEquationSeparate    uintptr
	gpBlendFuncSeparate        uintptr
	gpBufferData               uintptr
	gpBufferSubData            uintptr
	gpCheckFramebufferStatus   uintptr
	gpClear                    uintptr
	gpColorMask                uintptr
	gpCompileShader            uintptr
	gpCreateProgram            uintptr
	gpCreateShader             uintptr
	gpDeleteBuffers            uintptr
	gpDeleteFramebuffers       uintptr
	gpDeleteProgram            uintptr
	gpDeleteRenderbuffers      uintptr
	gpDeleteShader             uintptr
	gpDeleteTextures           uintptr
	gpDeleteVertexArrays       uintptr
	gpDisable                  uintptr
	gpDisableVertexAttribArray uintptr
	gpDrawElements             uintptr
	gpEnable                   uintptr
	gpEnableVertexAttribArray  uintptr
	gpFlush                    uintptr
	gpFramebufferRenderbuffer  uintptr
	gpFramebufferTexture2D     uintptr
	gpGenBuffers               uintptr
	gpGenFramebuffers          uintptr
	gpGenRenderbuffers         uintptr
	gpGenTextures              uintptr
	gpGenVertexArrays          uintptr
	gpGetError                 uintptr
	gpGetIntegerv              uintptr
	gpGetProgramInfoLog        uintptr
	gpGetProgramiv             uintptr
	gpGetShaderInfoLog         uintptr
	gpGetShaderiv              uintptr
	gpGetUniformLocation       uintptr
	gpIsProgram                uintptr
	gpLinkProgram              uintptr
	gpPixelStorei              uintptr
	gpReadPixels               uintptr
	gpRenderbufferStorage      uintptr
	gpScissor                  uintptr
	gpShaderSource             uintptr
	gpStencilFunc              uintptr
	gpStencilOpSeparate        uintptr
	gpTexImage2D               uintptr
	gpTexParameteri            uintptr
	gpTexSubImage2D            uintptr
	gpUniform1fv               uintptr
	gpUniform1i                uintptr
	gpUniform1iv               uintptr
	gpUniform2fv               uintptr
	gpUniform2iv               uintptr
	gpUniform3fv               uintptr
	gpUniform3iv               uintptr
	gpUniform4fv               uintptr
	gpUniform4iv               uintptr
	gpUniformMatrix2fv         uintptr
	gpUniformMatrix3fv         uintptr
	gpUniformMatrix4fv         uintptr
	gpUseProgram               uintptr
	gpVertexAttribPointer      uintptr
	gpViewport                 uintptr

	isES bool

	glBufferData        uintptr
	glClearColor        uintptr
	glDrawArrays        uintptr
	glUniform1f         uintptr
	glBlendFunc         uintptr
	glGetActiveUniform  uintptr
	glGetActiveAttrib   uintptr
	glGetAttribLocation uintptr
}

// BufferData(target Enum, data []byte, usage Enum)
func (fn *Functions) BufferData(target uint32, data []byte, usage uint32) {
	purego.SyscallN(fn.glBufferData, uintptr(target), uintptr(len(data)), uintptr(unsafe.Pointer(&data[0])), uintptr(usage))
}

// ClearColor(red float32, green float32, blue float32, alpha float32)
func (fn *Functions) ClearColor(red, green, blue, alpha float32) {
	purego.SyscallN(fn.glClearColor, uintptr(red), uintptr(green), uintptr(blue), uintptr(alpha))
}

// DrawArrays(mode Enum, first int, count int)
func (fn *Functions) DrawArrays(mode uint32, first int, count int) {
	purego.SyscallN(fn.glDrawArrays, uintptr(mode), uintptr(first), uintptr(count))
}

// Uniform1fv(dst Uniform, src []float32)
func (fn *Functions) Uniform1f(dst int32, value float32) {
	purego.SyscallN(fn.glUniform1f, uintptr(dst), uintptr(value))
}

// BlendFunc sets the pixel blending factors.
// BlendFunc(sfactor, dfactor Enum)
func (fn *Functions) BlendFunc(sFactor, dFactor uint32) {
	purego.SyscallN(fn.glBlendFunc, uintptr(sFactor), uintptr(dFactor))
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
// GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum)
func (fn *Functions) GetActiveUniform(program uint32, index uint32) (string, int32, uint32) {
	bufSize := int32(256)
	name := make([]byte, bufSize)

	var length, size int32
	var typ uint32

	purego.SyscallN(
		fn.glGetActiveUniform,
		uintptr(program),
		uintptr(index),
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&length)),
		uintptr(unsafe.Pointer(&size)),
		uintptr(unsafe.Pointer(&typ)),
		uintptr(unsafe.Pointer(&name[0])),
	)

	// Convert byte buffer to actual string
	uniformName := string(name[:length])

	return uniformName, size, typ
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
// GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum)
func (fn *Functions) GetActiveAttrib(program uint32, index uint32) (string, int32, uint32) {
	bufSize := int32(256)
	name := make([]byte, bufSize)

	var length, size int32
	var typ uint32

	purego.SyscallN(
		fn.glGetActiveAttrib,
		uintptr(program),
		uintptr(index),
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&length)),
		uintptr(unsafe.Pointer(&size)),
		uintptr(unsafe.Pointer(&typ)),
		uintptr(unsafe.Pointer(&name[0])),
	)

	// Convert only the valid bytes
	attribName := string(name[:length])

	return attribName, size, typ
}

// GetAttribLocation(p Program, name string) Attrib
func (fn *Functions) GetAttribLocation(program uint32, name string) uint32 {
	cname, free := cStr(name)
	defer free()
	ret, _, _ := purego.SyscallN(fn.glGetAttribLocation, uintptr(program), uintptr(unsafe.Pointer(cname)))
	return uint32(ret)
}

func NewFunctions() (*Functions, error) {
	ctx := &Functions{}
	if err := ctx.init(); err != nil {
		return nil, err
	}

	if err := ctx.LoadFunctions(); err != nil {
		return nil, err
	}
	
	return ctx, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (fn *Functions) IsES() bool {
	return fn.isES
}

// ActiveTexture(texture Enum)
func (fn *Functions) ActiveTexture(texture uint32) {
	purego.SyscallN(fn.gpActiveTexture, uintptr(texture))
}

// AttachShader(p Program, s Shader)
func (fn *Functions) AttachShader(program uint32, shader uint32) {
	purego.SyscallN(fn.gpAttachShader, uintptr(program), uintptr(shader))
}

// BindAttribLocation(p Program, a Attrib, name string)
func (fn *Functions) BindAttribLocation(program uint32, index uint32, name string) {
	cname, free := cStr(name)
	defer free()
	purego.SyscallN(fn.gpBindAttribLocation, uintptr(program), uintptr(index), uintptr(unsafe.Pointer(cname)))
}

// BindBuffer(target Enum, b Buffer)
func (fn *Functions) BindBuffer(target uint32, buffer uint32) {
	purego.SyscallN(fn.gpBindBuffer, uintptr(target), uintptr(buffer))
}

// BindFramebuffer(target Enum, fb Framebuffer)
func (fn *Functions) BindFramebuffer(target uint32, framebuffer uint32) {
	purego.SyscallN(fn.gpBindFramebuffer, uintptr(target), uintptr(framebuffer))
}

// BindRenderbuffer(target Enum, fb Renderbuffer)
func (fn *Functions) BindRenderbuffer(target uint32, renderbuffer uint32) {
	purego.SyscallN(fn.gpBindRenderbuffer, uintptr(target), uintptr(renderbuffer))
}

// BindTexture(target Enum, t Texture)
func (fn *Functions) BindTexture(target uint32, texture uint32) {
	purego.SyscallN(fn.gpBindTexture, uintptr(target), uintptr(texture))
}

// BindVertexArray(a VertexArray)
func (fn *Functions) BindVertexArray(array uint32) {
	purego.SyscallN(fn.gpBindVertexArray, uintptr(array))
}

func (fn *Functions) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {
	purego.SyscallN(fn.gpBlendEquationSeparate, uintptr(modeRGB), uintptr(modeAlpha))
}

// BlendFuncSeparate(srcRGB, dstRGB, srcA, dstA Enum)
func (fn *Functions) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	purego.SyscallN(fn.gpBlendFuncSeparate, uintptr(srcRGB), uintptr(dstRGB), uintptr(srcAlpha), uintptr(dstAlpha))
}

func (fn *Functions) BufferInit(target uint32, size int, usage uint32) {
	purego.SyscallN(fn.gpBufferData, uintptr(target), uintptr(size), 0, uintptr(usage))
}

// TODO : use  go:uintptrescapes

// BufferSubData(target Enum, offset int, src []byte)
func (fn *Functions) BufferSubData(target uint32, offset int, data []byte) {
	purego.SyscallN(fn.gpBufferSubData, uintptr(target), uintptr(offset), uintptr(len(data)), uintptr(unsafe.Pointer(&data[0])))
	runtime.KeepAlive(data)
}

// CheckFramebufferStatus(target Enum) Enum
func (fn *Functions) CheckFramebufferStatus(target uint32) uint32 {
	ret, _, _ := purego.SyscallN(fn.gpCheckFramebufferStatus, uintptr(target))
	return uint32(ret)
}

// Clear(mask Enum)
func (fn *Functions) Clear(mask uint32) {
	purego.SyscallN(fn.gpClear, uintptr(mask))
}

func (fn *Functions) ColorMask(red bool, green bool, blue bool, alpha bool) {
	purego.SyscallN(fn.gpColorMask, uintptr(boolToInt(red)), uintptr(boolToInt(green)), uintptr(boolToInt(blue)), uintptr(boolToInt(alpha)))
}

// CompileShader(s Shader)
func (fn *Functions) CompileShader(shader uint32) {
	purego.SyscallN(fn.gpCompileShader, uintptr(shader))
}

// CreateBuffer() Buffer
func (fn *Functions) CreateBuffer() uint32 {
	var buffer uint32
	purego.SyscallN(fn.gpGenBuffers, 1, uintptr(unsafe.Pointer(&buffer)))
	return buffer
}

// CreateFramebuffer() Framebuffer
func (fn *Functions) CreateFramebuffer() uint32 {
	var framebuffer uint32
	purego.SyscallN(fn.gpGenFramebuffers, 1, uintptr(unsafe.Pointer(&framebuffer)))
	return framebuffer
}

// CreateProgram() Program
func (fn *Functions) CreateProgram() uint32 {
	ret, _, _ := purego.SyscallN(fn.gpCreateProgram)
	return uint32(ret)
}

// CreateRenderbuffer() Renderbuffer
func (fn *Functions) CreateRenderbuffer() uint32 {
	var renderbuffer uint32
	purego.SyscallN(fn.gpGenRenderbuffers, 1, uintptr(unsafe.Pointer(&renderbuffer)))
	return renderbuffer
}

// CreateShader(ty Enum) Shader
func (fn *Functions) CreateShader(xtype uint32) uint32 {
	ret, _, _ := purego.SyscallN(fn.gpCreateShader, uintptr(xtype))
	return uint32(ret)
}

// CreateTexture() Texture
func (fn *Functions) CreateTexture() uint32 {
	var texture uint32
	purego.SyscallN(fn.gpGenTextures, 1, uintptr(unsafe.Pointer(&texture)))
	return texture
}

// CreateVertexArray() VertexArray
func (fn *Functions) CreateVertexArray() uint32 {
	var array uint32
	purego.SyscallN(fn.gpGenVertexArrays, 1, uintptr(unsafe.Pointer(&array)))
	return array
}

// DeleteBuffer(v Buffer)
func (fn *Functions) DeleteBuffer(buffer uint32) {
	purego.SyscallN(fn.gpDeleteBuffers, 1, uintptr(unsafe.Pointer(&buffer)))
}

// DeleteFramebuffer(v Framebuffer)
func (fn *Functions) DeleteFramebuffer(framebuffer uint32) {
	purego.SyscallN(fn.gpDeleteFramebuffers, 1, uintptr(unsafe.Pointer(&framebuffer)))
}

// DeleteProgram(p Program)
func (fn *Functions) DeleteProgram(program uint32) {
	purego.SyscallN(fn.gpDeleteProgram, uintptr(program))
}

// DeleteRenderbuffer(v Renderbuffer)
func (fn *Functions) DeleteRenderbuffer(renderbuffer uint32) {
	purego.SyscallN(fn.gpDeleteRenderbuffers, 1, uintptr(unsafe.Pointer(&renderbuffer)))
}

// DeleteShader(s Shader)
func (fn *Functions) DeleteShader(shader uint32) {
	purego.SyscallN(fn.gpDeleteShader, uintptr(shader))
}

// DeleteTexture(v Texture)
func (fn *Functions) DeleteTexture(texture uint32) {
	purego.SyscallN(fn.gpDeleteTextures, 1, uintptr(unsafe.Pointer(&texture)))
}

// DeleteVertexArray(array VertexArray)
func (fn *Functions) DeleteVertexArray(array uint32) {
	purego.SyscallN(fn.gpDeleteVertexArrays, 1, uintptr(unsafe.Pointer(&array)))
}

// Disable(cap Enum)
func (fn *Functions) Disable(cap uint32) {
	purego.SyscallN(fn.gpDisable, uintptr(cap))
}

// DisableVertexAttribArray(a Attrib)
func (fn *Functions) DisableVertexAttribArray(index uint32) {
	purego.SyscallN(fn.gpDisableVertexAttribArray, uintptr(index))
}

// DrawElements(mode Enum, count int, ty Enum, offset int)
func (fn *Functions) DrawElements(mode uint32, count int32, xtype uint32, offset int) {
	purego.SyscallN(fn.gpDrawElements, uintptr(mode), uintptr(count), uintptr(xtype), uintptr(offset))
}

// Enable(cap Enum)
func (fn *Functions) Enable(cap uint32) {
	purego.SyscallN(fn.gpEnable, uintptr(cap))
}

// EnableVertexAttribArray(a Attrib)
func (fn *Functions) EnableVertexAttribArray(index uint32) {
	purego.SyscallN(fn.gpEnableVertexAttribArray, uintptr(index))
}

// Flush()
func (fn *Functions) Flush() {
	purego.SyscallN(fn.gpFlush)
}

// FramebufferRenderbuffer(target, attachment, renderbuffertarget Enum, renderbuffer Renderbuffer)
func (fn *Functions) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	purego.SyscallN(fn.gpFramebufferRenderbuffer, uintptr(target), uintptr(attachment), uintptr(renderbuffertarget), uintptr(renderbuffer))
}

// FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int)
func (fn *Functions) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	purego.SyscallN(fn.gpFramebufferTexture2D, uintptr(target), uintptr(attachment), uintptr(textarget), uintptr(texture), uintptr(level))
}

// GetError()
func (fn *Functions) GetError() uint32 {
	ret, _, _ := purego.SyscallN(fn.gpGetError)
	return uint32(ret)
}

func (fn *Functions) GetExtension(name string) any {
	return nil
}

// GetInteger(pname Enum) int
func (fn *Functions) GetInteger(pname uint32) int {
	var dst int32
	purego.SyscallN(fn.gpGetIntegerv, uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetProgramInfoLog(p Program) string
func (fn *Functions) GetProgramInfoLog(program uint32) string {
	bufSize := fn.GetProgrami(program, INFO_LOG_LENGTH)
	if bufSize == 0 {
		return ""
	}
	infoLog := make([]byte, bufSize)
	purego.SyscallN(fn.gpGetProgramInfoLog, uintptr(program), uintptr(bufSize), 0, uintptr(unsafe.Pointer(&infoLog[0])))
	return string(infoLog)
}

// GetProgrami(p Program, pname Enum) int
func (fn *Functions) GetProgrami(program uint32, pname uint32) int {
	var dst int32
	purego.SyscallN(fn.gpGetProgramiv, uintptr(program), uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetShaderInfoLog(s Shader) string
func (fn *Functions) GetShaderInfoLog(shader uint32) string {
	bufSize := fn.GetShaderi(shader, INFO_LOG_LENGTH)
	if bufSize == 0 {
		return ""
	}
	infoLog := make([]byte, bufSize)
	purego.SyscallN(fn.gpGetShaderInfoLog, uintptr(shader), uintptr(bufSize), 0, uintptr(unsafe.Pointer(&infoLog[0])))
	return string(infoLog)
}

// GetShaderi(s Shader, pname Enum) int
func (fn *Functions) GetShaderi(shader uint32, pname uint32) int {
	var dst int32
	purego.SyscallN(fn.gpGetShaderiv, uintptr(shader), uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetUniformLocation(p Program, name string) Uniform
func (fn *Functions) GetUniformLocation(program uint32, name string) int32 {
	cname, free := cStr(name)
	defer free()
	ret, _, _ := purego.SyscallN(fn.gpGetUniformLocation, uintptr(program), uintptr(unsafe.Pointer(cname)))
	return int32(ret)
}

func (fn *Functions) IsProgram(program uint32) bool {
	ret, _, _ := purego.SyscallN(fn.gpIsProgram, uintptr(program))
	return byte(ret) != 0
}

// LinkProgram(p Program)
func (fn *Functions) LinkProgram(program uint32) {
	purego.SyscallN(fn.gpLinkProgram, uintptr(program))
}

// PixelStorei(pname Enum, param int)
func (fn *Functions) PixelStorei(pname uint32, param int32) {
	purego.SyscallN(fn.gpPixelStorei, uintptr(pname), uintptr(param))
}

// ReadPixels(x, y, width, height int, format, ty Enum, data []byte)
func (fn *Functions) ReadPixels(dst []byte, x int32, y int32, width int32, height int32, format uint32, xtype uint32) {
	purego.SyscallN(fn.gpReadPixels, uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(&dst[0])))
}

// RenderbufferStorage(target, internalformat Enum, width, height int)
func (fn *Functions) RenderbufferStorage(target uint32, internalformat uint32, width int32, height int32) {
	purego.SyscallN(fn.gpRenderbufferStorage, uintptr(target), uintptr(internalformat), uintptr(width), uintptr(height))
}

// Scissor(x, y, width, height int32)
func (fn *Functions) Scissor(x int32, y int32, width int32, height int32) {
	purego.SyscallN(fn.gpScissor, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}

// ShaderSource(s Shader, src string)
func (fn *Functions) ShaderSource(shader uint32, xstring string) {
	cstring, free := cStr(xstring)
	defer free()
	purego.SyscallN(fn.gpShaderSource, uintptr(shader), 1, uintptr(unsafe.Pointer(&cstring)), 0)
}

func (fn *Functions) StencilFunc(xfunc uint32, ref int32, mask uint32) {
	purego.SyscallN(fn.gpStencilFunc, uintptr(xfunc), uintptr(ref), uintptr(mask))
}

func (fn *Functions) StencilOpSeparate(face uint32, fail uint32, zfail uint32, zpass uint32) {
	purego.SyscallN(fn.gpStencilOpSeparate, uintptr(face), uintptr(fail), uintptr(zfail), uintptr(zpass))
}

// TexImage2D(target Enum, level int, width int, height int, format Enum, ty Enum, data []byte)
func (fn *Functions) TexImage2D(target uint32, level int32, width int32, height int32, format uint32, xtype uint32, pixels []byte) {
	var ptr *byte
	if len(pixels) > 0 {
		ptr = &pixels[0]
	}
	purego.SyscallN(fn.gpTexImage2D, uintptr(target), uintptr(level), uintptr(format), uintptr(width), uintptr(height), 0, uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(ptr)))
	runtime.KeepAlive(pixels)
}

// TexParameteri(target, pname Enum, param int)
func (fn *Functions) TexParameteri(target uint32, pname uint32, param int32) {
	purego.SyscallN(fn.gpTexParameteri, uintptr(target), uintptr(pname), uintptr(param))
}

// TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte)
func (fn *Functions) TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels []byte) {
	purego.SyscallN(fn.gpTexSubImage2D, uintptr(target), uintptr(level), uintptr(xoffset), uintptr(yoffset), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(&pixels[0])))
	runtime.KeepAlive(pixels)
}

// Uniform1fv(dst Uniform, src []float32)
func (fn *Functions) Uniform1fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniform1fv, uintptr(location), uintptr(len(value)), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform1i(dst Uniform, v int)
func (fn *Functions) Uniform1i(location int32, v0 int32) {
	purego.SyscallN(fn.gpUniform1i, uintptr(location), uintptr(v0))
}

func (fn *Functions) Uniform1iv(location int32, value []int32) {
	purego.SyscallN(fn.gpUniform1iv, uintptr(location), uintptr(len(value)), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform2fv(dst Uniform, src []float32)
func (fn *Functions) Uniform2fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniform2fv, uintptr(location), uintptr(len(value)/2), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (fn *Functions) Uniform2iv(location int32, value []int32) {
	purego.SyscallN(fn.gpUniform2iv, uintptr(location), uintptr(len(value)/2), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform3fv(dst Uniform, src []float32)
func (fn *Functions) Uniform3fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniform3fv, uintptr(location), uintptr(len(value)/3), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (fn *Functions) Uniform3iv(location int32, value []int32) {
	purego.SyscallN(fn.gpUniform3iv, uintptr(location), uintptr(len(value)/3), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform4fv(dst Uniform, src []float32)
func (fn *Functions) Uniform4fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniform4fv, uintptr(location), uintptr(len(value)/4), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (fn *Functions) Uniform4iv(location int32, value []int32) {
	purego.SyscallN(fn.gpUniform4iv, uintptr(location), uintptr(len(value)/4), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix2fv(dst Uniform, src []float32)
func (fn *Functions) UniformMatrix2fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniformMatrix2fv, uintptr(location), uintptr(len(value)/4), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix3fv(dst Uniform, src []float32)
func (fn *Functions) UniformMatrix3fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniformMatrix3fv, uintptr(location), uintptr(len(value)/9), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix4fv(dst Uniform, src []float32)
func (fn *Functions) UniformMatrix4fv(location int32, value []float32) {
	purego.SyscallN(fn.gpUniformMatrix4fv, uintptr(location), uintptr(len(value)/16), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UseProgram(p Program)
func (fn *Functions) UseProgram(program uint32) {
	purego.SyscallN(fn.gpUseProgram, uintptr(program))
}

// VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int)
func (fn *Functions) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	purego.SyscallN(fn.gpVertexAttribPointer, uintptr(index), uintptr(size), uintptr(xtype), uintptr(boolToInt(normalized)), uintptr(stride), uintptr(offset))
}

// Viewport(x int, y int, width int, height int)
func (fn *Functions) Viewport(x int32, y int32, width int32, height int32) {
	purego.SyscallN(fn.gpViewport, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}

func (fn *Functions) LoadFunctions() error {
	g := procAddressGetter{ctx: fn}

	fn.gpActiveTexture = g.get("glActiveTexture")
	fn.gpAttachShader = g.get("glAttachShader")
	fn.gpBindAttribLocation = g.get("glBindAttribLocation")
	fn.gpBindBuffer = g.get("glBindBuffer")
	fn.gpBindFramebuffer = g.get("glBindFramebuffer")
	fn.gpBindRenderbuffer = g.get("glBindRenderbuffer")
	fn.gpBindTexture = g.get("glBindTexture")
	fn.gpBindVertexArray = g.get("glBindVertexArray")
	fn.gpBlendEquationSeparate = g.get("glBlendEquationSeparate")
	fn.gpBlendFuncSeparate = g.get("glBlendFuncSeparate")
	fn.gpBufferData = g.get("glBufferData")
	fn.gpBufferSubData = g.get("glBufferSubData")
	fn.gpCheckFramebufferStatus = g.get("glCheckFramebufferStatus")
	fn.gpClear = g.get("glClear")
	fn.gpColorMask = g.get("glColorMask")
	fn.gpCompileShader = g.get("glCompileShader")
	fn.gpCreateProgram = g.get("glCreateProgram")
	fn.gpCreateShader = g.get("glCreateShader")
	fn.gpDeleteBuffers = g.get("glDeleteBuffers")
	fn.gpDeleteFramebuffers = g.get("glDeleteFramebuffers")
	fn.gpDeleteProgram = g.get("glDeleteProgram")
	fn.gpDeleteRenderbuffers = g.get("glDeleteRenderbuffers")
	fn.gpDeleteShader = g.get("glDeleteShader")
	fn.gpDeleteTextures = g.get("glDeleteTextures")
	fn.gpDeleteVertexArrays = g.get("glDeleteVertexArrays")
	fn.gpDisable = g.get("glDisable")
	fn.gpDisableVertexAttribArray = g.get("glDisableVertexAttribArray")
	fn.gpDrawElements = g.get("glDrawElements")
	fn.gpEnable = g.get("glEnable")
	fn.gpEnableVertexAttribArray = g.get("glEnableVertexAttribArray")
	fn.gpFlush = g.get("glFlush")
	fn.gpFramebufferRenderbuffer = g.get("glFramebufferRenderbuffer")
	fn.gpFramebufferTexture2D = g.get("glFramebufferTexture2D")
	fn.gpGenBuffers = g.get("glGenBuffers")
	fn.gpGenFramebuffers = g.get("glGenFramebuffers")
	fn.gpGenRenderbuffers = g.get("glGenRenderbuffers")
	fn.gpGenTextures = g.get("glGenTextures")
	fn.gpGenVertexArrays = g.get("glGenVertexArrays")
	fn.gpGetError = g.get("glGetError")
	fn.gpGetIntegerv = g.get("glGetIntegerv")
	fn.gpGetProgramInfoLog = g.get("glGetProgramInfoLog")
	fn.gpGetProgramiv = g.get("glGetProgramiv")
	fn.gpGetShaderInfoLog = g.get("glGetShaderInfoLog")
	fn.gpGetShaderiv = g.get("glGetShaderiv")
	fn.gpGetUniformLocation = g.get("glGetUniformLocation")
	fn.gpIsProgram = g.get("glIsProgram")
	fn.gpLinkProgram = g.get("glLinkProgram")
	fn.gpPixelStorei = g.get("glPixelStorei")
	fn.gpReadPixels = g.get("glReadPixels")
	fn.gpRenderbufferStorage = g.get("glRenderbufferStorage")
	fn.gpScissor = g.get("glScissor")
	fn.gpShaderSource = g.get("glShaderSource")
	fn.gpStencilFunc = g.get("glStencilFunc")
	fn.gpStencilOpSeparate = g.get("glStencilOpSeparate")
	fn.gpTexImage2D = g.get("glTexImage2D")
	fn.gpTexParameteri = g.get("glTexParameteri")
	fn.gpTexSubImage2D = g.get("glTexSubImage2D")
	fn.gpUniform1fv = g.get("glUniform1fv")
	fn.gpUniform1i = g.get("glUniform1i")
	fn.gpUniform1iv = g.get("glUniform1iv")
	fn.gpUniform2fv = g.get("glUniform2fv")
	fn.gpUniform2iv = g.get("glUniform2iv")
	fn.gpUniform3fv = g.get("glUniform3fv")
	fn.gpUniform3iv = g.get("glUniform3iv")
	fn.gpUniform4fv = g.get("glUniform4fv")
	fn.gpUniform4iv = g.get("glUniform4iv")
	fn.gpUniformMatrix2fv = g.get("glUniformMatrix2fv")
	fn.gpUniformMatrix3fv = g.get("glUniformMatrix3fv")
	fn.gpUniformMatrix4fv = g.get("glUniformMatrix4fv")
	fn.gpUseProgram = g.get("glUseProgram")
	fn.gpVertexAttribPointer = g.get("glVertexAttribPointer")
	fn.gpViewport = g.get("glViewport")

	fn.glBufferData = g.get("glBufferData")
	fn.glClearColor = g.get("glClearColor")
	fn.glDrawArrays = g.get("glDrawArrays")
	fn.glUniform1f = g.get("glUniform1f")
	fn.glBlendFunc = g.get("glBlendFunc")
	fn.glGetActiveUniform = g.get("glGetActiveUniform")
	fn.glGetActiveAttrib = g.get("glGetActiveAttrib")
	fn.glGetAttribLocation = g.get("glGetAttribLocation")

	return g.error()
}

// cStr takes a Go string (with or without null-termination)
// and returns the C counterpart.
//
// The returned free function must be called once you are done using the string
// in order to free the memory.
func cStr(str string) (cstr *byte, free func()) {
	bs := []byte(str)
	if len(bs) == 0 || bs[len(bs)-1] != 0 {
		bs = append(bs, 0)
	}
	return &bs[0], func() {
		runtime.KeepAlive(bs)
		bs = nil
	}
}
