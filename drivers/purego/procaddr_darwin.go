package purego

import (
	"fmt"

	"github.com/ebitengine/purego"
)

var (
	opengl uintptr
)

func (fn *Functions) init() error {
	lib, errGLES := purego.Dlopen("/System/Library/Frameworks/OpenGLES.framework/OpenGLES", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if errGLES == nil {
		fn.isES = true
		opengl = lib
		return nil
	}

	lib, errGL := purego.Dlopen("/System/Library/Frameworks/OpenGL.framework/OpenGL", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if errGL == nil {
		opengl = lib
		return nil
	}

	return fmt.Errorf("gl: failed to load: OpenGL.framework: %w, OpenGLES.framework: %w", errGL, errGLES)
}

func (fn *Functions) getProcAddress(name string) (uintptr, error) {
	proc, err := purego.Dlsym(opengl, name)
	if err != nil {
		return 0, err
	}
	return proc, nil
}
