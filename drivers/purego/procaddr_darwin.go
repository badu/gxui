package purego

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func (f *Functions) init() error {
	var errGLES error
	f.libGLES, errGLES = purego.Dlopen("/System/Library/Frameworks/OpenGLES.framework/OpenGLES", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if errGLES == nil {
		f.isES = true
		return nil
	}

	var errGL error
	f.libGL, errGL = purego.Dlopen("/System/Library/Frameworks/OpenGL.framework/OpenGL", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if errGL == nil {
		return nil
	}

	return fmt.Errorf("gl: failed to load: OpenGL.framework: %w, OpenGLES.framework: %w", errGL, errGLES)
}

func (f *Functions) getProcAddress(name string) (uintptr, error) {
	proc, err := purego.Dlsym(opengl, name)
	if err != nil {
		return 0, err
	}

	return proc, nil
}
