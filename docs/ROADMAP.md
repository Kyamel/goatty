# Refactoring Roadmap

The plan for turning the current codebase into the platform described in
[VISION.md](VISION.md). Phases are ordered by dependency, not by priority --
each one is mostly meaningless until the previous one has landed.

This document is expected to change as phases complete.

The semantic layer has a separate, independent path of its own -- see the
incremental steps in [PROTOCOL.md](PROTOCOL.md). It is not blocked on these
phases, but it gets much cheaper after Phase 2 separates the terminal core.

---

## Phase 1 -- Understand the current architecture

Before moving anything, map what is actually there.

Goals:

- document every subsystem
- identify responsibilities
- identify coupling
- identify global state
- identify hidden assumptions

Deliverables:

- architecture diagrams
- dependency graph
- ownership map

---

## Phase 2 -- Separate the terminal core

Extract, with no dependency on rendering:

- screen model
- parser
- cell representation
- grid
- cursor
- scrollback

---

## Phase 3 -- Define stable interfaces

No package should depend on a concrete implementation. For example:

```go
type Renderer interface {
    Draw(Frame)
}

type PTY interface {
    Read([]byte)
    Write([]byte)
}

type Clipboard interface {
    Copy(string)
    Paste() string
}
```

---

## Phase 4 -- Introduce a rendering abstraction

The terminal core produces a `RenderFrame` containing:

- glyph runs
- background runs
- cursor
- decorations
- images

The renderer only consumes that structure.

---

## Phase 5 -- Replace the renderer

Once Phase 4 holds, the backend becomes replaceable: Ebitengine, OpenGL,
Vulkan, WebGPU, or a headless renderer.

A headless renderer is worth building early -- it is what makes the core
testable without a GPU.

---

## Phase 6 -- UI layer

The terminal stays independent. The application adds tabs, a command palette,
settings, notifications and dialogs on top of it. This is what keeps the core
embeddable.

---

## Phase 7 -- Configuration

Configuration stays declarative, in TOML. Configuration is static; behaviour
belongs to extensions.

---

## Phase 8 -- Extension system

Rather than embedding a scripting language deep in the emulator, expose a
stable RPC interface and let extensions be independent processes.

That buys crash isolation, language independence, easy SDKs, a security
boundary, and far easier debugging. SDKs could follow for Go, Rust, Python and
JavaScript.

---

## Phase 9 -- Public Go packages

The repository ends up exposing reusable modules:

```text
goatty/parser
goatty/grid
goatty/renderer
goatty/config
goatty/pty
goatty/image
goatty/protocol
```

Applications should be able to embed Goatty as a library.
