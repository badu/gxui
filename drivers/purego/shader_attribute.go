package purego

type shaderAttribute struct {
	name       string
	size       int32
	shaderType shaderDataType
	glAttr     uint32
}

func (a shaderAttribute) enableArray(fn *Functions) {
	fn.EnableVertexAttribArray(a.glAttr)
}

func (a shaderAttribute) disableArray(fn *Functions) {
	fn.DisableVertexAttribArray(a.glAttr)
}

func (a shaderAttribute) attribPointer(fn *Functions, size int32, ty uint32, normalized bool, stride int32, offset int) {
	fn.VertexAttribPointer(a.glAttr, size, ty, normalized, stride, offset)
}
