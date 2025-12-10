//go:build freebsd || linux || netbsd || openbsd

package purego

import (
	"errors"
	"fmt"
	"os"

	"github.com/ebitengine/purego"
)

var (
	libGL   uintptr
	libGLES uintptr
)

func (fn *Functions) init() error {
	var errs []error

	// Try OpenGL ES first. Some machines like Android and Raspberry Pi might work only with OpenGL ES.
	//
	// Do not use OpenGL ES for Steam, as overlays might not work properly (#3338).
	// With Steam, OpenGL (not ES) should be available anyway.
	if os.Getenv("SteamEnv") != "1" {
		for _, name := range []string{"libGLESv2.so", "libGLESv2.so.2", "libGLESv2.so.1", "libGLESv2.so.0"} {
			lib, err := purego.Dlopen(name, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
			if err == nil {
				libGLES = lib
				fn.isES = true
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
			libGL = lib
			return nil
		}
		errs = append(errs, fmt.Errorf("gl: Dlopen failed: name: %s: %w", name, err))
	}

	errs = append([]error{fmt.Errorf("gl: failed to load libGL.so and libGLESv2.so: ")}, errs...)
	return errors.Join(errs...)
}

func (fn *Functions) getProcAddress(name string) (uintptr, error) {
	if fn.isES {
		return getProcAddressGLES(name)
	}
	return getProcAddressGL(name)
}

var glXGetProcAddress func(name string) uintptr

func getProcAddressGL(name string) (uintptr, error) {
	if glXGetProcAddress == nil {
		if _, err := purego.Dlsym(libGL, "glXGetProcAddress"); err == nil {
			purego.RegisterLibFunc(&glXGetProcAddress, libGL, "glXGetProcAddress")
		} else if _, err := purego.Dlsym(libGL, "glXGetProcAddressARB"); err == nil {
			purego.RegisterLibFunc(&glXGetProcAddress, libGL, "glXGetProcAddressARB")
		}
	}

	if glXGetProcAddress == nil {
		return 0, fmt.Errorf("gl: failed to find glXGetProcAddress or glXGetProcAddressARB in libGL.so")
	}

	return glXGetProcAddress(name), nil
}

func getProcAddressGLES(name string) (uintptr, error) {
	proc, err := purego.Dlsym(libGLES, name)
	if err != nil {
		return 0, err
	}
	return proc, nil
}
