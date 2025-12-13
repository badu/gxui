package purego

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	opengl32              = windows.NewLazySystemDLL("opengl32")
	procWglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
)

func (c *Functions) init() error {
	return nil
}

func (c *Functions) getProcAddress(namea string) (uintptr, error) {
	cname, err := windows.BytePtrFromString(namea)
	if err != nil {
		return 0, err
	}

	r, _, err := procWglGetProcAddress.Call(uintptr(unsafe.Pointer(cname)))
	if r != 0 {
		return r, nil
	}
	if err != nil && err != windows.ERROR_SUCCESS {
		if errors.Is(err, windows.ERROR_PROC_NOT_FOUND) {
			return 0, fmt.Errorf("gl: wglGetProcAddress not found for %s: %w", namea, err)
		}
		return 0, fmt.Errorf("gl: wglGetProcAddress failed for %s: %w", namea, err)
	}

	p := opengl32.NewProc(namea)
	if err := p.Find(); err != nil {
		return 0, err
	}

	c.libGL = p.Addr()

	return p.Addr(), nil
}
