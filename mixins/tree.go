// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type TreeOuter interface {
	ListOuter
	PaintUnexpandedSelection(c gxui.Canvas, r math.Rect)
}

type Tree struct {
	List
	FocusablePart
	outer       TreeOuter
	treeAdapter gxui.TreeAdapter
	listAdapter *TreeToListAdapter
	creator     TreeControlCreator
}

func (t *Tree) Init(outer TreeOuter, theme gxui.Theme) {
	t.List.Init(outer, theme)
	t.FocusablePart.Init(outer)
	t.outer = outer
	t.creator = defaultTreeControlCreator{}
}

func (t *Tree) SetControlCreator(control TreeControlCreator) {
	t.creator = control
	if t.treeAdapter != nil {
		t.listAdapter = CreateTreeToListAdapter(t.treeAdapter, t.creator)
		t.DataReplaced()
	}
}

// gxui.Tree complaince
func (t *Tree) SetAdapter(adapter gxui.TreeAdapter) {
	if t.treeAdapter == adapter {
		return
	}

	if adapter != nil {
		t.treeAdapter = adapter
		t.listAdapter = CreateTreeToListAdapter(adapter, t.creator)
		t.List.SetAdapter(t.listAdapter)
	} else {
		t.listAdapter = nil
		t.treeAdapter = nil
		t.List.SetAdapter(nil)
	}
}

func (t *Tree) Adapter() gxui.TreeAdapter {
	return t.treeAdapter
}

func (t *Tree) Show(item gxui.AdapterItem) {
	t.listAdapter.ExpandItem(item)
	t.List.ScrollTo(item)
}

func (t *Tree) ContainsItem(item gxui.AdapterItem) bool {
	return t.listAdapter != nil && t.listAdapter.Contains(item)
}

func (t *Tree) ExpandAll() {
	t.listAdapter.ExpandAll()
}

func (t *Tree) CollapseAll() {
	t.listAdapter.CollapseAll()
}

func (t *Tree) PaintUnexpandedSelection(canvas gxui.Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, gxui.CreatePen(1, gxui.Gray50), gxui.TransparentBrush)
}

// List override
func (t *Tree) PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) {
	t.List.PaintChild(canvas, child, idx)
	if t.selectedItem != nil {
		if deepest := t.listAdapter.DeepestNode(t.selectedItem); deepest != nil {
			if item := deepest.Item(); item != t.selectedItem {
				// The selected item is hidden by an unexpanded node.
				// Highlight the deepest visible node instead.
				if details, found := t.details[item]; found {
					if child == details.child {
						b := child.Bounds().Expand(child.Control.Margin())
						t.outer.PaintUnexpandedSelection(canvas, b)
					}
				}
			}
		}
	}
}

// InputEventHandlerPart override
func (t *Tree) KeyPress(event gxui.KeyboardEvent) bool {
	switch event.Key {
	case gxui.KeyLeft:
		if item := t.Selected(); item != nil {
			node := t.listAdapter.DeepestNode(item)
			if node.Collapse() {
				return true
			}
			if p := node.Parent(); p != nil {
				return t.Select(p.Item())
			}
		}

	case gxui.KeyRight:
		if item := t.Selected(); item != nil {
			node := t.listAdapter.DeepestNode(item)
			if node.Expand() {
				return true
			}
		}
	}

	return t.List.KeyPress(event)
}

type defaultTreeControlCreator struct{}

func (defaultTreeControlCreator) Create(theme gxui.Theme, control gxui.Control, node *TreeToListNode) gxui.Control {
	ll := theme.CreateLinearLayout()
	ll.SetDirection(gxui.LeftToRight)

	btn := theme.CreateButton()
	btn.SetBackgroundBrush(gxui.TransparentBrush)
	btn.SetBorderPen(gxui.CreatePen(1, gxui.Gray30))
	btn.SetMargin(math.Spacing{L: 2, R: 2, T: 1, B: 1})
	btn.OnClick(func(ev gxui.MouseEvent) {
		if ev.Button == gxui.MouseButtonLeft {
			node.ToggleExpanded()
		}
	})

	update := func() {
		btn.SetVisible(!node.IsLeaf())
		if node.IsExpanded() {
			btn.SetText("-")
		} else {
			btn.SetText("+")
		}
	}
	update()

	gxui.WhileAttached(btn, node.OnChange, update)

	ll.AddChild(btn)
	ll.AddChild(control)
	ll.SetPadding(math.Spacing{L: 16 * node.Depth()})
	return ll
}

func (defaultTreeControlCreator) Size(theme gxui.Theme, treeControlSize math.Size) math.Size {
	return treeControlSize
}
