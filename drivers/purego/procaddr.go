//go:build !js

package purego

import (
	"fmt"
)

type procAddressGetter struct {
	ctx *Functions
	err error
}

func (p *procAddressGetter) get(name string) uintptr {
	proc, err := p.ctx.getProcAddress(name)
	if err != nil {
		p.err = fmt.Errorf("gl: %s is missing: %w", name, err)
		return 0
	}

	if proc == 0 {
		p.err = fmt.Errorf("gl: %s is missing", name)
		return 0
	}

	return proc
}

func (p *procAddressGetter) error() error {
	return p.err
}
