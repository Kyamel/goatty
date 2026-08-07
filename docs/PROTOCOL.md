# The Semantic Layer

Goatty's most consequential design decision is not which renderer it uses. It
is **where the line sits between the terminal and the application.**

The goal is a terminal with a *semantic layer*: output can carry structured
meaning instead of being an opaque stream of ANSI bytes, and the terminal
decides how to present it. The terminal is never responsible for the
application's interface as a whole.

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

## The core idea: a parallel channel

The semantic layer never replaces ANSI. It runs alongside it:

```text
program
   ├── stdout ANSI/text   -> every terminal understands this
   └── semantic objects   -> a capable terminal enriches it
```

A program keeps printing text exactly as it does today, and *additionally* may
emit structured objects. Every object carries a textual fallback, so the output
still works over SSH, through pipes, in log files and in terminals that have
never heard of this protocol.

That constraint is what separates this from an in-product feature. Warp solves
this inside its own application; the point here is an **open protocol** that
any terminal, or any other renderer, can implement.

### Example

A compiler could emit, semantically:

```text
diagnostic
  severity: error
  file:     src/main.go
  line:     42
  message:  undefined: foo
```

The terminal decides the presentation: plain text, a card, something clickable
that opens an editor, errors grouped by file. The compiler never draws a UI and
never simulates one with ANSI.

---

## Anatomy of an object

Every object carries roughly:

```text
type            what kind of thing this is
id              stable handle, so it can be updated later
fallback        the textual rendering, always present
metadata        type-specific fields
actions         what the user can do with it
state           for objects that change over time
relationships   how it connects to other objects
```

`fallback` is not optional. An object without one is a bug: it would make the
program's output depend on the terminal, which is the failure mode this design
exists to avoid.

---

## Semantic scrollback

Today a terminal owns:

```text
line
line
line
line
```

With command boundaries it can instead own:

```text
Command {
    command:     "go test ./..."
    started_at:  ...
    exit_code:   1

    outputs:     [...]
    diagnostics: [...]
    artifacts:   [...]
}
```

Each execution becomes a closed, addressable unit with a beginning, updates and
an end — instead of characters irrevocably dumped into a grid.

The first step here needs no new protocol at all: **OSC 133 shell integration**
already delimits prompt, command, output and exit status. That is the cheapest
possible entry point and it interoperates with terminals that already support
it.

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
canvas
```

A canvas can contain rectangles, text, lines, paths, images and shaders.

If someone wants a Git graph, a pixel editor or a minimap, they use canvas.
There is no need for a "bar chart" protocol, a "pie chart" protocol or a "node
editor" protocol. All of it is just drawing.

Canvas is the one object type that cannot have a meaningful textual fallback.
That is the price of having an escape hatch, and it is the reason canvas should
stay a leaf: applications that reach for it lose portability, and they should
be able to feel that.

---

## Intents should be few

Roughly 15–20 object types, at most. For example:

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
- diagnostic
- artifact
- file-list
- test-suite

Beyond that the set starts getting too specific. A calendar, for instance, does
not get a `calendar` type — it is text, or it is canvas.

### The rule for admitting a type

> **If two independent programs would probably implement this component the
> same way, it deserves to be an object type. Otherwise it belongs to the
> application.**

Applying it:

| Component      | Object type? | Why                                  |
| -------------- | ------------ | ------------------------------------ |
| Progress bar   | yes          | practically everyone builds it alike |
| Spinner        | yes          | same                                 |
| Table          | yes          | same                                 |
| Hyperlink      | yes          | same                                 |
| Diagnostic     | yes          | compilers already agree on the shape |
| File tree      | no           | every editor does it differently     |
| Git explorer   | no           | different everywhere                 |
| Sidebar        | no           | different everywhere                 |

The ones at the bottom stay the application's responsibility.

---

## Immediate-mode feel, retained-mode wire

There is a real tension to resolve here.

The *application-facing API* should feel like an immediate-mode graphics API —
closer to Dear ImGui than to HTML. Few primitives, no layout engine, no
retained widget tree for the application to manage:

```go
BeginFrame()

Text(...)
Image(...)
Progress(...)
Canvas(...)

EndFrame()
```

But the *wire protocol* cannot be immediate-mode. Objects live in the
scrollback, a progress bar updates over time, and a diagnostic emitted an hour
ago is still addressable. That is why objects carry `id` and `state`.

So: immediate-mode ergonomics for the caller, retained and addressable
semantics on the wire. The SDK reconciles the two. Do not let the convenience
of the first leak into the design of the second.

---

## Open question: one protocol or two?

Two things have been discussed that are not obviously the same:

1. **Command and output objects** — diagnostics, tables, artifacts, command
   blocks. The semantics of *what the output means*.
2. **Drawing intents** — the immediate-mode surface above, ending in canvas.
   The semantics of *what to paint*.

They differ on lifetime (persistent vs per-frame), addressing (stable `id` vs
positional), fallback (mandatory vs impossible), and whether they survive a
pipe or an SSH session at all.

**Current recommendation: one protocol, with drawing nested inside it.** The
object protocol is the primary thing; `canvas` is one object type among others.
That buys a single capability negotiation, a single transport and a single
scrollback model, and it keeps drawing as a leaf rather than a parallel
universe.

Two sibling protocols would mean two handshakes, two transports, two fallback
stories, and applications having to choose between them. Merging everything
into drawing would throw away the fallback and every non-terminal renderer,
which is precisely what makes the specification worth writing.

This is not settled. It should be, before step 3 below.

---

## The test

Whenever a new capability is proposed, ask:

> **If I wanted to implement this protocol in another terminal tomorrow, would
> it still look like a terminal?**

If the answer is "no, now I have to implement half a browser", the line has
been crossed.

---

## Beyond terminals

If the producer emits a stream of semantic objects and something else decides
the presentation, the consumer does not have to be a terminal:

```text
producer
    ↓
semantic object stream
    ↓
renderer  ->  terminal | web | native app | TUI | editor
              agent | remote renderer | accessibility layer
```

The same description renders differently in each. This is speculative, but it
is the reason the design should not bake terminal-specific assumptions into the
object model where it can be avoided.

### The specification should not belong to Goatty

Goatty should be the *first implementation*, not the owner. A protocol that
belongs to one terminal becomes one more vendor extension that nobody else
adopts — which is what happened to most terminal protocols historically.

---

## Incremental path

Do not try to design the universal protocol first. In order:

1. **Semantic scrollback** — the terminal stops owning only lines.
2. **Command blocks** — via OSC 133, which already exists.
3. **A few object types** — `diagnostic`, `progress`, `table`, `artifact`.
   Resolve the one-protocol-or-two question before this step.
4. **Standardized actions** — what a user can do with an object.
5. **Updates and chunking** — objects that change after being emitted.
6. **Negotiation and capabilities** — so a program can discover what the
   terminal understands, and degrade to text when it does not.
7. **A specification independent of Goatty.**

---

## What this preserves

Holding this boundary protects the two properties that have kept terminals
relevant for decades:

- the application keeps control of its own behaviour and layout;
- tools stay composable and easy to wire together — shell, pipes, SSH, tmux.

Goatty gains modern capabilities without giving up the simplicity and
hackability that made terminals worth keeping.
