# Goatty Vision

> A simple, hackable and modern terminal emulator platform written in Go.

This document describes what Goatty is trying to become and the principles that
decisions should be measured against. Two companion documents cover the parts
that need more detail:

- [ROADMAP.md](ROADMAP.md) -- the phased plan for getting there.
- [PROTOCOL.md](PROTOCOL.md) -- the semantic layer: structured objects alongside
  ANSI, and where the line sits between the terminal and the application. This
  is the single most consequential design decision here, and the one genuinely
  novel thing the project is attempting.

---

## Philosophy

Goatty is **not** just another terminal emulator.

Its goal is to become a **terminal platform**: a modular ecosystem where the
emulator itself is only one component.

The project prioritizes:

- Simplicity over cleverness.
- Readability over micro-optimizations.
- Extensibility over hardcoded features.
- Small, composable packages.
- A contributor experience where understanding the codebase takes days, not
  months.

Performance matters, but maintainability comes first. The project should avoid
architectural decisions that make contributors afraid to touch the code.

---

## Long-term goals

Goatty should provide:

- a modern GPU renderer
- excellent Unicode support
- high compatibility with the xterm and kitty protocols
- image protocol support
- configurable appearance
- customizable shaders
- native-feeling desktop integration
- an extension ecosystem
- embeddable terminal widgets for Go applications

Eventually, Goatty should become something similar to:

> "The Ebitengine of terminal emulators."

Not just an application, but a reusable platform.

---

## Core Principles

### 1. Everything should be replaceable

The renderer should not know about the parser.
The parser should not know about the renderer.
The renderer should not know about the PTY.
The PTY should not know about the UI.

Every subsystem communicates through small interfaces.

### 2. Small packages

Instead of:

```text
internal/
    terminal/
```

prefer:

```text
cell/
grid/
parser/
pty/
renderer/
font/
unicode/
selection/
clipboard/
input/
config/
```

Each package should have one clear responsibility.

### 3. No giant files

Avoid 3000-line files. Prefer many files of roughly 100-300 lines. Large files
become psychological barriers for contributors.

### 4. Data-oriented architecture

The terminal is fundamentally a grid of cells, a set of state transitions, and
drawing. Keep rendering data contiguous and simple; avoid deeply nested object
graphs.

### 5. Minimal dependencies

The project should avoid accumulating frameworks. A dependency should exist
only when it clearly reduces maintenance.

### 6. Go idioms

Prefer interfaces, composition and explicit code. Avoid abstractions borrowed
from Java or C++ that Go does not need.

---

## Architecture

```text
            +-----------------------+
            |      Application      |
            +-----------+-----------+
                        |
            +-----------v-----------+
            |        Window         |
            +-----------+-----------+
                        |
        +---------------+----------------+
        |                                |
+-------v------+                +--------v--------+
|     UI       |                | Terminal View   |
+--------------+                +--------+--------+
                                         |
                            +------------v------------+
                            |      Renderer API       |
                            +------------+------------+
                                         |
                           +-------------v--------------+
                           | GPU Backend (Ebitengine)   |
                           +----------------------------+

                    Terminal Core

+------------------------------------------------------+
|                    Screen Model                      |
+------------------------------------------------------+
| Parser | Grid | PTY | Scrollback | Selection | Input |
+------------------------------------------------------+
```

The renderer only ever receives immutable render data. It never owns terminal
state.

---

## Why Darktile?

Darktile already proves that a GPU terminal in Go is viable. What it does not
do is separate concerns that should eventually become reusable libraries.

Rather than continuing as a monolithic application, Goatty evolves that
codebase into a modular platform.

---

## Non-goals

Goatty should avoid becoming:

- an IDE
- a shell
- a window manager
- a scripting language
- a plugin host with unrestricted in-process execution
- a UI toolkit or a browser

The last one is the easiest to drift into and the hardest to reverse.
[PROTOCOL.md](PROTOCOL.md) exists to keep that line visible.

---

## Deferred: the transactional terminal

A separate and more radical idea, deliberately parked rather than rejected.

It treats a command not as `command -> stdout/stderr -> exit code` but as a
**transaction over the environment**:

```text
Transaction #42
command: make install

changes:
  files:      + /usr/local/bin/foo
              ~ ~/.config/foo/config.toml
  processes:  + pid 12345
  env:        PATH changed
  network:    GET https://...
```

From which actions follow: undo, rerun, inspect changes, diff the filesystem
before and after. Roughly *shell + audit log + light sandbox + undo*.

This is a different axis from [PROTOCOL.md](PROTOCOL.md):

```text
semantic objects      -> the meaning of the output
transactional model   -> the meaning of the execution's effects
```

**Why it is deferred:** implementing it generally means sandboxing, namespaces,
filesystem overlays or snapshots, syscall tracing and child-process
accounting -- deep OS integration, with a different answer on every platform.
The semantic object protocol is incrementally implementable and just as novel,
so it goes first.

Revisit once the object protocol has shipped and command blocks exist, since a
transaction is a natural extension of a command block.

---

## Contributor Experience

A new contributor should be able to:

- clone the repository
- understand the architecture in one afternoon
- modify one subsystem without learning every other subsystem
- run tests quickly
- build without complicated setup

Good documentation is a feature, not overhead.

---

## Success Metrics

The project is succeeding if contributors say:

- "I found the code surprisingly easy to understand."
- "Adding a feature didn't require touching unrelated packages."
- "The architecture guided me toward the correct implementation."
- "I could replace a subsystem without rewriting the entire application."

Those outcomes are worth more than squeezing a few extra frames per second out
of the renderer.

---

## Vision Statement

Goatty aims to become the most approachable modern terminal emulator codebase.

Not necessarily the fastest. Not necessarily the most feature-rich. But the one
developers genuinely enjoy reading, extending, and building upon.
