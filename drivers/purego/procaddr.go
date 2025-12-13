//go:build !js

package purego

import (
	"fmt"
)

func (f *Functions) get(name string) (uintptr, error) {
	proc, err := f.getProcAddress(name)
	if err != nil {
		return 0, fmt.Errorf("gl: %s is missing: %w", name, err)
	}

	if proc == 0 {
		return 0, fmt.Errorf("gl: %s is missing", name)
	}

	return proc, nil
}
