//go:build darwin || linux || freebsd || openbsd
// +build darwin linux freebsd openbsd

package purego

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

/*
#cgo CFLAGS: -Werror
#cgo linux freebsd LDFLAGS: -ldl

#include <stdint.h>
#include <stdlib.h>
#include <sys/types.h>
#define __USE_GNU
#include <dlfcn.h>

typedef unsigned int GLenum;
typedef unsigned int GLuint;
typedef char GLchar;
typedef float GLfloat;
typedef ssize_t GLsizeiptr;
typedef intptr_t GLintptr;
typedef unsigned int GLbitfield;
typedef int GLint;
typedef unsigned char GLboolean;
typedef int GLsizei;
typedef uint8_t GLubyte;

typedef void (*_glActiveTexture)(GLenum texture);
typedef void (*_glAttachShader)(GLuint program, GLuint shader);
typedef void (*_glBindAttribLocation)(GLuint program, GLuint index, const GLchar *name);
typedef void (*_glBindBuffer)(GLenum target, GLuint buffer);
typedef void (*_glBindFramebuffer)(GLenum target, GLuint framebuffer);
typedef void (*_glBindRenderbuffer)(GLenum target, GLuint renderbuffer);
typedef void (*_glBindTexture)(GLenum target, GLuint texture);
typedef void (*_glBlendEquation)(GLenum mode);
typedef void (*_glBlendFuncSeparate)(GLenum srcRGB, GLenum dstRGB, GLenum srcA, GLenum dstA);
typedef void (*_glBufferData)(GLenum target, GLsizeiptr size, const void *data, GLenum usage);
typedef void (*_glBufferSubData)(GLenum target, GLintptr offset, GLsizeiptr size, const void *data);
typedef GLenum (*_glCheckFramebufferStatus)(GLenum target);
typedef void (*_glClear)(GLbitfield mask);
typedef void (*_glClearColor)(GLfloat red, GLfloat green, GLfloat blue, GLfloat alpha);
typedef void (*_glClearDepthf)(GLfloat d);
typedef void (*_glCompileShader)(GLuint shader);
typedef void (*_glCopyTexSubImage2D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint x, GLint y, GLsizei width, GLsizei height);
typedef GLuint (*_glCreateProgram)(void);
typedef GLuint (*_glCreateShader)(GLenum type);
typedef void (*_glDeleteBuffers)(GLsizei n, const GLuint *buffers);
typedef void (*_glDeleteFramebuffers)(GLsizei n, const GLuint *framebuffers);
typedef void (*_glDeleteProgram)(GLuint program);
typedef void (*_glDeleteRenderbuffers)(GLsizei n, const GLuint *renderbuffers);
typedef void (*_glDeleteShader)(GLuint shader);
typedef void (*_glDeleteTextures)(GLsizei n, const GLuint *textures);
typedef void (*_glDepthFunc)(GLenum func);
typedef void (*_glDepthMask)(GLboolean flag);
typedef void (*_glDisable)(GLenum cap);
typedef void (*_glDisableVertexAttribArray)(GLuint index);
typedef void (*_glDrawArrays)(GLenum mode, GLint first, GLsizei count);
typedef void (*_glDrawElements)(GLenum mode, GLsizei count, GLenum type, const void *indices);
typedef void (*_glEnable)(GLenum cap);
typedef void (*_glEnableVertexAttribArray)(GLuint index);
typedef void (*_glFinish)(void);
typedef void (*_glFlush)(void);
typedef void (*_glFramebufferRenderbuffer)(GLenum target, GLenum attachment, GLenum renderbuffertarget, GLuint renderbuffer);
typedef void (*_glFramebufferTexture2D)(GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level);
typedef void (*_glGenBuffers)(GLsizei n, GLuint *buffers);
typedef void (*_glGenerateMipmap)(GLenum target);
typedef void (*_glGenFramebuffers)(GLsizei n, GLuint *framebuffers);
typedef void (*_glGenRenderbuffers)(GLsizei n, GLuint *renderbuffers);
typedef void (*_glGenTextures)(GLsizei n, GLuint *textures);
typedef GLenum (*_glGetError)(void);
typedef void (*_glGetFramebufferAttachmentParameteriv)(GLenum target, GLenum attachment, GLenum pname, GLint *params);
typedef void (*_glGetFloatv)(GLenum pname, GLfloat *data);
typedef void (*_glGetIntegerv)(GLenum pname, GLint *data);
typedef void (*_glGetIntegeri_v)(GLenum pname, GLuint idx, GLint *data);
typedef void (*_glGetProgramiv)(GLuint program, GLenum pname, GLint *params);
typedef void (*_glGetProgramInfoLog)(GLuint program, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
typedef void (*_glGetRenderbufferParameteriv)(GLenum target, GLenum pname, GLint *params);
typedef void (*_glGetShaderiv)(GLuint shader, GLenum pname, GLint *params);
typedef void (*_glGetShaderInfoLog)(GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog);
typedef const GLubyte *(*_glGetString)(GLenum name);
typedef GLint (*_glGetUniformLocation)(GLuint program, const GLchar *name);
typedef void (*_glGetVertexAttribiv)(GLuint index, GLenum pname, GLint *params);
typedef void (*_glGetVertexAttribPointerv)(GLuint index, GLenum pname, void **params);
typedef GLboolean (*_glIsEnabled)(GLenum cap);
typedef void (*_glLinkProgram)(GLuint program);
typedef void (*_glPixelStorei)(GLenum pname, GLint param);
typedef void (*_glReadPixels)(GLint x, GLint y, GLsizei width, GLsizei height, GLenum format, GLenum type, void *pixels);
typedef void (*_glRenderbufferStorage)(GLenum target, GLenum internalformat, GLsizei width, GLsizei height);
typedef void (*_glScissor)(GLint x, GLint y, GLsizei width, GLsizei height);
typedef void (*_glShaderSource)(GLuint shader, GLsizei count, const GLchar *const*string, const GLint *length);
typedef void (*_glTexImage2D)(GLenum target, GLint level, GLint internalformat, GLsizei width, GLsizei height, GLint border, GLenum format, GLenum type, const void *pixels);
typedef void (*_glTexParameteri)(GLenum target, GLenum pname, GLint param);
typedef void (*_glTexSubImage2D)(GLenum target, GLint level, GLint xoffset, GLint yoffset, GLsizei width, GLsizei height, GLenum format, GLenum type, const void *pixels);
typedef void (*_glUniform1f)(GLint location, GLfloat v0);
typedef void (*_glUniform1i)(GLint location, GLint v0);
typedef void (*_glUniform2f)(GLint location, GLfloat v0, GLfloat v1);
typedef void (*_glUniform3f)(GLint location, GLfloat v0, GLfloat v1, GLfloat v2);
typedef void (*_glUniform4f)(GLint location, GLfloat v0, GLfloat v1, GLfloat v2, GLfloat v3);
typedef void (*_glUseProgram)(GLuint program);
typedef void (*_glVertexAttribPointer)(GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, const void *pointer);
typedef void (*_glViewport)(GLint x, GLint y, GLsizei width, GLsizei height);
typedef void (*_glBindVertexArray)(GLuint array);
typedef void (*_glBindBufferBase)(GLenum target, GLuint index, GLuint buffer);
typedef GLuint (*_glGetUniformBlockIndex)(GLuint program, const GLchar *uniformBlockName);
typedef void (*_glUniformBlockBinding)(GLuint program, GLuint uniformBlockIndex, GLuint uniformBlockBinding);
typedef void (*_glInvalidateFramebuffer)(GLenum target, GLsizei numAttachments, const GLenum *attachments);
typedef void (*_glBeginQuery)(GLenum target, GLuint id);
typedef void (*_glDeleteQueries)(GLsizei n, const GLuint *ids);
typedef void (*_glDeleteVertexArrays)(GLsizei n, const GLuint *ids);
typedef void (*_glEndQuery)(GLenum target);
typedef void (*_glGenQueries)(GLsizei n, GLuint *ids);
typedef void (*_glGenVertexArrays)(GLsizei n, GLuint *ids);
typedef void (*_glGetProgramBinary)(GLuint program, GLsizei bufsize, GLsizei *length, GLenum *binaryFormat, void *binary);
typedef void (*_glGetQueryObjectuiv)(GLuint id, GLenum pname, GLuint *params);
typedef const GLubyte* (*_glGetStringi)(GLenum name, GLuint index);
typedef void (*_glDispatchCompute)(GLuint x, GLuint y, GLuint z);
typedef void (*_glMemoryBarrier)(GLbitfield barriers);
typedef void* (*_glMapBufferRange)(GLenum target, GLintptr offset, GLsizeiptr length, GLbitfield access);
typedef GLboolean (*_glUnmapBuffer)(GLenum target);
typedef void (*_glBindImageTexture)(GLuint unit, GLuint texture, GLint level, GLboolean layered, GLint layer, GLenum access, GLenum format);
typedef void (*_glTexStorage2D)(GLenum target, GLsizei levels, GLenum internalformat, GLsizei width, GLsizei height);
typedef void (*_glBlitFramebuffer)(GLint srcX0, GLint srcY0, GLint srcX1, GLint srcY1, GLint dstX0, GLint dstY0, GLint dstX1, GLint dstY1, GLbitfield mask, GLenum filter);

typedef void  (*_glGetActiveUniform)(GLuint  program, GLuint  index, GLsizei  bufSize, GLsizei * length, GLint * size, GLenum * type, GLchar * name);
typedef void  (*_glGetActiveAttrib)(GLuint  program, GLuint  index, GLsizei  bufSize, GLsizei * length, GLint * size, GLenum * type, GLchar * name);
typedef GLint  (*_glGetAttribLocation)(GLuint  program, const GLchar * name);
typedef void  (*_glBlendFunc)(GLenum  sfactor, GLenum  dfactor);
typedef void  (*_glUniformMatrix2fv)(GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value);
typedef void  (*_glUniformMatrix3fv)(GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value);
typedef void  (*_glUniformMatrix4fv)(GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value);
typedef void  (*_glUniform1fv)(GLint  location, GLsizei  count, const GLfloat * value);
typedef void  (*_glUniform2fv)(GLint  location, GLsizei  count, const GLfloat * value);
typedef void  (*_glUniform3fv)(GLint  location, GLsizei  count, const GLfloat * value);
typedef void  (*_glUniform4fv)(GLint  location, GLsizei  count, const GLfloat * value);

static void glActiveTexture(_glActiveTexture f, GLenum texture) {
	f(texture);
}

static void glAttachShader(_glAttachShader f, GLuint program, GLuint shader) {
	f(program, shader);
}

static void glBindAttribLocation(_glBindAttribLocation f, GLuint program, GLuint index, const GLchar *name) {
	f(program, index, name);
}

static void glBindBuffer(_glBindBuffer f, GLenum target, GLuint buffer) {
	f(target, buffer);
}

static void glBindFramebuffer(_glBindFramebuffer f, GLenum target, GLuint framebuffer) {
	f(target, framebuffer);
}

static void glBindRenderbuffer(_glBindRenderbuffer f, GLenum target, GLuint renderbuffer) {
	f(target, renderbuffer);
}

static void glBindTexture(_glBindTexture f, GLenum target, GLuint texture) {
	f(target, texture);
}

static void glBindVertexArray(_glBindVertexArray f, GLuint array) {
	f(array);
}

static void glBlendEquation(_glBlendEquation f, GLenum mode) {
	f(mode);
}

static void glBlendFuncSeparate(_glBlendFuncSeparate f, GLenum srcRGB, GLenum dstRGB, GLenum srcA, GLenum dstA) {
	f(srcRGB, dstRGB, srcA, dstA);
}

static void glBufferData(_glBufferData f, GLenum target, GLsizeiptr size, const void *data, GLenum usage) {
	f(target, size, data, usage);
}

static void glBufferSubData(_glBufferSubData f, GLenum target, GLintptr offset, GLsizeiptr size, const void *data) {
	f(target, offset, size, data);
}

static GLenum glCheckFramebufferStatus(_glCheckFramebufferStatus f, GLenum target) {
	return f(target);
}

static void glClear(_glClear f, GLbitfield mask) {
	f(mask);
}

static void glClearColor(_glClearColor f, GLfloat red, GLfloat green, GLfloat blue, GLfloat alpha) {
	f(red, green, blue, alpha);
}

static void glClearDepthf(_glClearDepthf f, GLfloat d) {
	f(d);
}

static void glCompileShader(_glCompileShader f, GLuint shader) {
	f(shader);
}

static void glCopyTexSubImage2D(_glCopyTexSubImage2D f, GLenum target, GLint level, GLint xoffset, GLint yoffset, GLint x, GLint y, GLsizei width, GLsizei height) {
	f(target, level, xoffset, yoffset, x, y, width, height);
}

static GLuint glCreateProgram(_glCreateProgram f) {
	return f();
}

static GLuint glCreateShader(_glCreateShader f, GLenum type) {
	return f(type);
}

static void glDeleteBuffers(_glDeleteBuffers f, GLsizei n, const GLuint *buffers) {
	f(n, buffers);
}

static void glDeleteFramebuffers(_glDeleteFramebuffers f, GLsizei n, const GLuint *framebuffers) {
	f(n, framebuffers);
}

static void glDeleteProgram(_glDeleteProgram f, GLuint program) {
	f(program);
}

static void glDeleteRenderbuffers(_glDeleteRenderbuffers f, GLsizei n, const GLuint *renderbuffers) {
	f(n, renderbuffers);
}

static void glDeleteShader(_glDeleteShader f, GLuint shader) {
	f(shader);
}

static void glDeleteTextures(_glDeleteTextures f, GLsizei n, const GLuint *textures) {
	f(n, textures);
}

static void glDepthFunc(_glDepthFunc f, GLenum func) {
	f(func);
}

static void glDepthMask(_glDepthMask f, GLboolean flag) {
	f(flag);
}

static void glDisable(_glDisable f, GLenum cap) {
	f(cap);
}

static void glDisableVertexAttribArray(_glDisableVertexAttribArray f, GLuint index) {
	f(index);
}

static void glDrawArrays(_glDrawArrays f, GLenum mode, GLint first, GLsizei count) {
	f(mode, first, count);
}

// offset is defined as an uintptr_t to omit Cgo pointer checks.
static void glDrawElements(_glDrawElements f, GLenum mode, GLsizei count, GLenum type, const uintptr_t offset) {
	f(mode, count, type, (const void *)offset);
}

static void glEnable(_glEnable f, GLenum cap) {
	f(cap);
}

static void glEnableVertexAttribArray(_glEnableVertexAttribArray f, GLuint index) {
	f(index);
}

static void glFinish(_glFinish f) {
	f();
}

static void glFlush(_glFlush f) {
	f();
}

static void glFramebufferRenderbuffer(_glFramebufferRenderbuffer f, GLenum target, GLenum attachment, GLenum renderbuffertarget, GLuint renderbuffer) {
	f(target, attachment, renderbuffertarget, renderbuffer);
}

static void glFramebufferTexture2D(_glFramebufferTexture2D f, GLenum target, GLenum attachment, GLenum textarget, GLuint texture, GLint level) {
	f(target, attachment, textarget, texture, level);
}

static void glGenBuffers(_glGenBuffers f, GLsizei n, GLuint *buffers) {
	f(n, buffers);
}

static void glGenerateMipmap(_glGenerateMipmap f, GLenum target) {
	f(target);
}

static void glGenFramebuffers(_glGenFramebuffers f, GLsizei n, GLuint *framebuffers) {
	f(n, framebuffers);
}

static void glGenRenderbuffers(_glGenRenderbuffers f, GLsizei n, GLuint *renderbuffers) {
	f(n, renderbuffers);
}

static void glGenTextures(_glGenTextures f, GLsizei n, GLuint *textures) {
	f(n, textures);
}

static GLenum glGetError(_glGetError f) {
	return f();
}

static void glGetFramebufferAttachmentParameteriv(_glGetFramebufferAttachmentParameteriv f, GLenum target, GLenum attachment, GLenum pname, GLint *params) {
	f(target, attachment, pname, params);
}

static void glGetIntegerv(_glGetIntegerv f, GLenum pname, GLint *data) {
	f(pname, data);
}

static void glGetFloatv(_glGetFloatv f, GLenum pname, GLfloat *data) {
	f(pname, data);
}

static void glGetIntegeri_v(_glGetIntegeri_v f, GLenum pname, GLuint idx, GLint *data) {
	f(pname, idx, data);
}

static void glGetProgramiv(_glGetProgramiv f, GLuint program, GLenum pname, GLint *params) {
	f(program, pname, params);
}

static void glGetProgramInfoLog(_glGetProgramInfoLog f, GLuint program, GLsizei bufSize, GLsizei *length, GLchar *infoLog) {
	f(program, bufSize, length, infoLog);
}

static void glGetRenderbufferParameteriv(_glGetRenderbufferParameteriv f, GLenum target, GLenum pname, GLint *params) {
	f(target, pname, params);
}

static void glGetShaderiv(_glGetShaderiv f, GLuint shader, GLenum pname, GLint *params) {
	f(shader, pname, params);
}

static void glGetShaderInfoLog(_glGetShaderInfoLog f, GLuint shader, GLsizei bufSize, GLsizei *length, GLchar *infoLog) {
	f(shader, bufSize, length, infoLog);
}

static const GLubyte *glGetString(_glGetString f, GLenum name) {
	return f(name);
}

static GLint glGetUniformLocation(_glGetUniformLocation f, GLuint program, const GLchar *name) {
	return f(program, name);
}

static void glGetVertexAttribiv(_glGetVertexAttribiv f, GLuint index, GLenum pname, GLint *data) {
	f(index, pname, data);
}

// Return uintptr_t to avoid Cgo pointer check.
static uintptr_t glGetVertexAttribPointerv(_glGetVertexAttribPointerv f, GLuint index, GLenum pname) {
	void *ptrs;
	f(index, pname, &ptrs);
	return (uintptr_t)ptrs;
}

static GLboolean glIsEnabled(_glIsEnabled f, GLenum cap) {
	return f(cap);
}

static void glLinkProgram(_glLinkProgram f, GLuint program) {
	f(program);
}

static void glPixelStorei(_glPixelStorei f, GLenum pname, GLint param) {
	f(pname, param);
}

static void glReadPixels(_glReadPixels f, GLint x, GLint y, GLsizei width, GLsizei height, GLenum format, GLenum type, void *pixels) {
	f(x, y, width, height, format, type, pixels);
}

static void glRenderbufferStorage(_glRenderbufferStorage f, GLenum target, GLenum internalformat, GLsizei width, GLsizei height) {
	f(target, internalformat, width, height);
}

static void glScissor(_glScissor f, GLint x, GLint y, GLsizei width, GLsizei height) {
	f(x, y, width, height);
}

static void glShaderSource(_glShaderSource f, GLuint shader, GLsizei count, const GLchar *const*string, const GLint *length) {
	f(shader, count, string, length);
}

static void glTexImage2D(_glTexImage2D f, GLenum target, GLint level, GLint internalformat, GLsizei width, GLsizei height, GLint border, GLenum format, GLenum type, const void *pixels) {
	f(target, level, internalformat, width, height, border, format, type, pixels);
}

static void glTexParameteri(_glTexParameteri f, GLenum target, GLenum pname, GLint param) {
	f(target, pname, param);
}

static void glTexSubImage2D(_glTexSubImage2D f, GLenum target, GLint level, GLint xoffset, GLint yoffset, GLsizei width, GLsizei height, GLenum format, GLenum type, const void *pixels) {
	f(target, level, xoffset, yoffset, width, height, format, type, pixels);
}

static void glUniform1f(_glUniform1f f, GLint location, GLfloat v0) {
	f(location, v0);
}

static void glUniform1i(_glUniform1i f, GLint location, GLint v0) {
	f(location, v0);
}

static void glUniform2f(_glUniform2f f, GLint location, GLfloat v0, GLfloat v1) {
	f(location, v0, v1);
}

static void glUniform3f(_glUniform3f f, GLint location, GLfloat v0, GLfloat v1, GLfloat v2) {
	f(location, v0, v1, v2);
}

static void glUniform4f(_glUniform4f f, GLint location, GLfloat v0, GLfloat v1, GLfloat v2, GLfloat v3) {
	f(location, v0, v1, v2, v3);
}

static void glUseProgram(_glUseProgram f, GLuint program) {
	f(program);
}

// offset is defined as an uintptr_t to omit Cgo pointer checks.
static void glVertexAttribPointer(_glVertexAttribPointer f, GLuint index, GLint size, GLenum type, GLboolean normalized, GLsizei stride, uintptr_t offset) {
	f(index, size, type, normalized, stride, (const void *)offset);
}

static void glViewport(_glViewport f, GLint x, GLint y, GLsizei width, GLsizei height) {
	f(x, y, width, height);
}

static void glBindBufferBase(_glBindBufferBase f, GLenum target, GLuint index, GLuint buffer) {
	f(target, index, buffer);
}

static void glUniformBlockBinding(_glUniformBlockBinding f, GLuint program, GLuint uniformBlockIndex, GLuint uniformBlockBinding) {
	f(program, uniformBlockIndex, uniformBlockBinding);
}

static GLuint glGetUniformBlockIndex(_glGetUniformBlockIndex f, GLuint program, const GLchar *uniformBlockName) {
	return f(program, uniformBlockName);
}

static void glInvalidateFramebuffer(_glInvalidateFramebuffer f, GLenum target, GLenum attachment) {
	// Framebuffer invalidation is just a hint and can safely be ignored.
	if (f != NULL) {
		f(target, 1, &attachment);
	}
}

static void glBeginQuery(_glBeginQuery f, GLenum target, GLenum attachment) {
	f(target, attachment);
}

static void glDeleteQueries(_glDeleteQueries f, GLsizei n, const GLuint *ids) {
	f(n, ids);
}

static void glDeleteVertexArrays(_glDeleteVertexArrays f, GLsizei n, const GLuint *ids) {
	f(n, ids);
}

static void glEndQuery(_glEndQuery f, GLenum target) {
	f(target);
}

static const GLubyte* glGetStringi(_glGetStringi f, GLenum name, GLuint index) {
	return f(name, index);
}

static void glGenQueries(_glGenQueries f, GLsizei n, GLuint *ids) {
	f(n, ids);
}

static void glGenVertexArrays(_glGenVertexArrays f, GLsizei n, GLuint *ids) {
	f(n, ids);
}

static void glGetProgramBinary(_glGetProgramBinary f, GLuint program, GLsizei bufsize, GLsizei *length, GLenum *binaryFormat, void *binary) {
	f(program, bufsize, length, binaryFormat, binary);
}

static void glGetQueryObjectuiv(_glGetQueryObjectuiv f, GLuint id, GLenum pname, GLuint *params) {
	f(id, pname, params);
}

static void glMemoryBarrier(_glMemoryBarrier f, GLbitfield barriers) {
	f(barriers);
}

static void glDispatchCompute(_glDispatchCompute f, GLuint x, GLuint y, GLuint z) {
	f(x, y, z);
}

static void *glMapBufferRange(_glMapBufferRange f, GLenum target, GLintptr offset, GLsizeiptr length, GLbitfield access) {
	return f(target, offset, length, access);
}

static GLboolean glUnmapBuffer(_glUnmapBuffer f, GLenum target) {
	return f(target);
}

static void glBindImageTexture(_glBindImageTexture f, GLuint unit, GLuint texture, GLint level, GLboolean layered, GLint layer, GLenum access, GLenum format) {
	f(unit, texture, level, layered, layer, access, format);
}

static void glTexStorage2D(_glTexStorage2D f, GLenum target, GLsizei levels, GLenum internalFormat, GLsizei width, GLsizei height) {
	f(target, levels, internalFormat, width, height);
}

static void glBlitFramebuffer(_glBlitFramebuffer f, GLint srcX0, GLint srcY0, GLint srcX1, GLint srcY1, GLint dstX0, GLint dstY0, GLint dstX1, GLint dstY1, GLbitfield mask, GLenum filter) {
	f(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter);
}

static void  glowGetActiveUniform(_glGetActiveUniform f, GLuint  program, GLuint  index, GLsizei  bufSize, GLsizei * length, GLint * size, GLenum * type, GLchar * name) {
	f(program, index, bufSize, length, size, type, name);
}

static void  glowGetActiveAttrib(_glGetActiveAttrib f, GLuint  program, GLuint  index, GLsizei  bufSize, GLsizei * length, GLint * size, GLenum * type, GLchar * name) {
   f(program, index, bufSize, length, size, type, name);
}

static GLint  glowGetAttribLocation(_glGetAttribLocation f, GLuint  program, const GLchar * name) {
   return f(program, name);
}

static void  glowBlendFunc(_glBlendFunc f, GLenum  sfactor, GLenum  dfactor) {
   f(sfactor, dfactor);
}

static void  glowUniformMatrix2fv(_glUniformMatrix2fv f, GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value) {
   f(location, count, transpose, value);
}

static void  glowUniformMatrix3fv(_glUniformMatrix3fv f, GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value) {
   f(location, count, transpose, value);
}

static void  glowUniformMatrix4fv(_glUniformMatrix4fv f, GLint  location, GLsizei  count, GLboolean  transpose, const GLfloat * value) {
   f(location, count, transpose, value);
}

static void  glowUniform1fv(_glUniform1fv f, GLint  location, GLsizei  count, const GLfloat * value) {
   f(location, count, value);
}

static void  glowUniform2fv(_glUniform2fv f, GLint  location, GLsizei  count, const GLfloat * value) {
   f(location, count, value);
}

static void  glowUniform3fv(_glUniform3fv f, GLint  location, GLsizei  count, const GLfloat * value) {
   f(location, count, value);
}

static void  glowUniform4fv(_glUniform4fv f, GLint  location, GLsizei  count, const GLfloat * value) {
   f(location, count, value);
}
*/
import "C"

type Context interface{}

type Functions struct {
	// Query caches.
	uints  [100]C.GLuint
	ints   [100]C.GLint
	floats [100]C.GLfloat

	glActiveTexture                       C._glActiveTexture
	glAttachShader                        C._glAttachShader
	glBindAttribLocation                  C._glBindAttribLocation
	glBindBuffer                          C._glBindBuffer
	glBindFramebuffer                     C._glBindFramebuffer
	glBindRenderbuffer                    C._glBindRenderbuffer
	glBindTexture                         C._glBindTexture
	glBlendEquation                       C._glBlendEquation
	glBlendFuncSeparate                   C._glBlendFuncSeparate
	glBufferData                          C._glBufferData
	glBufferSubData                       C._glBufferSubData
	glCheckFramebufferStatus              C._glCheckFramebufferStatus
	glClear                               C._glClear
	glClearColor                          C._glClearColor
	glClearDepthf                         C._glClearDepthf
	glCompileShader                       C._glCompileShader
	glCopyTexSubImage2D                   C._glCopyTexSubImage2D
	glCreateProgram                       C._glCreateProgram
	glCreateShader                        C._glCreateShader
	glDeleteBuffers                       C._glDeleteBuffers
	glDeleteFramebuffers                  C._glDeleteFramebuffers
	glDeleteProgram                       C._glDeleteProgram
	glDeleteRenderbuffers                 C._glDeleteRenderbuffers
	glDeleteShader                        C._glDeleteShader
	glDeleteTextures                      C._glDeleteTextures
	glDepthFunc                           C._glDepthFunc
	glDepthMask                           C._glDepthMask
	glDisable                             C._glDisable
	glDisableVertexAttribArray            C._glDisableVertexAttribArray
	glDrawArrays                          C._glDrawArrays
	glDrawElements                        C._glDrawElements
	glEnable                              C._glEnable
	glEnableVertexAttribArray             C._glEnableVertexAttribArray
	glFinish                              C._glFinish
	glFlush                               C._glFlush
	glFramebufferRenderbuffer             C._glFramebufferRenderbuffer
	glFramebufferTexture2D                C._glFramebufferTexture2D
	glGenBuffers                          C._glGenBuffers
	glGenerateMipmap                      C._glGenerateMipmap
	glGenFramebuffers                     C._glGenFramebuffers
	glGenRenderbuffers                    C._glGenRenderbuffers
	glGenTextures                         C._glGenTextures
	glGetError                            C._glGetError
	glGetFramebufferAttachmentParameteriv C._glGetFramebufferAttachmentParameteriv
	glGetFloatv                           C._glGetFloatv
	glGetIntegerv                         C._glGetIntegerv
	glGetIntegeri_v                       C._glGetIntegeri_v
	glGetProgramiv                        C._glGetProgramiv
	glGetProgramInfoLog                   C._glGetProgramInfoLog
	glGetRenderbufferParameteriv          C._glGetRenderbufferParameteriv
	glGetShaderiv                         C._glGetShaderiv
	glGetShaderInfoLog                    C._glGetShaderInfoLog
	glGetString                           C._glGetString
	glGetUniformLocation                  C._glGetUniformLocation
	glGetVertexAttribiv                   C._glGetVertexAttribiv
	glGetVertexAttribPointerv             C._glGetVertexAttribPointerv
	glIsEnabled                           C._glIsEnabled
	glLinkProgram                         C._glLinkProgram
	glPixelStorei                         C._glPixelStorei
	glReadPixels                          C._glReadPixels
	glRenderbufferStorage                 C._glRenderbufferStorage
	glScissor                             C._glScissor
	glShaderSource                        C._glShaderSource
	glTexImage2D                          C._glTexImage2D
	glTexParameteri                       C._glTexParameteri
	glTexSubImage2D                       C._glTexSubImage2D
	glUniform1f                           C._glUniform1f
	glUniform1i                           C._glUniform1i
	glUniform2f                           C._glUniform2f
	glUniform3f                           C._glUniform3f
	glUniform4f                           C._glUniform4f
	glUseProgram                          C._glUseProgram
	glVertexAttribPointer                 C._glVertexAttribPointer
	glViewport                            C._glViewport
	glBindVertexArray                     C._glBindVertexArray
	glBindBufferBase                      C._glBindBufferBase
	glGetUniformBlockIndex                C._glGetUniformBlockIndex
	glUniformBlockBinding                 C._glUniformBlockBinding
	glInvalidateFramebuffer               C._glInvalidateFramebuffer
	glBeginQuery                          C._glBeginQuery
	glDeleteQueries                       C._glDeleteQueries
	glDeleteVertexArrays                  C._glDeleteVertexArrays
	glEndQuery                            C._glEndQuery
	glGenQueries                          C._glGenQueries
	glGenVertexArrays                     C._glGenVertexArrays
	glGetProgramBinary                    C._glGetProgramBinary
	glGetQueryObjectuiv                   C._glGetQueryObjectuiv
	glGetStringi                          C._glGetStringi
	glDispatchCompute                     C._glDispatchCompute
	glMemoryBarrier                       C._glMemoryBarrier
	glMapBufferRange                      C._glMapBufferRange
	glUnmapBuffer                         C._glUnmapBuffer
	glBindImageTexture                    C._glBindImageTexture
	glTexStorage2D                        C._glTexStorage2D
	glBlitFramebuffer                     C._glBlitFramebuffer

	glowGetActiveUniform  C._glGetActiveUniform
	glowGetActiveAttrib   C._glGetActiveAttrib
	glowGetAttribLocation C._glGetAttribLocation
	glowBlendFunc         C._glBlendFunc
	glowUniformMatrix2fv  C._glUniformMatrix2fv
	glowUniformMatrix3fv  C._glUniformMatrix3fv
	glowUniformMatrix4fv  C._glUniformMatrix4fv
	glowUniform1fv        C._glUniform1fv
	glowUniform2fv        C._glUniform2fv
	glowUniform3fv        C._glUniform3fv
	glowUniform4fv        C._glUniform4fv
}

// https://github.com/YouROK/go-mpv
func NewFunctions(ctx Context, forceES bool) (*Functions, error) {
	if ctx != nil {
		panic("non-nil context")
	}
	f := new(Functions)

	err := f.load(forceES)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func dlsym(handle unsafe.Pointer, s string) unsafe.Pointer {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return C.dlsym(handle, cs)
}

func dlopen(lib string) unsafe.Pointer {
	clib := C.CString(lib)
	defer C.free(unsafe.Pointer(clib))
	return C.dlopen(clib, C.RTLD_NOW|C.RTLD_LOCAL)
}

func (fn *Functions) load(forceES bool) error {
	var (
		loadErr  error
		libNames []string
		handles  []unsafe.Pointer
	)
	switch {
	case runtime.GOOS == "darwin" && !forceES:
		libNames = []string{"/System/Library/Frameworks/OpenGL.framework/OpenGL"}
	case runtime.GOOS == "darwin" && forceES:
		libNames = []string{"libGLESv2.dylib"}
	case runtime.GOOS == "ios":
		libNames = []string{"/System/Library/Frameworks/OpenGLES.framework/OpenGLES"}
	case runtime.GOOS == "android":
		libNames = []string{"libGLESv2.so", "libGLESv3.so"}
	default:
		libNames = []string{"libGLESv2.so.2", "libGLESv2.so.3.0"}
	}

	for _, lib := range libNames {
		if h := dlopen(lib); h != nil {
			handles = append(handles, h)
		}
	}
	if len(handles) == 0 {
		return fmt.Errorf("gl: no OpenGL implementation could be loaded (tried %q)", libNames)
	}

	load := func(s string) *[0]byte {
		for _, h := range handles {
			if f := dlsym(h, s); f != nil {
				return (*[0]byte)(f)
			}
		}

		// Try glGetProcAddress second
		if f := GetProcAddress(s); f != nil {
			fmt.Println("glGetProcAddress:", s, "->", C.GoString((*C.char)(f)))
			return (*[0]byte)(f)
		}

		return nil
	}

	must := func(s string) *[0]byte {
		ptr := load(s)
		if ptr == nil {
			loadErr = fmt.Errorf("gl: failed to load symbol %q", s)
		}
		return ptr
	}

	// GL ES 2.0 functions.
	fn.glActiveTexture = must("glActiveTexture")
	fn.glAttachShader = must("glAttachShader")
	fn.glBindAttribLocation = must("glBindAttribLocation")
	fn.glBindBuffer = must("glBindBuffer")
	fn.glBindFramebuffer = must("glBindFramebuffer")
	fn.glBindRenderbuffer = must("glBindRenderbuffer")
	fn.glBindTexture = must("glBindTexture")
	fn.glBlendEquation = must("glBlendEquation")
	fn.glBlendFuncSeparate = must("glBlendFuncSeparate")
	fn.glBufferData = must("glBufferData")
	fn.glBufferSubData = must("glBufferSubData")
	fn.glCheckFramebufferStatus = must("glCheckFramebufferStatus")
	fn.glClear = must("glClear")
	fn.glClearColor = must("glClearColor")
	fn.glClearDepthf = must("glClearDepthf")
	fn.glCompileShader = must("glCompileShader")
	fn.glCopyTexSubImage2D = must("glCopyTexSubImage2D")
	fn.glCreateProgram = must("glCreateProgram")
	fn.glCreateShader = must("glCreateShader")
	fn.glDeleteBuffers = must("glDeleteBuffers")
	fn.glDeleteFramebuffers = must("glDeleteFramebuffers")
	fn.glDeleteProgram = must("glDeleteProgram")
	fn.glDeleteRenderbuffers = must("glDeleteRenderbuffers")
	fn.glDeleteShader = must("glDeleteShader")
	fn.glDeleteTextures = must("glDeleteTextures")
	fn.glDepthFunc = must("glDepthFunc")
	fn.glDepthMask = must("glDepthMask")
	fn.glDisable = must("glDisable")
	fn.glDisableVertexAttribArray = must("glDisableVertexAttribArray")
	fn.glDrawArrays = must("glDrawArrays")
	fn.glDrawElements = must("glDrawElements")
	fn.glEnable = must("glEnable")
	fn.glEnableVertexAttribArray = must("glEnableVertexAttribArray")
	fn.glFinish = must("glFinish")
	fn.glFlush = must("glFlush")
	fn.glFramebufferRenderbuffer = must("glFramebufferRenderbuffer")
	fn.glFramebufferTexture2D = must("glFramebufferTexture2D")
	fn.glGenBuffers = must("glGenBuffers")
	fn.glGenerateMipmap = must("glGenerateMipmap")
	fn.glGenFramebuffers = must("glGenFramebuffers")
	fn.glGenRenderbuffers = must("glGenRenderbuffers")
	fn.glGenTextures = must("glGenTextures")
	fn.glGetError = must("glGetError")
	fn.glGetFramebufferAttachmentParameteriv = must("glGetFramebufferAttachmentParameteriv")
	fn.glGetIntegerv = must("glGetIntegerv")
	fn.glGetFloatv = must("glGetFloatv")
	fn.glGetProgramiv = must("glGetProgramiv")
	fn.glGetProgramInfoLog = must("glGetProgramInfoLog")
	fn.glGetRenderbufferParameteriv = must("glGetRenderbufferParameteriv")
	fn.glGetShaderiv = must("glGetShaderiv")
	fn.glGetShaderInfoLog = must("glGetShaderInfoLog")
	fn.glGetString = must("glGetString")
	fn.glGetUniformLocation = must("glGetUniformLocation")
	fn.glGetVertexAttribiv = must("glGetVertexAttribiv")
	fn.glGetVertexAttribPointerv = must("glGetVertexAttribPointerv")
	fn.glIsEnabled = must("glIsEnabled")
	fn.glLinkProgram = must("glLinkProgram")
	fn.glPixelStorei = must("glPixelStorei")
	fn.glReadPixels = must("glReadPixels")
	fn.glRenderbufferStorage = must("glRenderbufferStorage")
	fn.glScissor = must("glScissor")
	fn.glShaderSource = must("glShaderSource")
	fn.glTexImage2D = must("glTexImage2D")
	fn.glTexParameteri = must("glTexParameteri")
	fn.glTexSubImage2D = must("glTexSubImage2D")
	fn.glUniform1f = must("glUniform1f")
	fn.glUniform1i = must("glUniform1i")
	fn.glUniform2f = must("glUniform2f")
	fn.glUniform3f = must("glUniform3f")
	fn.glUniform4f = must("glUniform4f")
	fn.glUseProgram = must("glUseProgram")
	fn.glVertexAttribPointer = must("glVertexAttribPointer")
	fn.glViewport = must("glViewport")

	// Extensions and GL ES 3 functions.
	fn.glBindBufferBase = load("glBindBufferBase")
	fn.glBindVertexArray = load("glBindVertexArray")
	fn.glGetIntegeri_v = load("glGetIntegeri_v")
	fn.glGetUniformBlockIndex = load("glGetUniformBlockIndex")
	fn.glUniformBlockBinding = load("glUniformBlockBinding")
	fn.glInvalidateFramebuffer = load("glInvalidateFramebuffer")
	fn.glGetStringi = load("glGetStringi")
	// Fall back to EXT_invalidate_framebuffer if available.
	if fn.glInvalidateFramebuffer == nil {
		fn.glInvalidateFramebuffer = load("glDiscardFramebufferEXT")
	}

	fn.glBeginQuery = load("glBeginQuery")
	if fn.glBeginQuery == nil {
		fn.glBeginQuery = load("glBeginQueryEXT")
	}
	fn.glDeleteQueries = load("glDeleteQueries")
	if fn.glDeleteQueries == nil {
		fn.glDeleteQueries = load("glDeleteQueriesEXT")
	}
	fn.glEndQuery = load("glEndQuery")
	if fn.glEndQuery == nil {
		fn.glEndQuery = load("glEndQueryEXT")
	}
	fn.glGenQueries = load("glGenQueries")
	if fn.glGenQueries == nil {
		fn.glGenQueries = load("glGenQueriesEXT")
	}
	fn.glGetQueryObjectuiv = load("glGetQueryObjectuiv")
	if fn.glGetQueryObjectuiv == nil {
		fn.glGetQueryObjectuiv = load("glGetQueryObjectuivEXT")
	}

	fn.glDeleteVertexArrays = load("glDeleteVertexArrays")
	fn.glGenVertexArrays = load("glGenVertexArrays")
	fn.glMemoryBarrier = load("glMemoryBarrier")
	fn.glDispatchCompute = load("glDispatchCompute")
	fn.glMapBufferRange = load("glMapBufferRange")
	fn.glUnmapBuffer = load("glUnmapBuffer")
	fn.glBindImageTexture = load("glBindImageTexture")
	fn.glTexStorage2D = load("glTexStorage2D")
	fn.glBlitFramebuffer = load("glBlitFramebuffer")
	fn.glGetProgramBinary = load("glGetProgramBinary")

	fn.glowGetActiveUniform = load("glGetActiveUniform")
	fn.glowGetActiveAttrib = load("glGetActiveAttrib")
	fn.glowGetAttribLocation = load("glGetAttribLocation")
	fn.glowBlendFunc = load("glBlendFunc")
	fn.glowUniformMatrix2fv = load("glUniformMatrix2fv")
	fn.glowUniformMatrix3fv = load("glUniformMatrix3fv")
	fn.glowUniformMatrix4fv = load("glUniformMatrix4fv")
	fn.glowUniform1fv = load("glUniform1fv")
	fn.glowUniform2fv = load("glUniform2fv")
	fn.glowUniform3fv = load("glUniform3fv")
	fn.glowUniform4fv = load("glUniform4fv")

	return loadErr
}

func (fn *Functions) ActiveTexture(texture Enum) {
	C.glActiveTexture(fn.glActiveTexture, C.GLenum(texture))
}

func (fn *Functions) AttachShader(p Program, s Shader) {
	C.glAttachShader(fn.glAttachShader, C.GLuint(p.V), C.GLuint(s.V))
}

func (fn *Functions) BeginQuery(target Enum, query Query) {
	C.glBeginQuery(fn.glBeginQuery, C.GLenum(target), C.GLenum(query.V))
}

func (fn *Functions) BindAttribLocation(p Program, a Attrib, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.glBindAttribLocation(fn.glBindAttribLocation, C.GLuint(p.V), C.GLuint(a), cname)
}

func (fn *Functions) BindBufferBase(target Enum, index int, b Buffer) {
	C.glBindBufferBase(fn.glBindBufferBase, C.GLenum(target), C.GLuint(index), C.GLuint(b.V))
}

func (fn *Functions) BindBuffer(target Enum, b Buffer) {
	C.glBindBuffer(fn.glBindBuffer, C.GLenum(target), C.GLuint(b.V))
}

func (fn *Functions) BindFramebuffer(target Enum, fb Framebuffer) {
	C.glBindFramebuffer(fn.glBindFramebuffer, C.GLenum(target), C.GLuint(fb.V))
}

func (fn *Functions) BindRenderbuffer(target Enum, fb Renderbuffer) {
	C.glBindRenderbuffer(fn.glBindRenderbuffer, C.GLenum(target), C.GLuint(fb.V))
}

func (fn *Functions) BindImageTexture(unit int, t Texture, level int, layered bool, layer int, access, format Enum) {
	l := C.GLboolean(FALSE)
	if layered {
		l = TRUE
	}
	C.glBindImageTexture(fn.glBindImageTexture, C.GLuint(unit), C.GLuint(t.V), C.GLint(level), l, C.GLint(layer), C.GLenum(access), C.GLenum(format))
}

func (fn *Functions) BindTexture(target Enum, t Texture) {
	C.glBindTexture(fn.glBindTexture, C.GLenum(target), C.GLuint(t.V))
}

func (fn *Functions) BindVertexArray(a VertexArray) {
	C.glBindVertexArray(fn.glBindVertexArray, C.GLuint(a.V))
}

func (fn *Functions) BlendEquation(mode Enum) {
	C.glBlendEquation(fn.glBlendEquation, C.GLenum(mode))
}

func (fn *Functions) BlendFuncSeparate(srcRGB, dstRGB, srcA, dstA Enum) {
	C.glBlendFuncSeparate(fn.glBlendFuncSeparate, C.GLenum(srcRGB), C.GLenum(dstRGB), C.GLenum(srcA), C.GLenum(dstA))
}

func (fn *Functions) BlitFramebuffer(sx0, sy0, sx1, sy1, dx0, dy0, dx1, dy1 int, mask Enum, filter Enum) {
	C.glBlitFramebuffer(fn.glBlitFramebuffer,
		C.GLint(sx0), C.GLint(sy0), C.GLint(sx1), C.GLint(sy1),
		C.GLint(dx0), C.GLint(dy0), C.GLint(dx1), C.GLint(dy1),
		C.GLenum(mask), C.GLenum(filter),
	)
}

func (fn *Functions) BufferData(target Enum, data []byte, usage Enum) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glBufferData(fn.glBufferData, C.GLenum(target), C.GLsizeiptr(len(data)), p, C.GLenum(usage))
}

func (fn *Functions) BufferSubData(target Enum, offset int, src []byte) {
	var p unsafe.Pointer
	if len(src) > 0 {
		p = unsafe.Pointer(&src[0])
	}
	C.glBufferSubData(fn.glBufferSubData, C.GLenum(target), C.GLintptr(offset), C.GLsizeiptr(len(src)), p)
}

func (fn *Functions) CheckFramebufferStatus(target Enum) Enum {
	return Enum(C.glCheckFramebufferStatus(fn.glCheckFramebufferStatus, C.GLenum(target)))
}

func (fn *Functions) Clear(mask Enum) {
	C.glClear(fn.glClear, C.GLbitfield(mask))
}

func (fn *Functions) ClearColor(red float32, green float32, blue float32, alpha float32) {
	C.glClearColor(fn.glClearColor, C.GLfloat(red), C.GLfloat(green), C.GLfloat(blue), C.GLfloat(alpha))
}

func (fn *Functions) ClearDepthf(d float32) {
	C.glClearDepthf(fn.glClearDepthf, C.GLfloat(d))
}

func (fn *Functions) CompileShader(s Shader) {
	C.glCompileShader(fn.glCompileShader, C.GLuint(s.V))
}

func (fn *Functions) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	C.glCopyTexSubImage2D(fn.glCopyTexSubImage2D, C.GLenum(target), C.GLint(level), C.GLint(xoffset), C.GLint(yoffset), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func (fn *Functions) CreateBuffer() Buffer {
	C.glGenBuffers(fn.glGenBuffers, 1, &fn.uints[0])
	return Buffer{uint(fn.uints[0])}
}

func (fn *Functions) CreateFramebuffer() Framebuffer {
	C.glGenFramebuffers(fn.glGenFramebuffers, 1, &fn.uints[0])
	return Framebuffer{uint(fn.uints[0])}
}

func (fn *Functions) CreateProgram() Program {
	return Program{uint(C.glCreateProgram(fn.glCreateProgram))}
}

func (fn *Functions) CreateQuery() Query {
	C.glGenQueries(fn.glGenQueries, 1, &fn.uints[0])
	return Query{uint(fn.uints[0])}
}

func (fn *Functions) CreateRenderbuffer() Renderbuffer {
	C.glGenRenderbuffers(fn.glGenRenderbuffers, 1, &fn.uints[0])
	return Renderbuffer{uint(fn.uints[0])}
}

func (fn *Functions) CreateShader(ty Enum) Shader {
	return Shader{uint(C.glCreateShader(fn.glCreateShader, C.GLenum(ty)))}
}

func (fn *Functions) CreateTexture() Texture {
	C.glGenTextures(fn.glGenTextures, 1, &fn.uints[0])
	return Texture{uint(fn.uints[0])}
}

func (fn *Functions) CreateVertexArray() VertexArray {
	C.glGenVertexArrays(fn.glGenVertexArrays, 1, &fn.uints[0])
	return VertexArray{uint(fn.uints[0])}
}

func (fn *Functions) DeleteBuffer(v Buffer) {
	fn.uints[0] = C.GLuint(v.V)
	C.glDeleteBuffers(fn.glDeleteBuffers, 1, &fn.uints[0])
}

func (fn *Functions) DeleteFramebuffer(v Framebuffer) {
	fn.uints[0] = C.GLuint(v.V)
	C.glDeleteFramebuffers(fn.glDeleteFramebuffers, 1, &fn.uints[0])
}

func (fn *Functions) DeleteProgram(p Program) {
	C.glDeleteProgram(fn.glDeleteProgram, C.GLuint(p.V))
}

func (fn *Functions) DeleteQuery(query Query) {
	fn.uints[0] = C.GLuint(query.V)
	C.glDeleteQueries(fn.glDeleteQueries, 1, &fn.uints[0])
}

func (fn *Functions) DeleteVertexArray(array VertexArray) {
	fn.uints[0] = C.GLuint(array.V)
	C.glDeleteVertexArrays(fn.glDeleteVertexArrays, 1, &fn.uints[0])
}

func (fn *Functions) DeleteRenderbuffer(v Renderbuffer) {
	fn.uints[0] = C.GLuint(v.V)
	C.glDeleteRenderbuffers(fn.glDeleteRenderbuffers, 1, &fn.uints[0])
}

func (fn *Functions) DeleteShader(s Shader) {
	C.glDeleteShader(fn.glDeleteShader, C.GLuint(s.V))
}

func (fn *Functions) DeleteTexture(v Texture) {
	fn.uints[0] = C.GLuint(v.V)
	C.glDeleteTextures(fn.glDeleteTextures, 1, &fn.uints[0])
}

func (fn *Functions) DepthFunc(v Enum) {
	C.glDepthFunc(fn.glDepthFunc, C.GLenum(v))
}

func (fn *Functions) DepthMask(mask bool) {
	m := C.GLboolean(FALSE)
	if mask {
		m = C.GLboolean(TRUE)
	}
	C.glDepthMask(fn.glDepthMask, m)
}

func (fn *Functions) DisableVertexAttribArray(a Attrib) {
	C.glDisableVertexAttribArray(fn.glDisableVertexAttribArray, C.GLuint(a))
}

func (fn *Functions) Disable(cap Enum) {
	C.glDisable(fn.glDisable, C.GLenum(cap))
}

func (fn *Functions) DrawArrays(mode Enum, first int, count int) {
	C.glDrawArrays(fn.glDrawArrays, C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}

func (fn *Functions) DrawElements(mode Enum, count int, ty Enum, offset int) {
	C.glDrawElements(fn.glDrawElements, C.GLenum(mode), C.GLsizei(count), C.GLenum(ty), C.uintptr_t(offset))
}

func (fn *Functions) DispatchCompute(x, y, z int) {
	C.glDispatchCompute(fn.glDispatchCompute, C.GLuint(x), C.GLuint(y), C.GLuint(z))
}

func (fn *Functions) Enable(cap Enum) {
	C.glEnable(fn.glEnable, C.GLenum(cap))
}

func (fn *Functions) EndQuery(target Enum) {
	C.glEndQuery(fn.glEndQuery, C.GLenum(target))
}

func (fn *Functions) EnableVertexAttribArray(a Attrib) {
	C.glEnableVertexAttribArray(fn.glEnableVertexAttribArray, C.GLuint(a))
}

func (fn *Functions) Finish() {
	C.glFinish(fn.glFinish)
}

func (fn *Functions) Flush() {
	C.glFlush(fn.glFinish)
}

func (fn *Functions) FramebufferRenderbuffer(target, attachment, renderbuffertarget Enum, renderbuffer Renderbuffer) {
	C.glFramebufferRenderbuffer(fn.glFramebufferRenderbuffer, C.GLenum(target), C.GLenum(attachment), C.GLenum(renderbuffertarget), C.GLuint(renderbuffer.V))
}

func (fn *Functions) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	C.glFramebufferTexture2D(fn.glFramebufferTexture2D, C.GLenum(target), C.GLenum(attachment), C.GLenum(texTarget), C.GLuint(t.V), C.GLint(level))
}

func (fn *Functions) GenerateMipmap(target Enum) {
	C.glGenerateMipmap(fn.glGenerateMipmap, C.GLenum(target))
}

func (fn *Functions) GetBinding(pname Enum) Object {
	return Object{uint(fn.GetInteger(pname))}
}

func (fn *Functions) GetBindingi(pname Enum, idx int) Object {
	return Object{uint(fn.GetIntegeri(pname, idx))}
}

func (fn *Functions) GetError() Enum {
	return Enum(C.glGetError(fn.glGetError))
}

func (fn *Functions) GetRenderbufferParameteri(target, pname Enum) int {
	C.glGetRenderbufferParameteriv(fn.glGetRenderbufferParameteriv, C.GLenum(target), C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	C.glGetFramebufferAttachmentParameteriv(fn.glGetFramebufferAttachmentParameteriv, C.GLenum(target), C.GLenum(attachment), C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetFloat4(pname Enum) [4]float32 {
	C.glGetFloatv(fn.glGetFloatv, C.GLenum(pname), &fn.floats[0])
	var r [4]float32
	for i := range r {
		r[i] = float32(fn.floats[i])
	}
	return r
}

func (fn *Functions) GetFloat(pname Enum) float32 {
	C.glGetFloatv(fn.glGetFloatv, C.GLenum(pname), &fn.floats[0])
	return float32(fn.floats[0])
}

func (fn *Functions) GetInteger4(pname Enum) [4]int {
	C.glGetIntegerv(fn.glGetIntegerv, C.GLenum(pname), &fn.ints[0])
	var r [4]int
	for i := range r {
		r[i] = int(fn.ints[i])
	}
	return r
}

func (fn *Functions) GetInteger(pname Enum) int {
	C.glGetIntegerv(fn.glGetIntegerv, C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetIntegeri(pname Enum, idx int) int {
	C.glGetIntegeri_v(fn.glGetIntegeri_v, C.GLenum(pname), C.GLuint(idx), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetProgrami(p Program, pname Enum) int {
	C.glGetProgramiv(fn.glGetProgramiv, C.GLuint(p.V), C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetProgramBinary(p Program) []byte {
	sz := fn.GetProgrami(p, PROGRAM_BINARY_LENGTH)
	if sz == 0 {
		return nil
	}
	buf := make([]byte, sz)
	var format C.GLenum
	C.glGetProgramBinary(fn.glGetProgramBinary, C.GLuint(p.V), C.GLsizei(sz), nil, &format, unsafe.Pointer(&buf[0]))
	return buf
}

func (fn *Functions) GetProgramInfoLog(p Program) string {
	n := fn.GetProgrami(p, INFO_LOG_LENGTH)
	buf := make([]byte, n)
	C.glGetProgramInfoLog(fn.glGetProgramInfoLog, C.GLuint(p.V), C.GLsizei(len(buf)), nil, (*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf)
}

func (fn *Functions) GetQueryObjectuiv(query Query, pname Enum) uint {
	C.glGetQueryObjectuiv(fn.glGetQueryObjectuiv, C.GLuint(query.V), C.GLenum(pname), &fn.uints[0])
	return uint(fn.uints[0])
}

func (fn *Functions) GetShaderi(s Shader, pname Enum) int {
	C.glGetShaderiv(fn.glGetShaderiv, C.GLuint(s.V), C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetShaderInfoLog(s Shader) string {
	n := fn.GetShaderi(s, INFO_LOG_LENGTH)
	buf := make([]byte, n)
	C.glGetShaderInfoLog(fn.glGetShaderInfoLog, C.GLuint(s.V), C.GLsizei(len(buf)), nil, (*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf)
}

func (fn *Functions) getStringi(pname Enum, index int) string {
	str := C.glGetStringi(fn.glGetStringi, C.GLenum(pname), C.GLuint(index))
	if str == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(str)))
}

func (fn *Functions) GetString(pname Enum) string {
	switch {
	case runtime.GOOS == "darwin" && pname == EXTENSIONS:
		// macOS OpenGL 3 core profile doesn't support glGetString(GL_EXTENSIONS).
		// Use glGetStringi(GL_EXTENSIONS, <index>).
		var exts []string
		nexts := fn.GetInteger(NUM_EXTENSIONS)
		for i := 0; i < nexts; i++ {
			ext := fn.getStringi(EXTENSIONS, i)
			exts = append(exts, ext)
		}
		return strings.Join(exts, " ")
	default:
		str := C.glGetString(fn.glGetString, C.GLenum(pname))
		return C.GoString((*C.char)(unsafe.Pointer(str)))
	}
}

func (fn *Functions) GetUniformBlockIndex(p Program, name string) uint {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint(C.glGetUniformBlockIndex(fn.glGetUniformBlockIndex, C.GLuint(p.V), cname))
}

func (fn *Functions) GetUniformLocation(p Program, name string) Uniform {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Uniform{int(C.glGetUniformLocation(fn.glGetUniformLocation, C.GLuint(p.V), cname))}
}

func (fn *Functions) GetVertexAttrib(index int, pname Enum) int {
	C.glGetVertexAttribiv(fn.glGetVertexAttribiv, C.GLuint(index), C.GLenum(pname), &fn.ints[0])
	return int(fn.ints[0])
}

func (fn *Functions) GetVertexAttribBinding(index int, pname Enum) Object {
	return Object{uint(fn.GetVertexAttrib(index, pname))}
}

func (fn *Functions) GetVertexAttribPointer(index int, pname Enum) uintptr {
	ptr := C.glGetVertexAttribPointerv(fn.glGetVertexAttribPointerv, C.GLuint(index), C.GLenum(pname))
	return uintptr(ptr)
}

func (fn *Functions) InvalidateFramebuffer(target, attachment Enum) {
	C.glInvalidateFramebuffer(fn.glInvalidateFramebuffer, C.GLenum(target), C.GLenum(attachment))
}

func (fn *Functions) IsEnabled(cap Enum) bool {
	return C.glIsEnabled(fn.glIsEnabled, C.GLenum(cap)) == TRUE
}

func (fn *Functions) LinkProgram(p Program) {
	C.glLinkProgram(fn.glLinkProgram, C.GLuint(p.V))
}

func (fn *Functions) PixelStorei(pname Enum, param int) {
	C.glPixelStorei(fn.glPixelStorei, C.GLenum(pname), C.GLint(param))
}

func (fn *Functions) MemoryBarrier(barriers Enum) {
	C.glMemoryBarrier(fn.glMemoryBarrier, C.GLbitfield(barriers))
}

func (fn *Functions) MapBufferRange(target Enum, offset, length int, access Enum) []byte {
	p := C.glMapBufferRange(fn.glMapBufferRange, C.GLenum(target), C.GLintptr(offset), C.GLsizeiptr(length), C.GLbitfield(access))
	if p == nil {
		return nil
	}
	return (*[1 << 30]byte)(p)[:length:length]
}

func (fn *Functions) Scissor(x, y, width, height int32) {
	C.glScissor(fn.glScissor, C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func (fn *Functions) ReadPixels(x, y, width, height int, format, ty Enum, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glReadPixels(fn.glReadPixels, C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLenum(format), C.GLenum(ty), p)
}

func (fn *Functions) RenderbufferStorage(target, internalformat Enum, width, height int) {
	C.glRenderbufferStorage(fn.glRenderbufferStorage, C.GLenum(target), C.GLenum(internalformat), C.GLsizei(width), C.GLsizei(height))
}

func (fn *Functions) ShaderSource(s Shader, src string) {
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	strlen := C.GLint(len(src))
	C.glShaderSource(fn.glShaderSource, C.GLuint(s.V), 1, &csrc, &strlen)
}

func (fn *Functions) TexImage2D(target Enum, level int, width int, height int, format Enum, ty Enum, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glTexImage2D(fn.glTexImage2D, C.GLenum(target), C.GLint(level), C.GLint(format), C.GLsizei(width), C.GLsizei(height), 0, C.GLenum(format), C.GLenum(ty), p)
}

func (fn *Functions) TexStorage2D(target Enum, levels int, internalFormat Enum, width, height int) {
	C.glTexStorage2D(fn.glTexStorage2D, C.GLenum(target), C.GLsizei(levels), C.GLenum(internalFormat), C.GLsizei(width), C.GLsizei(height))
}

func (fn *Functions) TexSubImage2D(target Enum, level int, x int, y int, width int, height int, format Enum, ty Enum, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glTexSubImage2D(fn.glTexSubImage2D, C.GLenum(target), C.GLint(level), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLenum(format), C.GLenum(ty), p)
}

func (fn *Functions) TexParameteri(target, pname Enum, param int) {
	C.glTexParameteri(fn.glTexParameteri, C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

func (fn *Functions) UniformBlockBinding(p Program, uniformBlockIndex uint, uniformBlockBinding uint) {
	C.glUniformBlockBinding(fn.glUniformBlockBinding, C.GLuint(p.V), C.GLuint(uniformBlockIndex), C.GLuint(uniformBlockBinding))
}

func (fn *Functions) Uniform1f(dst Uniform, v float32) {
	C.glUniform1f(fn.glUniform1f, C.GLint(dst.V), C.GLfloat(v))
}

func (fn *Functions) Uniform1i(dst Uniform, v int) {
	C.glUniform1i(fn.glUniform1i, C.GLint(dst.V), C.GLint(v))
}

func (fn *Functions) Uniform2f(dst Uniform, v0 float32, v1 float32) {
	C.glUniform2f(fn.glUniform2f, C.GLint(dst.V), C.GLfloat(v0), C.GLfloat(v1))
}

func (fn *Functions) Uniform3f(dst Uniform, v0 float32, v1 float32, v2 float32) {
	C.glUniform3f(fn.glUniform3f, C.GLint(dst.V), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
}

func (fn *Functions) Uniform4f(dst Uniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	C.glUniform4f(fn.glUniform4f, C.GLint(dst.V), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
}

func (fn *Functions) UseProgram(p Program) {
	C.glUseProgram(fn.glUseProgram, C.GLuint(p.V))
}

func (fn *Functions) UnmapBuffer(target Enum) bool {
	r := C.glUnmapBuffer(fn.glUnmapBuffer, C.GLenum(target))
	return r == TRUE
}

func (fn *Functions) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride int, offset int) {
	var n C.GLboolean = FALSE
	if normalized {
		n = TRUE
	}
	C.glVertexAttribPointer(fn.glVertexAttribPointer, C.GLuint(dst), C.GLint(size), C.GLenum(ty), n, C.GLsizei(stride), C.uintptr_t(offset))
}

func (fn *Functions) Viewport(x int, y int, width int, height int) {
	C.glViewport(fn.glViewport, C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func (fn *Functions) BlendFunc(sfactor, dfactor Enum) {
	C.glowBlendFunc(fn.glowBlendFunc, C.GLenum(sfactor), C.GLenum(dfactor))
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func (fn *Functions) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	var length, si C.GLint
	var typ C.GLenum
	name = strings.Repeat("\x00", 256)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.glowGetActiveUniform(fn.glowGetActiveUniform, C.GLuint(p.V), C.GLuint(index), C.GLint(len(name)-1), &length, &si, &typ, cname)
	name = name[:strings.IndexRune(name, 0)]
	return name, int(si), Enum(typ)

}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func (fn *Functions) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	var length, si C.GLint
	var typ C.GLenum
	name = strings.Repeat("\x00", 256)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.glowGetActiveAttrib(fn.glowGetActiveAttrib, C.GLuint(p.V), C.GLuint(index), C.GLint(len(name)-1), &length, &si, &typ, cname)
	name = name[:strings.IndexRune(name, 0)]
	return name, int(si), Enum(typ)
}

func (fn *Functions) GetAttribLocation(p Program, name string) Attrib {
	cname := C.CString(name + "\x00")
	defer C.free(unsafe.Pointer(cname))
	return Attrib(uint(C.glowGetAttribLocation(fn.glowGetAttribLocation, C.GLuint(p.V), cname)))
}

func (fn *Functions) UniformMatrix2fv(dst Uniform, src []float32) {
	C.glowUniformMatrix2fv(fn.glowUniformMatrix2fv, C.GLint(dst.V), C.GLsizei(len(src)/(2*2)), C.GLboolean(FALSE), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) UniformMatrix3fv(dst Uniform, src []float32) {
	C.glowUniformMatrix3fv(fn.glowUniformMatrix3fv, C.GLint(dst.V), C.GLsizei(len(src)/(3*3)), C.GLboolean(FALSE), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) UniformMatrix4fv(dst Uniform, src []float32) {
	C.glowUniformMatrix4fv(fn.glowUniformMatrix4fv, C.GLint(dst.V), C.GLsizei(len(src)/(4*4)), C.GLboolean(FALSE), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) Uniform1fv(dst Uniform, src []float32) {
	C.glowUniform1fv(fn.glowUniform1fv, C.GLint(dst.V), C.GLsizei(len(src)), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) Uniform2fv(dst Uniform, src []float32) {
	C.glowUniform2fv(fn.glowUniform2fv, C.GLint(dst.V), C.GLsizei(len(src)/2), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) Uniform3fv(dst Uniform, src []float32) {
	C.glowUniform3fv(fn.glowUniform3fv, C.GLint(dst.V), C.GLsizei(len(src)/3), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}

func (fn *Functions) Uniform4fv(dst Uniform, src []float32) {
	C.glowUniform4fv(fn.glowUniform4fv, C.GLint(dst.V), C.GLsizei(len(src)/4), (*C.GLfloat)(unsafe.Pointer(&src[0])))
}
