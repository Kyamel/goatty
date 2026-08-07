# The Semantic Layer

Goatty's most consequential design decision is not which renderer it uses. It
is **where the line sits between the terminal and the application.**

The goal is a terminal with a *semantic layer*: it understands that something
is a progress bar, an image, a table or a canvas, but it is never responsible
for the application's interface as a whole.

---

## Why the line matters

A terminal is valuable today because it is, essentially, a surface an
application draws on. That is why all of this works without anyone
coordinating:

```sh
cat log.txt
grep foo
less
vim
```

Start adding flexbox, constraint solving, rich mouse events, a focus tree,
widgets and declarative animation, and you have stopped building a terminal.
You are building a browser, or a UI toolkit.

The line should not be drawn based on how complex an object is. It should be
drawn based on **who controls the layout.**

---

## Where the line sits

The protocol **may** create:

- text
- images
- canvas
- tables
- progress bars
- hyperlinks
- notifications
- scrollable regions
- overlays
- simple popups

The protocol **may not** create:

- buttons
- menus
- checkboxes
- forms
- widget trees
- flexbox or grid layout
- focus management
- a rich event system

The reason is not aesthetic. Every item in the second list ties the application
to Goatty specifically. The first list does not.

---

## Canvas is the escape hatch

Rather than adding a new object for every case, there is one generic object:

```text
Canvas
```

A canvas can contain rectangles, text, lines, paths, images and shaders.

If someone wants a Git graph, a pixel editor or a minimap, they use Canvas.
There is no need for a "bar chart" protocol, a "pie chart" protocol or a "node
editor" protocol. All of it is just drawing.

---

## Intents should be few

Roughly 15–20 intents, at most. For example:

- prompt
- progress
- spinner
- error
- warning
- success
- notification
- table
- image
- code
- markdown
- diff

Beyond that the set starts getting too specific. A calendar, for instance, does
not get a `calendar` intent — it is text, or it is Canvas.

### The rule for admitting an intent

> **If two independent programs would probably implement this component the
> same way, it deserves to be an intent. Otherwise it belongs to the
> application.**

Applying it:

| Component      | Intent? | Why                                  |
| -------------- | ------- | ------------------------------------ |
| Progress bar   | yes     | practically everyone builds it alike |
| Spinner        | yes     | same                                 |
| Table          | yes     | same                                 |
| Hyperlink      | yes     | same                                 |
| File tree      | no      | every editor does it differently     |
| Git explorer   | no      | different everywhere                 |
| Sidebar        | no      | different everywhere                 |

The three at the bottom stay the application's responsibility.

---

## Shape of the API

The protocol should feel more like an immediate-mode graphics API — closer to
Dear ImGui than to HTML:

```go
BeginFrame()

Text(...)
Image(...)
Progress(...)
Canvas(...)

EndFrame()
```

The terminal renders. The application keeps owning its logic.

---

## The test

Whenever a new capability is proposed, ask:

> **If I wanted to implement this protocol in another terminal tomorrow, would
> it still look like a terminal?**

If the answer is "no, now I have to implement half a browser", the line has
been crossed.

---

## What this preserves

Holding this boundary protects the two properties that have kept terminals
relevant for decades:

- the application keeps control of its own behaviour and layout;
- tools stay composable and easy to wire together — shell, pipes, SSH, tmux.

Goatty gains modern capabilities without giving up the simplicity and
hackability that made terminals worth keeping.
