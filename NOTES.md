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