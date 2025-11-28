26.11.2025
---
TODO : purego Windows driver requires testing

- BlendFunc
- GetActiveUniform
- GetActiveAttrib
- GetAttribLocation
- UniformMatrix2fv
- UniformMatrix3fv
- UniformMatrix4fv
- Uniform1fv
- Uniform2fv
- Uniform3fv
- Uniform4fv

25.11.2025
---
github.com/goxjs/gl unique calls :
gl.ActiveTexture
gl.AttachShader
gl.BindBuffer
gl.BindTexture
gl.BlendFunc
gl.BufferData
gl.Clear
gl.ClearColor
gl.CompileShader
gl.CreateBuffer
gl.CreateProgram
gl.CreateShader
gl.CreateTexture
gl.DeleteBuffer
gl.DeleteProgram
gl.DeleteTexture
gl.Disable
gl.DisableVertexAttribArray
gl.DrawArrays
gl.DrawElements
gl.Enable
gl.EnableVertexAttribArray
gl.Enum
gl.GetActiveAttrib
gl.GetActiveUniform
gl.GetAttribLocation
gl.GetError
gl.GetProgramInfoLog
gl.GetProgrami
gl.GetShaderInfoLog
gl.GetShaderi
gl.GetUniformLocation
gl.LinkProgram
gl.Scissor
gl.ShaderSource
gl.TexImage2D
gl.TexParameteri
gl.Uniform1f
gl.Uniform1fv
gl.Uniform1i
gl.Uniform2fv
gl.Uniform3fv
gl.Uniform4fv
gl.UniformMatrix2fv
gl.UniformMatrix3fv
gl.UniformMatrix4fv
gl.UseProgram
gl.VertexAttribPointer
gl.Viewport

28.12.2024
---

First thing I needed to do is to add `go.mod` file. Of course, since all imports were pointing to [
`gxui`](https://github.com/google/gxui), I had to replace to match the [`go.mod`](https://github.com/badu/gxui).

Running first example, triggered the first error, since `GetClipboardString` return has changed in `glfw`.

I've played around with the code (in the same time with `fyne` and `gio`) and I can congratulate Ben Clayton for the
simplicity and clarity of the code. What I am really impressed is the cascading of the `Init()` functions, which allows
inheritance in the same time with
ability to compose.

By looking of modifications that other developers brought to their forks, I can totally say that it is easy to break the
code, mostly when you try to mess with the way `driver` was designed.

Additions:

[x] font size is now settable through the `flags` (in `samples` folder). Currently, is default to 24. It is exposed to
be accessed by descendants of `AdapterBase`.

[x] Light and dark themes, are setting display width and height, so `samples` resize the window to right half of the
monitor.

[x] Increasing readability : renaming parameters, properties, returns and eliminate named returns

30.12.2024
---

[x] Moved `mixins` into main package, because I am planning to reduce the number of returned interfaces.

[x] Promoted `theme` as `DefaultApp`, trying to break the interface dependency into a direct dependency.

[x] `CreateBubbleOverlay`, `CreateButton`, `CreateCodeEditor`, `CreateDropDownList`, `CreateImage`, `CreateLabel`,
`CreateLinearLayout`, `CreateList`, `CreatePanelHolder`, `CreateProgressBar`, `CreateScrollBar`, `CreateScrollLayout`,
`CreateSplitterLayout`, `CreateTableLayout`, `CreateTextBox`, `CreateTree`, `CreateWindow` should NOT depend on the
application itself, but on the `Driver` interface and a `StyleManager` (which currently doesn't exist).

The first thing that an application is building is a `Window`, which holds focus, keyboard and mouse managers.

[x] Application interface was removed

31.12.2024
---

Creating such a framework requires "only" this much (taken, shameless, from a CSS post):

1. Box Model
   • width
   • height
   • margin
   • padding
   • border
   • box-sizing

2. Positioning
   • position
   • top
   • right
   • bottom
   • left
   • float
   • clear
   • z-index

3. Typography
   • font-family
   • font-size
   • font-weight
   • font-style
   • color
   • line-height
   • letter-spacing
   • text-align
   • text-decoration
   • text-transform

4. Visual Formatting
   • background-color
   • background-image
   • background-repeat
   • background-position
   • background-size
   • color
   • display
   • visibility
   • overflow
   • opacity

5. Flexbox
   • display: flex;
   • flex-direction
   • justify-content
   • align-items
   • align-self
   • flex

6. Grid Layout
   • display: grid;
   • grid-template-columns
   • grid-template-rows
   • grid-column
   • grid-row
   • grid-gap
   • grid-area

7. Transforms and Animations
   • transform
   • transition
   • animation

8. Others
   • cursor
   • list-style
   • outline
   • user-select
   • pointer-events