package purego

import (
	"errors"
	"fmt"
	"strings"
)

func CreateProgram(ctx *Functions, vsSrc, fsSrc string, attribs []string) (uint32, error) {
	vs, err := CreateShader(ctx, VERTEX_SHADER, vsSrc)
	if err != nil {
		return 0, err
	}
	defer ctx.DeleteShader(vs)
	fs, err := CreateShader(ctx, FRAGMENT_SHADER, fsSrc)
	if err != nil {
		return 0, err
	}
	defer ctx.DeleteShader(fs)
	prog := ctx.CreateProgram()
	if prog == 0 {
		return 0, errors.New("glCreateProgram failed")
	}
	ctx.AttachShader(prog, vs)
	ctx.AttachShader(prog, fs)
	for i, a := range attribs {
		ctx.BindAttribLocation(prog, Attrib(i), a)
	}
	ctx.LinkProgram(prog)
	if ctx.GetProgrami(prog, LINK_STATUS) == 0 {
		log := ctx.GetProgramInfoLog(prog)
		ctx.DeleteProgram(prog)
		return 0, fmt.Errorf("program link failed: %s", strings.TrimSpace(log))
	}
	return prog, nil
}

func CreateComputeProgram(ctx *Functions, src string) (uint32, error) {
	cs, err := CreateShader(ctx, COMPUTE_SHADER, src)
	if err != nil {
		return 0, err
	}
	defer ctx.DeleteShader(cs)
	prog := ctx.CreateProgram()
	if prog == 0 {
		return 0, errors.New("glCreateProgram failed")
	}
	ctx.AttachShader(prog, cs)
	ctx.LinkProgram(prog)
	if ctx.GetProgrami(prog, LINK_STATUS) == 0 {
		log := ctx.GetProgramInfoLog(prog)
		ctx.DeleteProgram(prog)
		return 0, fmt.Errorf("program link failed: %s", strings.TrimSpace(log))
	}
	return prog, nil
}

func CreateShader(ctx *Functions, typ Enum, src string) (uint32, error) {
	sh := ctx.CreateShader(typ)
	if sh == 0 {
		return 0, errors.New("glCreateShader failed")
	}
	ctx.ShaderSource(sh, src)
	ctx.CompileShader(sh)
	if ctx.GetShaderi(sh, COMPILE_STATUS) == 0 {
		log := ctx.GetShaderInfoLog(sh)
		ctx.DeleteShader(sh)
		return 0, fmt.Errorf("shader compilation failed: %s", strings.TrimSpace(log))
	}
	return sh, nil
}

func ParseGLVersion(glVer string) (version [2]int, gles bool, err error) {
	var ver [2]int
	if _, err := fmt.Sscanf(glVer, "OpenGL ES %d.%d", &ver[0], &ver[1]); err == nil {
		return ver, true, nil
	} else if _, err := fmt.Sscanf(glVer, "WebGL %d.%d", &ver[0], &ver[1]); err == nil {
		// WebGL major version v corresponds to OpenGL ES version v + 1
		ver[0]++
		return ver, true, nil
	} else if _, err := fmt.Sscanf(glVer, "%d.%d", &ver[0], &ver[1]); err == nil {
		return ver, false, nil
	}
	return ver, false, fmt.Errorf("failed to parse OpenGL ES version (%s)", glVer)
}
