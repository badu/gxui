// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/badu/gxui/pkg/math"
)

type ParentPoint struct {
	Control Parent
	Point   math.Point
}

type ControlPoint struct {
	Control Control
	Point   math.Point
}

type ControlPointList []ControlPoint

func (l ControlPointList) Contains(control Control) bool {
	_, found := l.Find(control)
	return found
}

func (l ControlPointList) Find(control Control) (math.Point, bool) {
	for _, i := range l {
		if i.Control == control {
			return i.Point, true
		}
	}
	return math.Point{}, false
}

func ValidateHierarchy(parent Parent) {
	for _, child := range parent.Children() {
		if parent != child.Control.Parent() {
			panic(fmt.Errorf("Child's parent is not as expected.\nChild: %s\nExpected parent: %s",
				Path(child.Control), Path(parent)))
		}
		if cp, ok := child.Control.(Parent); ok {
			ValidateHierarchy(cp)
		}
	}
}

func CommonAncestor(control, otherControl Control) Parent {
	seen := make(map[Parent]bool)
	if parent, _ := control.(Parent); parent != nil {
		seen[parent] = true
	}
	for control != nil {
		p := control.Parent()
		seen[p] = true
		control, _ = p.(Control)
	}
	if c, _ := otherControl.(Parent); c != nil {
		if seen[c] {
			return c
		}
	}
	for otherControl != nil {
		p := otherControl.Parent()
		if seen[p] {
			return p
		}
		otherControl, _ = p.(Control)
	}
	return nil
}

func TopControlsUnder(point math.Point, parent Parent) ControlPointList {
	children := parent.Children()
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		cp := point.Sub(child.Offset)
		if child.Control.ContainsPoint(cp) {
			l := ControlPointList{ControlPoint{Control: child.Control, Point: cp}}
			if cc, ok := child.Control.(Parent); ok {
				l = append(l, TopControlsUnder(cp, cc)...)
			}
			return l
		}
	}
	return ControlPointList{}
}

func ControlsUnder(point math.Point, parent Parent) ControlPointList {
	toVisit := []ParentPoint{{Control: parent, Point: point}}
	l := ControlPointList{}
	for len(toVisit) > 0 {
		parent = toVisit[0].Control
		point = toVisit[0].Point
		toVisit = toVisit[1:]
		for _, child := range parent.Children() {
			cp := point.Sub(child.Offset)
			if child.Control.ContainsPoint(cp) {
				l = append(l, ControlPoint{child.Control, cp})
				if cc, ok := child.Control.(Parent); ok {
					toVisit = append(toVisit, ParentPoint{cc, cp})
				}
			}
		}
	}
	return l
}

func WindowToChild(point math.Point, control Control) math.Point {
	ctrl := control
	for {
		parent := ctrl.Parent()
		if parent == nil {
			panic("Control's parent was nil")
		}

		child := parent.Children().Find(ctrl)
		if child == nil {
			Dump(parent)
			panic(fmt.Errorf("Control's parent (%p %T) did not contain control (%p %T).", &parent, parent, &ctrl, ctrl))
		}

		point = point.Sub(child.Offset)

		switch parent.(type) {
		case *WindowImpl:
			return point
		}

		var ok bool
		ctrl, ok = parent.(Control)
		if !ok {
			Dump(parent)
			panic(fmt.Errorf("WindowToChild (%p %T) -> (%p %T) reached non-control parent (%p %T).", &control, control, &parent, parent, &ctrl, ctrl))
		}
	}
}

func ChildToParent(point math.Point, onControl Control, parent Parent) math.Point {
	control := onControl
	for {
		parent := control.Parent()
		if parent == nil {
			panic(fmt.Errorf("Control detached: %s", Path(control)))
		}

		child := parent.Children().Find(control)
		if child == nil {
			Dump(parent)
			panic(fmt.Errorf("Control's parent (%p %T) did not contain control (%p %T).", &parent, parent, &control, control))
		}

		point = point.Add(child.Offset)
		if parent == parent {
			return point
		}

		var ok bool
		control, ok = parent.(Control)
		if !ok {
			Dump(parent)
			panic(fmt.Errorf("ChildToParent (%p %T) -> (%p %T) reached non-control parent (%p %T).", &onControl, onControl, &parent, parent, &parent, parent))
		}
	}
}

func ParentToChild(point math.Point, parent Parent, control Control) math.Point {
	return point.Sub(ChildToParent(math.ZeroPoint, control, parent))
}

func TransformCoordinate(point math.Point, fromControl, toControl Control) math.Point {
	if fromControl == toControl {
		return point
	}

	ancestor := CommonAncestor(fromControl, toControl)
	if ancestor == nil {
		panic(fmt.Errorf("no common ancestor between %s and %s", Path(fromControl), Path(toControl)))
	}

	if parent, ok := ancestor.(Control); !ok || parent != fromControl {
		point = ChildToParent(point, fromControl, ancestor)
	}

	if parent, ok := ancestor.(Control); !ok || parent != toControl {
		point = ParentToChild(point, ancestor, toControl)
	}

	return point
}

// FindControl performs a depth-first search of the controls starting from root, calling test with each visited control.
// If test returns true then the search is stopped and FindControl returns the Control passed to test.
// If no call to test returns true then FindControl returns nil.
func FindControl(root Parent, test func(Control) (found bool)) Control {
	if c, ok := root.(Control); ok && test(c) {
		return c
	}

	for _, child := range root.Children() {
		if test(child.Control) {
			return child.Control
		}

		if parent, ok := child.Control.(Parent); ok {
			if c := FindControl(parent, test); c != nil {
				return c
			}
		}
	}
	return nil
}

func WindowContaining(control Control) *WindowImpl {
	for {
		parent := control.Parent()
		if parent == nil {
			panic("Control's parent was nil")
		}

		switch window := parent.(type) {
		case *WindowImpl:
			return window
		}

		control = parent.(Control)
	}
}

func SetFocus(target Focusable) {
	window := WindowContaining(target)
	window.SetFocus(target)
}

func StringToRuneArray(str string) []rune {
	return bytes.Runes([]byte(str))
}

func RuneArrayToString(arr []rune) string {
	tmp := make([]byte, 8)
	enc := make([]byte, 0, len(arr))
	offset := 0
	for _, r := range arr {
		size := utf8.EncodeRune(tmp, r)
		enc = append(enc, tmp[:size]...)
		offset += size
	}
	return string(enc)
}
