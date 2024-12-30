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

[ ] `CreateBubbleOverlay`, `CreateButton`, `CreateCodeEditor`, `CreateDropDownList`, `CreateImage`, `CreateLabel`,
`CreateLinearLayout`, `CreateList`, `CreatePanelHolder`, `CreateProgressBar`, `CreateScrollBar`, `CreateScrollLayout`,
`CreateSplitterLayout`, `CreateTableLayout`, `CreateTextBox`, `CreateTree`, `CreateWindow` should NOT depend on the
application itself, but on the `Driver` interface and a `StyleManager` (which currently doesn't exist). 

The first thing that an application is building is a `Window`, which holds focus, keyboard and mouse managers.