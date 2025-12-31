//go:build darwin || linux || freebsd || openbsd || windows

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

	libGL             uintptr
	libGLES           uintptr
	glXGetProcAddress func(name string) uintptr
}

// BufferData(target Enum, data []byte, usage Enum)
func (f *Functions) BufferData(target uint32, data []byte, usage uint32) {
	purego.SyscallN(f.glBufferData, uintptr(target), uintptr(len(data)), uintptr(unsafe.Pointer(&data[0])), uintptr(usage))
}

// ClearColor(red float32, green float32, blue float32, alpha float32)
func (f *Functions) ClearColor(red, green, blue, alpha float32) {
	purego.SyscallN(f.glClearColor, uintptr(red), uintptr(green), uintptr(blue), uintptr(alpha))
}

// DrawArrays(mode Enum, first int, count int)
func (f *Functions) DrawArrays(mode uint32, first int, count int) {
	purego.SyscallN(f.glDrawArrays, uintptr(mode), uintptr(first), uintptr(count))
}

// Uniform1fv(dst Uniform, src []float32)
func (f *Functions) Uniform1f(dst int32, value float32) {
	purego.SyscallN(f.glUniform1f, uintptr(dst), uintptr(value))
}

// BlendFunc sets the pixel blending factors.
// BlendFunc(sfactor, dfactor Enum)
func (f *Functions) BlendFunc(sFactor, dFactor uint32) {
	purego.SyscallN(f.glBlendFunc, uintptr(sFactor), uintptr(dFactor))
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
// GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum)
func (f *Functions) GetActiveUniform(program uint32, index uint32) (string, int32, uint32) {
	bufSize := int32(256)
	name := make([]byte, bufSize)

	var length, size int32
	var typ uint32

	purego.SyscallN(
		f.glGetActiveUniform,
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
func (f *Functions) GetActiveAttrib(program uint32, index uint32) (string, int32, uint32) {
	bufSize := int32(256)
	name := make([]byte, bufSize)

	var length, size int32
	var typ uint32

	purego.SyscallN(
		f.glGetActiveAttrib,
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
func (f *Functions) GetAttribLocation(program uint32, name string) uint32 {
	cname, free := cStr(name)
	defer free()
	ret, _, _ := purego.SyscallN(f.glGetAttribLocation, uintptr(program), uintptr(unsafe.Pointer(cname)))
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

func (f *Functions) IsES() bool {
	return f.isES
}

// ActiveTexture(texture Enum)
func (f *Functions) ActiveTexture(texture uint32) {
	purego.SyscallN(f.gpActiveTexture, uintptr(texture))
}

// AttachShader(p Program, s Shader)
func (f *Functions) AttachShader(program uint32, shader uint32) {
	purego.SyscallN(f.gpAttachShader, uintptr(program), uintptr(shader))
}

// BindAttribLocation(p Program, a Attrib, name string)
func (f *Functions) BindAttribLocation(program uint32, index uint32, name string) {
	cname, free := cStr(name)
	defer free()
	purego.SyscallN(f.gpBindAttribLocation, uintptr(program), uintptr(index), uintptr(unsafe.Pointer(cname)))
}

// BindBuffer(target Enum, b Buffer)
func (f *Functions) BindBuffer(target uint32, buffer uint32) {
	purego.SyscallN(f.gpBindBuffer, uintptr(target), uintptr(buffer))
}

// BindFramebuffer(target Enum, fb Framebuffer)
func (f *Functions) BindFramebuffer(target uint32, framebuffer uint32) {
	purego.SyscallN(f.gpBindFramebuffer, uintptr(target), uintptr(framebuffer))
}

// BindRenderbuffer(target Enum, fb Renderbuffer)
func (f *Functions) BindRenderbuffer(target uint32, renderbuffer uint32) {
	purego.SyscallN(f.gpBindRenderbuffer, uintptr(target), uintptr(renderbuffer))
}

// BindTexture(target Enum, t Texture)
func (f *Functions) BindTexture(target uint32, texture uint32) {
	purego.SyscallN(f.gpBindTexture, uintptr(target), uintptr(texture))
}

// BindVertexArray(a VertexArray)
func (f *Functions) BindVertexArray(array uint32) {
	purego.SyscallN(f.gpBindVertexArray, uintptr(array))
}

func (f *Functions) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {
	purego.SyscallN(f.gpBlendEquationSeparate, uintptr(modeRGB), uintptr(modeAlpha))
}

// BlendFuncSeparate(srcRGB, dstRGB, srcA, dstA Enum)
func (f *Functions) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	purego.SyscallN(f.gpBlendFuncSeparate, uintptr(srcRGB), uintptr(dstRGB), uintptr(srcAlpha), uintptr(dstAlpha))
}

func (f *Functions) BufferInit(target uint32, size int, usage uint32) {
	purego.SyscallN(f.gpBufferData, uintptr(target), uintptr(size), 0, uintptr(usage))
}

// TODO : use  go:uintptrescapes

// BufferSubData(target Enum, offset int, src []byte)
func (f *Functions) BufferSubData(target uint32, offset int, data []byte) {
	purego.SyscallN(f.gpBufferSubData, uintptr(target), uintptr(offset), uintptr(len(data)), uintptr(unsafe.Pointer(&data[0])))
	runtime.KeepAlive(data)
}

// CheckFramebufferStatus(target Enum) Enum
func (f *Functions) CheckFramebufferStatus(target uint32) uint32 {
	ret, _, _ := purego.SyscallN(f.gpCheckFramebufferStatus, uintptr(target))
	return uint32(ret)
}

// Clear(mask Enum)
func (f *Functions) Clear(mask uint32) {
	purego.SyscallN(f.gpClear, uintptr(mask))
}

func (f *Functions) ColorMask(red bool, green bool, blue bool, alpha bool) {
	purego.SyscallN(f.gpColorMask, uintptr(boolToInt(red)), uintptr(boolToInt(green)), uintptr(boolToInt(blue)), uintptr(boolToInt(alpha)))
}

// CompileShader(s Shader)
func (f *Functions) CompileShader(shader uint32) {
	purego.SyscallN(f.gpCompileShader, uintptr(shader))
}

// CreateBuffer() Buffer
func (f *Functions) CreateBuffer() uint32 {
	var buffer uint32
	purego.SyscallN(f.gpGenBuffers, 1, uintptr(unsafe.Pointer(&buffer)))
	return buffer
}

// CreateFramebuffer() Framebuffer
func (f *Functions) CreateFramebuffer() uint32 {
	var framebuffer uint32
	purego.SyscallN(f.gpGenFramebuffers, 1, uintptr(unsafe.Pointer(&framebuffer)))
	return framebuffer
}

// CreateProgram() Program
func (f *Functions) CreateProgram() uint32 {
	ret, _, _ := purego.SyscallN(f.gpCreateProgram)
	return uint32(ret)
}

// CreateRenderbuffer() Renderbuffer
func (f *Functions) CreateRenderbuffer() uint32 {
	var renderbuffer uint32
	purego.SyscallN(f.gpGenRenderbuffers, 1, uintptr(unsafe.Pointer(&renderbuffer)))
	return renderbuffer
}

// CreateShader(ty Enum) Shader
func (f *Functions) CreateShader(xtype uint32) uint32 {
	ret, _, _ := purego.SyscallN(f.gpCreateShader, uintptr(xtype))
	return uint32(ret)
}

// CreateTexture() Texture
func (f *Functions) CreateTexture() uint32 {
	var texture uint32
	purego.SyscallN(f.gpGenTextures, 1, uintptr(unsafe.Pointer(&texture)))
	return texture
}

// CreateVertexArray() VertexArray
func (f *Functions) CreateVertexArray() uint32 {
	var array uint32
	purego.SyscallN(f.gpGenVertexArrays, 1, uintptr(unsafe.Pointer(&array)))
	return array
}

// DeleteBuffer(v Buffer)
func (f *Functions) DeleteBuffer(buffer uint32) {
	purego.SyscallN(f.gpDeleteBuffers, 1, uintptr(unsafe.Pointer(&buffer)))
}

// DeleteFramebuffer(v Framebuffer)
func (f *Functions) DeleteFramebuffer(framebuffer uint32) {
	purego.SyscallN(f.gpDeleteFramebuffers, 1, uintptr(unsafe.Pointer(&framebuffer)))
}

// DeleteProgram(p Program)
func (f *Functions) DeleteProgram(program uint32) {
	purego.SyscallN(f.gpDeleteProgram, uintptr(program))
}

// DeleteRenderbuffer(v Renderbuffer)
func (f *Functions) DeleteRenderbuffer(renderbuffer uint32) {
	purego.SyscallN(f.gpDeleteRenderbuffers, 1, uintptr(unsafe.Pointer(&renderbuffer)))
}

// DeleteShader(s Shader)
func (f *Functions) DeleteShader(shader uint32) {
	purego.SyscallN(f.gpDeleteShader, uintptr(shader))
}

// DeleteTexture(v Texture)
func (f *Functions) DeleteTexture(texture uint32) {
	purego.SyscallN(f.gpDeleteTextures, 1, uintptr(unsafe.Pointer(&texture)))
}

// DeleteVertexArray(array VertexArray)
func (f *Functions) DeleteVertexArray(array uint32) {
	purego.SyscallN(f.gpDeleteVertexArrays, 1, uintptr(unsafe.Pointer(&array)))
}

// Disable(cap Enum)
func (f *Functions) Disable(cap uint32) {
	purego.SyscallN(f.gpDisable, uintptr(cap))
}

// DisableVertexAttribArray(a Attrib)
func (f *Functions) DisableVertexAttribArray(index uint32) {
	purego.SyscallN(f.gpDisableVertexAttribArray, uintptr(index))
}

// DrawElements(mode Enum, count int, ty Enum, offset int)
func (f *Functions) DrawElements(mode uint32, count int32, xtype uint32, offset int) {
	purego.SyscallN(f.gpDrawElements, uintptr(mode), uintptr(count), uintptr(xtype), uintptr(offset))
}

// Enable(cap Enum)
func (f *Functions) Enable(cap uint32) {
	purego.SyscallN(f.gpEnable, uintptr(cap))
}

// EnableVertexAttribArray(a Attrib)
func (f *Functions) EnableVertexAttribArray(index uint32) {
	purego.SyscallN(f.gpEnableVertexAttribArray, uintptr(index))
}

// Flush()
func (f *Functions) Flush() {
	purego.SyscallN(f.gpFlush)
}

// FramebufferRenderbuffer(target, attachment, renderbuffertarget Enum, renderbuffer Renderbuffer)
func (f *Functions) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	purego.SyscallN(f.gpFramebufferRenderbuffer, uintptr(target), uintptr(attachment), uintptr(renderbuffertarget), uintptr(renderbuffer))
}

// FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int)
func (f *Functions) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	purego.SyscallN(f.gpFramebufferTexture2D, uintptr(target), uintptr(attachment), uintptr(textarget), uintptr(texture), uintptr(level))
}

// GetError()
func (f *Functions) GetError() uint32 {
	ret, _, _ := purego.SyscallN(f.gpGetError)
	return uint32(ret)
}

func (f *Functions) GetExtension(name string) any {
	return nil
}

// GetInteger(pname Enum) int
func (f *Functions) GetInteger(pname uint32) int {
	var dst int32
	purego.SyscallN(f.gpGetIntegerv, uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetProgramInfoLog(p Program) string
func (f *Functions) GetProgramInfoLog(program uint32) string {
	bufSize := f.GetProgrami(program, INFO_LOG_LENGTH)
	if bufSize == 0 {
		return ""
	}
	infoLog := make([]byte, bufSize)
	purego.SyscallN(f.gpGetProgramInfoLog, uintptr(program), uintptr(bufSize), 0, uintptr(unsafe.Pointer(&infoLog[0])))
	return string(infoLog)
}

// GetProgrami(p Program, pname Enum) int
func (f *Functions) GetProgrami(program uint32, pname uint32) int {
	var dst int32
	purego.SyscallN(f.gpGetProgramiv, uintptr(program), uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetShaderInfoLog(s Shader) string
func (f *Functions) GetShaderInfoLog(shader uint32) string {
	bufSize := f.GetShaderi(shader, INFO_LOG_LENGTH)
	if bufSize == 0 {
		return ""
	}
	infoLog := make([]byte, bufSize)
	purego.SyscallN(f.gpGetShaderInfoLog, uintptr(shader), uintptr(bufSize), 0, uintptr(unsafe.Pointer(&infoLog[0])))
	return string(infoLog)
}

// GetShaderi(s Shader, pname Enum) int
func (f *Functions) GetShaderi(shader uint32, pname uint32) int {
	var dst int32
	purego.SyscallN(f.gpGetShaderiv, uintptr(shader), uintptr(pname), uintptr(unsafe.Pointer(&dst)))
	return int(dst)
}

// GetUniformLocation(p Program, name string) Uniform
func (f *Functions) GetUniformLocation(program uint32, name string) int32 {
	cname, free := cStr(name)
	defer free()
	ret, _, _ := purego.SyscallN(f.gpGetUniformLocation, uintptr(program), uintptr(unsafe.Pointer(cname)))
	return int32(ret)
}

func (f *Functions) IsProgram(program uint32) bool {
	ret, _, _ := purego.SyscallN(f.gpIsProgram, uintptr(program))
	return byte(ret) != 0
}

// LinkProgram(p Program)
func (f *Functions) LinkProgram(program uint32) {
	purego.SyscallN(f.gpLinkProgram, uintptr(program))
}

// PixelStorei(pname Enum, param int)
func (f *Functions) PixelStorei(pname uint32, param int32) {
	purego.SyscallN(f.gpPixelStorei, uintptr(pname), uintptr(param))
}

// ReadPixels(x, y, width, height int, format, ty Enum, data []byte)
func (f *Functions) ReadPixels(dst []byte, x int32, y int32, width int32, height int32, format uint32, xtype uint32) {
	purego.SyscallN(f.gpReadPixels, uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(&dst[0])))
}

// RenderbufferStorage(target, internalformat Enum, width, height int)
func (f *Functions) RenderbufferStorage(target uint32, internalformat uint32, width int32, height int32) {
	purego.SyscallN(f.gpRenderbufferStorage, uintptr(target), uintptr(internalformat), uintptr(width), uintptr(height))
}

// Scissor(x, y, width, height int32)
func (f *Functions) Scissor(x int32, y int32, width int32, height int32) {
	purego.SyscallN(f.gpScissor, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}

// ShaderSource(s Shader, src string)
func (f *Functions) ShaderSource(shader uint32, xstring string) {
	cstring, free := cStr(xstring)
	defer free()
	purego.SyscallN(f.gpShaderSource, uintptr(shader), 1, uintptr(unsafe.Pointer(&cstring)), 0)
}

func (f *Functions) StencilFunc(xfunc uint32, ref int32, mask uint32) {
	purego.SyscallN(f.gpStencilFunc, uintptr(xfunc), uintptr(ref), uintptr(mask))
}

func (f *Functions) StencilOpSeparate(face uint32, fail uint32, zfail uint32, zpass uint32) {
	purego.SyscallN(f.gpStencilOpSeparate, uintptr(face), uintptr(fail), uintptr(zfail), uintptr(zpass))
}

// TexImage2D(target Enum, level int, width int, height int, format Enum, ty Enum, data []byte)
func (f *Functions) TexImage2D(target uint32, level int32, width int32, height int32, format uint32, xtype uint32, pixels []byte) {
	var ptr *byte
	if len(pixels) > 0 {
		ptr = &pixels[0]
	}
	purego.SyscallN(f.gpTexImage2D, uintptr(target), uintptr(level), uintptr(format), uintptr(width), uintptr(height), 0, uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(ptr)))
	runtime.KeepAlive(pixels)
}

// TexParameteri(target, pname Enum, param int)
func (f *Functions) TexParameteri(target uint32, pname uint32, param int32) {
	purego.SyscallN(f.gpTexParameteri, uintptr(target), uintptr(pname), uintptr(param))
}

// TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte)
func (f *Functions) TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels []byte) {
	purego.SyscallN(f.gpTexSubImage2D, uintptr(target), uintptr(level), uintptr(xoffset), uintptr(yoffset), uintptr(width), uintptr(height), uintptr(format), uintptr(xtype), uintptr(unsafe.Pointer(&pixels[0])))
	runtime.KeepAlive(pixels)
}

// Uniform1fv(dst Uniform, src []float32)
func (f *Functions) Uniform1fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniform1fv, uintptr(location), uintptr(len(value)), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform1i(dst Uniform, v int)
func (f *Functions) Uniform1i(location int32, v0 int32) {
	purego.SyscallN(f.gpUniform1i, uintptr(location), uintptr(v0))
}

func (f *Functions) Uniform1iv(location int32, value []int32) {
	purego.SyscallN(f.gpUniform1iv, uintptr(location), uintptr(len(value)), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform2fv(dst Uniform, src []float32)
func (f *Functions) Uniform2fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniform2fv, uintptr(location), uintptr(len(value)/2), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (f *Functions) Uniform2iv(location int32, value []int32) {
	purego.SyscallN(f.gpUniform2iv, uintptr(location), uintptr(len(value)/2), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform3fv(dst Uniform, src []float32)
func (f *Functions) Uniform3fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniform3fv, uintptr(location), uintptr(len(value)/3), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (f *Functions) Uniform3iv(location int32, value []int32) {
	purego.SyscallN(f.gpUniform3iv, uintptr(location), uintptr(len(value)/3), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// Uniform4fv(dst Uniform, src []float32)
func (f *Functions) Uniform4fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniform4fv, uintptr(location), uintptr(len(value)/4), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

func (f *Functions) Uniform4iv(location int32, value []int32) {
	purego.SyscallN(f.gpUniform4iv, uintptr(location), uintptr(len(value)/4), uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix2fv(dst Uniform, src []float32)
func (f *Functions) UniformMatrix2fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniformMatrix2fv, uintptr(location), uintptr(len(value)/4), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix3fv(dst Uniform, src []float32)
func (f *Functions) UniformMatrix3fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniformMatrix3fv, uintptr(location), uintptr(len(value)/9), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UniformMatrix4fv(dst Uniform, src []float32)
func (f *Functions) UniformMatrix4fv(location int32, value []float32) {
	purego.SyscallN(f.gpUniformMatrix4fv, uintptr(location), uintptr(len(value)/16), 0, uintptr(unsafe.Pointer(&value[0])))
	runtime.KeepAlive(value)
}

// UseProgram(p Program)
func (f *Functions) UseProgram(program uint32) {
	purego.SyscallN(f.gpUseProgram, uintptr(program))
}

// VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int)
// Do I have to call glVertexAttribPointer() for each VAO that I create?
// Yes : the VAO contains the state of the attrib pointers and a bunch of other stuff. This allows you to pull from two different buffer objects within one VAO.
func (f *Functions) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	purego.SyscallN(f.gpVertexAttribPointer, uintptr(index), uintptr(size), uintptr(xtype), uintptr(boolToInt(normalized)), uintptr(stride), uintptr(offset))
}

// Viewport(x int, y int, width int, height int)
func (f *Functions) Viewport(x int32, y int32, width int32, height int32) {
	purego.SyscallN(f.gpViewport, uintptr(x), uintptr(y), uintptr(width), uintptr(height))
}

func (f *Functions) LoadFunctions() error {
	var err error
	f.gpActiveTexture, err = f.get("glActiveTexture")
	if err != nil {
		return err
	}
	f.gpAttachShader, err = f.get("glAttachShader")
	if err != nil {
		return err
	}
	f.gpBindAttribLocation, err = f.get("glBindAttribLocation")
	if err != nil {
		return err
	}
	f.gpBindBuffer, err = f.get("glBindBuffer")
	if err != nil {
		return err
	}
	f.gpBindFramebuffer, err = f.get("glBindFramebuffer")
	if err != nil {
		return err
	}
	f.gpBindRenderbuffer, err = f.get("glBindRenderbuffer")
	if err != nil {
		return err
	}
	f.gpBindTexture, err = f.get("glBindTexture")
	if err != nil {
		return err
	}
	f.gpBindVertexArray, err = f.get("glBindVertexArray")
	if err != nil {
		return err
	}
	f.gpBlendEquationSeparate, err = f.get("glBlendEquationSeparate")
	if err != nil {
		return err
	}
	f.gpBlendFuncSeparate, err = f.get("glBlendFuncSeparate")
	if err != nil {
		return err
	}
	f.gpBufferData, err = f.get("glBufferData")
	if err != nil {
		return err
	}
	f.gpBufferSubData, err = f.get("glBufferSubData")
	if err != nil {
		return err
	}
	f.gpCheckFramebufferStatus, err = f.get("glCheckFramebufferStatus")
	if err != nil {
		return err
	}
	f.gpClear, err = f.get("glClear")
	if err != nil {
		return err
	}
	f.gpColorMask, err = f.get("glColorMask")
	if err != nil {
		return err
	}
	f.gpCompileShader, err = f.get("glCompileShader")
	if err != nil {
		return err
	}
	f.gpCreateProgram, err = f.get("glCreateProgram")
	if err != nil {
		return err
	}
	f.gpCreateShader, err = f.get("glCreateShader")
	if err != nil {
		return err
	}
	f.gpDeleteBuffers, err = f.get("glDeleteBuffers")
	if err != nil {
		return err
	}
	f.gpDeleteFramebuffers, err = f.get("glDeleteFramebuffers")
	if err != nil {
		return err
	}
	f.gpDeleteProgram, err = f.get("glDeleteProgram")
	if err != nil {
		return err
	}
	f.gpDeleteRenderbuffers, err = f.get("glDeleteRenderbuffers")
	if err != nil {
		return err
	}
	f.gpDeleteShader, err = f.get("glDeleteShader")
	if err != nil {
		return err
	}
	f.gpDeleteTextures, err = f.get("glDeleteTextures")
	if err != nil {
		return err
	}
	f.gpDeleteVertexArrays, err = f.get("glDeleteVertexArrays")
	if err != nil {
		return err
	}
	f.gpDisable, err = f.get("glDisable")
	if err != nil {
		return err
	}
	f.gpDisableVertexAttribArray, err = f.get("glDisableVertexAttribArray")
	if err != nil {
		return err
	}
	f.gpDrawElements, err = f.get("glDrawElements")
	if err != nil {
		return err
	}
	f.gpEnable, err = f.get("glEnable")
	if err != nil {
		return err
	}
	f.gpEnableVertexAttribArray, err = f.get("glEnableVertexAttribArray")
	if err != nil {
		return err
	}
	f.gpFlush, err = f.get("glFlush")
	if err != nil {
		return err
	}
	f.gpFramebufferRenderbuffer, err = f.get("glFramebufferRenderbuffer")
	if err != nil {
		return err
	}
	f.gpFramebufferTexture2D, err = f.get("glFramebufferTexture2D")
	if err != nil {
		return err
	}
	f.gpGenBuffers, err = f.get("glGenBuffers")
	if err != nil {
		return err
	}
	f.gpGenFramebuffers, err = f.get("glGenFramebuffers")
	if err != nil {
		return err
	}
	f.gpGenRenderbuffers, err = f.get("glGenRenderbuffers")
	if err != nil {
		return err
	}
	f.gpGenTextures, err = f.get("glGenTextures")
	if err != nil {
		return err
	}
	f.gpGenVertexArrays, err = f.get("glGenVertexArrays")
	if err != nil {
		return err
	}
	f.gpGetError, err = f.get("glGetError")
	if err != nil {
		return err
	}
	f.gpGetIntegerv, err = f.get("glGetIntegerv")
	if err != nil {
		return err
	}
	f.gpGetProgramInfoLog, err = f.get("glGetProgramInfoLog")
	if err != nil {
		return err
	}
	f.gpGetProgramiv, err = f.get("glGetProgramiv")
	if err != nil {
		return err
	}
	f.gpGetShaderInfoLog, err = f.get("glGetShaderInfoLog")
	if err != nil {
		return err
	}
	f.gpGetShaderiv, err = f.get("glGetShaderiv")
	if err != nil {
		return err
	}
	f.gpGetUniformLocation, err = f.get("glGetUniformLocation")
	if err != nil {
		return err
	}
	f.gpIsProgram, err = f.get("glIsProgram")
	if err != nil {
		return err
	}
	f.gpLinkProgram, err = f.get("glLinkProgram")
	if err != nil {
		return err
	}
	f.gpPixelStorei, err = f.get("glPixelStorei")
	if err != nil {
		return err
	}
	f.gpReadPixels, err = f.get("glReadPixels")
	if err != nil {
		return err
	}
	f.gpRenderbufferStorage, err = f.get("glRenderbufferStorage")
	if err != nil {
		return err
	}
	f.gpScissor, err = f.get("glScissor")
	if err != nil {
		return err
	}
	f.gpShaderSource, err = f.get("glShaderSource")
	if err != nil {
		return err
	}
	f.gpStencilFunc, err = f.get("glStencilFunc")
	if err != nil {
		return err
	}
	f.gpStencilOpSeparate, err = f.get("glStencilOpSeparate")
	if err != nil {
		return err
	}
	f.gpTexImage2D, err = f.get("glTexImage2D")
	if err != nil {
		return err
	}
	f.gpTexParameteri, err = f.get("glTexParameteri")
	if err != nil {
		return err
	}
	f.gpTexSubImage2D, err = f.get("glTexSubImage2D")
	if err != nil {
		return err
	}
	f.gpUniform1fv, err = f.get("glUniform1fv")
	if err != nil {
		return err
	}
	f.gpUniform1i, err = f.get("glUniform1i")
	if err != nil {
		return err
	}
	f.gpUniform1iv, err = f.get("glUniform1iv")
	if err != nil {
		return err
	}
	f.gpUniform2fv, err = f.get("glUniform2fv")
	if err != nil {
		return err
	}
	f.gpUniform2iv, err = f.get("glUniform2iv")
	if err != nil {
		return err
	}
	f.gpUniform3fv, err = f.get("glUniform3fv")
	if err != nil {
		return err
	}
	f.gpUniform3iv, err = f.get("glUniform3iv")
	if err != nil {
		return err
	}
	f.gpUniform4fv, err = f.get("glUniform4fv")
	if err != nil {
		return err
	}
	f.gpUniform4iv, err = f.get("glUniform4iv")
	if err != nil {
		return err
	}
	f.gpUniformMatrix2fv, err = f.get("glUniformMatrix2fv")
	if err != nil {
		return err
	}
	f.gpUniformMatrix3fv, err = f.get("glUniformMatrix3fv")
	if err != nil {
		return err
	}
	f.gpUniformMatrix4fv, err = f.get("glUniformMatrix4fv")
	if err != nil {
		return err
	}
	f.gpUseProgram, err = f.get("glUseProgram")
	if err != nil {
		return err
	}
	f.gpVertexAttribPointer, err = f.get("glVertexAttribPointer")
	if err != nil {
		return err
	}
	f.gpViewport, err = f.get("glViewport")
	if err != nil {
		return err
	}

	f.glBufferData, err = f.get("glBufferData")
	if err != nil {
		return err
	}

	f.glClearColor, err = f.get("glClearColor")
	if err != nil {
		return err
	}
	f.glDrawArrays, err = f.get("glDrawArrays")
	if err != nil {
		return err
	}
	f.glUniform1f, err = f.get("glUniform1f")
	if err != nil {
		return err
	}
	f.glBlendFunc, err = f.get("glBlendFunc")
	if err != nil {
		return err
	}
	f.glGetActiveUniform, err = f.get("glGetActiveUniform")
	if err != nil {
		return err
	}
	f.glGetActiveAttrib, err = f.get("glGetActiveAttrib")
	if err != nil {
		return err
	}
	f.glGetAttribLocation, err = f.get("glGetAttribLocation")
	if err != nil {
		return err
	}

	return nil
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
