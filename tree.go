// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

// Tree is the interface of all controls that visualize a hierarchical tree
// structure of items.
type Tree interface {
	Focusable

	// SetAdapter binds the specified TreeAdapter to this Tree control, replacing
	// any previously bound adapter.
	SetAdapter(TreeAdapter)

	// TreeAdapter returns the currently bound adapter.
	Adapter() TreeAdapter

	// Show makes the specified item visible, expanding the tree if necessary.
	Show(AdapterItem)

	// ExpandAll expands all tree nodes.
	ExpandAll()

	// CollapseAll collapses all tree nodes.
	CollapseAll()

	// Selected returns the currently selected item.
	Selected() AdapterItem

	// Select makes the specified item selected. The tree will not automatically
	// expand to the newly selected item. If the Tree does not contain the
	// specified item, then Select returns false and the previous selection
	// remains unaltered.
	Select(AdapterItem) bool

	// OnSelectionChanged registers the function f to be called when the selection
	// changes.
	OnSelectionChanged(f func(AdapterItem)) EventSubscription
}

// TreeNodeContainer is the interface used by nodes that can hold sub-nodes in the tree.
type TreeNodeContainer interface {
	// Count returns the number of immediate child nodes.
	Count() int

	// Node returns the i'th child TreeNode.
	NodeAt(i int) TreeNode

	// ItemIndex returns the index of the child equal to item, or the index of the
	// child that indirectly contains item, or if the item is not found under this
	// node, -1.
	ItemIndex(item AdapterItem) int
}

// TreeNode is the interface used by nodes in the tree.
type TreeNode interface {
	TreeNodeContainer

	// Item returns the AdapterItem this node.
	// It is important for the TreeNode to return consistent AdapterItems for
	// the same data, so that selections can be persisted, or re-ordering
	// animations can be played when the dataset changes.
	// The AdapterItem returned must be equality-unique across the entire Adapter.
	Item() AdapterItem

	// Create returns a Control visualizing this node.
	Create(theme App) Control
}

// TreeAdapter is an interface used to visualize a set of hierarchical items.
// Users of the TreeAdapter should presume the data is unchanged until the
// OnDataChanged or OnDataReplaced events are fired.
type TreeAdapter interface {
	TreeNodeContainer

	// Size returns the size that each of the item's controls will be displayed
	// at for the given theme.
	Size(App) math.Size

	// OnDataChanged registers f to be called when there is a partial change in
	// the items of the adapter. Scroll positions and selections should be
	// preserved if possible.
	// If recreateControls is true then each of the visible controls should be
	// recreated by re-calling Create().
	OnDataChanged(f func(recreateControls bool)) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete
	// replacement of items in the adapter.
	OnDataReplaced(f func()) EventSubscription
}

type TreeOuter interface {
	ListOuter
	PaintUnexpandedSelection(c Canvas, r math.Rect)
}

type TreeImpl struct {
	ListImpl
	FocusablePart
	outer       TreeOuter
	treeAdapter TreeAdapter
	listAdapter *TreeToListAdapter
	creator     TreeControlCreator
}

func (t *TreeImpl) Init(outer TreeOuter, theme App) {
	t.ListImpl.Init(outer, theme)
	t.FocusablePart.Init()
	t.outer = outer
	t.creator = defaultTreeControlCreator{}
}

func (t *TreeImpl) SetControlCreator(control TreeControlCreator) {
	t.creator = control
	if t.treeAdapter != nil {
		t.listAdapter = CreateTreeToListAdapter(t.treeAdapter, t.creator)
		t.DataReplaced()
	}
}

// gxui.Tree complaince
func (t *TreeImpl) SetAdapter(adapter TreeAdapter) {
	if t.treeAdapter == adapter {
		return
	}

	if adapter != nil {
		t.treeAdapter = adapter
		t.listAdapter = CreateTreeToListAdapter(adapter, t.creator)
		t.ListImpl.SetAdapter(t.listAdapter)
	} else {
		t.listAdapter = nil
		t.treeAdapter = nil
		t.ListImpl.SetAdapter(nil)
	}
}

func (t *TreeImpl) Adapter() TreeAdapter {
	return t.treeAdapter
}

func (t *TreeImpl) Show(item AdapterItem) {
	t.listAdapter.ExpandItem(item)
	t.ListImpl.ScrollTo(item)
}

func (t *TreeImpl) ContainsItem(item AdapterItem) bool {
	return t.listAdapter != nil && t.listAdapter.Contains(item)
}

func (t *TreeImpl) ExpandAll() {
	t.listAdapter.ExpandAll()
}

func (t *TreeImpl) CollapseAll() {
	t.listAdapter.CollapseAll()
}

func (t *TreeImpl) PaintUnexpandedSelection(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, CreatePen(1, Gray50), TransparentBrush)
}

// ListImpl override
func (t *TreeImpl) PaintChild(canvas Canvas, child *Child, idx int) {
	t.ListImpl.PaintChild(canvas, child, idx)
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
func (t *TreeImpl) KeyPress(event KeyboardEvent) bool {
	switch event.Key {
	case KeyLeft:
		if item := t.Selected(); item != nil {
			node := t.listAdapter.DeepestNode(item)
			if node.Collapse() {
				return true
			}
			if p := node.Parent(); p != nil {
				return t.Select(p.Item())
			}
		}

	case KeyRight:
		if item := t.Selected(); item != nil {
			node := t.listAdapter.DeepestNode(item)
			if node.Expand() {
				return true
			}
		}
	}

	return t.ListImpl.KeyPress(event)
}

type defaultTreeControlCreator struct{}

func (defaultTreeControlCreator) Create(theme App, control Control, node *TreeToListNode) Control {
	ll := theme.CreateLinearLayout()
	ll.SetDirection(LeftToRight)

	btn := theme.CreateButton()
	btn.SetBackgroundBrush(TransparentBrush)
	btn.SetBorderPen(CreatePen(1, Gray30))
	btn.SetMargin(math.Spacing{L: 2, R: 2, T: 1, B: 1})
	btn.OnClick(func(ev MouseEvent) {
		if ev.Button == MouseButtonLeft {
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

	WhileAttached(btn, node.OnChange, update)

	ll.AddChild(btn)
	ll.AddChild(control)
	ll.SetPadding(math.Spacing{L: 16 * node.Depth()})
	return ll
}

func (defaultTreeControlCreator) Size(theme App, treeControlSize math.Size) math.Size {
	return treeControlSize
}
