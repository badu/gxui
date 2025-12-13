//go:build freebsd || linux || netbsd || openbsd

package purego

import (
	"errors"
	"fmt"
	"os"

	"github.com/ebitengine/purego"
)

func (f *Functions) init() error {
	var errs []error

	// Try OpenGL ES first. Some machines like Android and Raspberry Pi might work only with OpenGL ES.
	//
	// Do not use OpenGL ES for Steam, as overlays might not work properly (#3338).
	// With Steam, OpenGL (not ES) should be available anyway.
	if os.Getenv("SteamEnv") != "1" {
		for _, name := range []string{"libGLESv2.so", "libGLESv2.so.2", "libGLESv2.so.1", "libGLESv2.so.0"} {
			lib, err := purego.Dlopen(name, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
			if err == nil {
				f.libGLES = lib
				f.isES = true
				return nil
			}
			errs = append(errs, fmt.Errorf("gl: Dlopen failed: name: %s: %w", name, err))
		}
	}

	// Try OpenGL next.
	// Usually libGL.so or libGL.so.1 is used. libGL.so.2 might exist only on NetBSD.
	// TODO: Should "libOpenGL.so.0" [1] and "libGLX.so.0" [2] be added? These were added as of GLFW 3.3.9.
	// [1] https://github.com/glfw/glfw/commit/55aad3c37b67f17279378db52da0a3ab81bbf26d
	// [2] https://github.com/glfw/glfw/commit/c18851f52ec9704eb06464058a600845ec1eada1
	for _, name := range []string{"libGL.so", "libGL.so.2", "libGL.so.1", "libGL.so.0"} {
		lib, err := purego.Dlopen(name, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
		if err == nil {
			f.libGL = lib
			return nil
		}
		errs = append(errs, fmt.Errorf("gl: Dlopen failed: name: %s: %w", name, err))
	}

	errs = append([]error{fmt.Errorf("gl: failed to load libGL.so and libGLESv2.so: ")}, errs...)
	return errors.Join(errs...)
}

func (f *Functions) getProcAddress(name string) (uintptr, error) {
	if f.isES {
		return f.getProcAddressGLES(name)
	}
	return f.getProcAddressGL(name)
}

func (f *Functions) getProcAddressGL(name string) (uintptr, error) {
	if f.glXGetProcAddress == nil {
		if _, err := purego.Dlsym(f.libGL, "glXGetProcAddress"); err == nil {
			purego.RegisterLibFunc(&f.glXGetProcAddress, f.libGL, "glXGetProcAddress")
		} else if _, err := purego.Dlsym(f.libGL, "glXGetProcAddressARB"); err == nil {
			purego.RegisterLibFunc(&f.glXGetProcAddress, f.libGL, "glXGetProcAddressARB")
		}
	}

	if f.glXGetProcAddress == nil {
		return 0, fmt.Errorf("gl: failed to find glXGetProcAddress or glXGetProcAddressARB in libGL.so")
	}

	return f.glXGetProcAddress(name), nil
}

func (f *Functions) getProcAddressGLES(name string) (uintptr, error) {
	proc, err := purego.Dlsym(f.libGLES, name)
	if err != nil {
		return 0, err
	}
	return proc, nil
}
