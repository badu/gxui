package purego

type shaderAttribute struct {
	fn         *Functions
	name       string
	size       int
	shaderType shaderDataType
	glAttr     Attrib
}

func (a shaderAttribute) enableArray() {
	a.fn.EnableVertexAttribArray(a.glAttr)
}

func (a shaderAttribute) disableArray() {
	a.fn.DisableVertexAttribArray(a.glAttr)
}

func (a shaderAttribute) attribPointer(size int32, ty uint32, normalized bool, stride int32, offset int) {
	a.fn.VertexAttribPointer(a.glAttr, int(size), Enum(ty), normalized, int(stride), offset)
}
