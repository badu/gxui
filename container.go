package gxui

import "github.com/badu/gxui/pkg/math"

type Container interface {
	Parent
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
}
