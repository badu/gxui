// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/cgo"
	"github.com/badu/gxui/pkg/math"
	"github.com/badu/gxui/samples/file_dlg/roots"
	"github.com/badu/gxui/samples/flags"
)

var (
	fileColor      = gxui.Color{R: 0.7, G: 0.8, B: 1.0, A: 1}
	directoryColor = gxui.Color{R: 0.8, G: 1.0, B: 0.7, A: 1}
)

// filesAt returns a list of all immediate files in the given directory path.
func filesAt(path string) []string {
	var files []string
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			files = append(files, subpath)
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return files
}

// filesAdapter is an implementation of the gxui.ListAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type filesAdapter struct {
	gxui.AdapterBase
	files []string // The absolute file paths
}

// SetFiles assigns the specified list of absolute-path files to this adapter.
func (a *filesAdapter) SetFiles(files []string) {
	a.files = files
	a.DataChanged(false)
}

func (a *filesAdapter) Count() int {
	return len(a.files)
}

func (a *filesAdapter) ItemAt(index int) gxui.AdapterItem {
	return a.files[index]
}

func (a *filesAdapter) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	for i, f := range a.files {
		if f == path {
			return i
		}
	}
	return -1 // Not found
}

func (a *filesAdapter) Create(driver gxui.Driver, styles *gxui.StyleDefs, index int) gxui.Control {
	path := a.files[index]
	_, name := filepath.Split(path)
	label := gxui.CreateLabel(driver, styles)
	label.SetText(name)
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		label.SetColor(directoryColor)
	} else {
		label.SetColor(fileColor)
	}
	return label
}

func (a *filesAdapter) Size(styles *gxui.StyleDefs) math.Size {
	return math.Size{Width: math.MaxSize.Width, Height: styles.FontSize + 4}
}

// directory implements the gxui.TreeNode interface to represent a directory
// node in a file-system.
type directory struct {
	path    string   // The absolute path of this directory.
	subdirs []string // The absolute paths of all immediate sub-directories.
}

// directoryAt returns a directory structure populated with the immediate
// subdirectories at the given path.
func directoryAt(path string) directory {
	result := directory{path: path}
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			if info.IsDir() {
				result.subdirs = append(result.subdirs, subpath)
				return filepath.SkipDir
			}
		}
		return nil
	})
	return result
}

// Count implements gxui.TreeNodeContainer.
func (d directory) Count() int {
	return len(d.subdirs)
}

// NodeAt implements gxui.TreeNodeContainer.
func (d directory) NodeAt(index int) gxui.TreeNode {
	return directoryAt(d.subdirs[index])
}

// ItemIndex implements gxui.TreeNodeContainer.
func (d directory) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	for i, subpath := range d.subdirs {
		subpath += string(filepath.Separator)
		if strings.HasPrefix(path, subpath) {
			return i
		}
	}
	return -1
}

// Item implements gxui.TreeNode.
func (d directory) Item() gxui.AdapterItem {
	return d.path
}

// Create implements gxui.TreeNode.
func (d directory) Create(driver gxui.Driver, styles *gxui.StyleDefs) gxui.Control {
	_, name := filepath.Split(d.path)
	if name == "" {
		name = d.path
	}
	l := gxui.CreateLabel(driver, styles)
	l.SetText(name)
	l.SetColor(directoryColor)
	return l
}

// directoryAdapter is an implementation of the gxui.TreeAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type directoryAdapter struct {
	directory
	gxui.AdapterBase
}

func (a directoryAdapter) Size(styles *gxui.StyleDefs) math.Size {
	return math.Size{Width: math.MaxSize.Width, Height: styles.FontSize + 4}
}

// Override directory.Create so that the full root is shown, unaltered.
func (a directoryAdapter) Create(driver gxui.Driver, styles *gxui.StyleDefs, index int) gxui.Control {
	l := gxui.CreateLabel(driver, styles)
	l.SetText(a.subdirs[index])
	l.SetColor(directoryColor)
	return l
}

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Open file...")
	window.SetScale(flags.DefaultScaleFactor)

	// fullpath is the textbox at the top of the window holding the current
	// selection's absolute file path.
	fullpath := gxui.CreateTextBox(driver, styles)
	fullpath.SetDesiredWidth(math.MaxSize.Width)

	// directories is the TreeImpl of directories on the left of the window.
	// It uses the directoryAdapter to show the entire system's directory
	// hierarchy.
	directories := gxui.CreateTree(driver, styles)
	directories.SetAdapter(
		&directoryAdapter{
			directory: directory{
				subdirs: roots.Roots(),
			},
		},
	)

	// filesAdapter is the adapter used to show the currently selected directory's
	// content. The adapter has its data changed whenever the selected directory
	// changes.
	adapter := &filesAdapter{}

	// files is the ListImpl of files in the selected directory to the right of the
	// window.
	files := gxui.CreateList(driver, styles)
	files.SetAdapter(adapter)

	open := gxui.CreateButton(driver, styles)
	open.SetText("Open...")
	open.OnClick(func(gxui.MouseEvent) {
		fmt.Printf("File '%s' selected!\n", files.Selected())
		window.Close()
	})

	// If the user hits the enter key while the fullpath control has focus,
	// attempt to select the directory.
	fullpath.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Key == gxui.KeyEnter || ev.Key == gxui.KeyKpEnter {
			path := fullpath.Text()
			if directories.Select(path) {
				directories.Show(path)
			}
		}
	})

	// When the directory selection changes, update the files list
	directories.OnSelectionChanged(func(item gxui.AdapterItem) {
		dir := item.(string)
		adapter.SetFiles(filesAt(dir))
		fullpath.SetText(dir)
	})

	// When the file selection changes, update the fullpath text
	files.OnSelectionChanged(func(item gxui.AdapterItem) {
		fullpath.SetText(item.(string))
	})

	// When the user double-clicks a directory in the file list, select it in the
	// directories tree view.
	files.OnDoubleClick(func(gxui.MouseEvent) {
		if path, ok := files.Selected().(string); ok {
			if fi, err := os.Stat(path); err == nil && fi.IsDir() {
				if directories.Select(path) {
					directories.Show(path)
				}
			} else {
				fmt.Printf("File '%s' selected!\n", path)
				window.Close()
			}
		}
	})

	// Start with the CWD selected and visible.
	if cwd, err := os.Getwd(); err == nil {
		if directories.Select(cwd) {
			directories.Show(directories.Selected())
		}
	}

	splitter := gxui.CreateSplitterLayout(driver, styles)
	splitter.SetOrientation(gxui.Horizontal)
	splitter.AddChild(directories)
	splitter.AddChild(files)

	topLayout := gxui.CreateLinearLayout(driver, styles)
	topLayout.SetDirection(gxui.TopToBottom)
	topLayout.AddChild(fullpath)
	topLayout.AddChild(splitter)

	btmLayout := gxui.CreateLinearLayout(driver, styles)
	btmLayout.SetDirection(gxui.BottomToTop)
	btmLayout.SetHorizontalAlignment(gxui.AlignRight)
	btmLayout.AddChild(open)
	btmLayout.AddChild(topLayout)

	window.AddChild(btmLayout)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{Left: 10, Top: 10, Right: 10, Bottom: 10})
}

func main() {
	cgo.StartDriver(appMain)
}
