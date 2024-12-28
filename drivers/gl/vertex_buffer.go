// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "fmt"

type vertexBuffer struct {
	streams map[string]*vertexStream
	count   int
}

func newVertexBuffer(streams ...*vertexStream) *vertexBuffer {
	vb := &vertexBuffer{streams: make(map[string]*vertexStream)}
	for index, stream := range streams {
		if index == 0 {
			vb.count = stream.count
		} else {
			if vb.count != stream.count {
				panic(fmt.Errorf("inconsistent vertex count in vertex buffer. %s has %d vertices, %s has %d", streams[index-1].name, streams[index-1].count, stream.name, stream.count))
			}
		}
		vb.streams[stream.name] = stream
	}
	return vb
}
