# Goatty

<p align="center">
  <strong>Go Open Architecture & Toolkit for TTYs</strong><br>
  <em>A modular, hackable terminal platform built in Go.</em>
</p>

<p align="center">
  <a href="https://github.com/kyamel/goatty/actions/workflows/test.yml"><img alt="Test status" src="https://github.com/kyamel/goatty/actions/workflows/test.yml/badge.svg"></a>
  <a href="https://github.com/kyamel/goatty/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/v/release/kyamel/goatty"></a>
  <a href="LICENSE"><img alt="MIT license" src="https://img.shields.io/github/license/kyamel/goatty"></a>
</p>

![Goatty terminal demo](demo.gif)

Goatty is not trying to become the terminal with the longest feature list.
It is trying to become the terminal codebase people actually want to read,
embed, extend and build upon.

The terminal emulator is the first application. The long-term product is the
platform underneath it: a reusable terminal core, replaceable renderers, small
Go packages and an open semantic protocol that can carry meaning without
breaking the TTY ecosystem.

> **Project status:** Goatty already runs as a GPU-rendered terminal emulator.
> The public toolkit and semantic protocol are active architecture work, not
> finished APIs. See the [roadmap](docs/ROADMAP.md) for the path from here.

## Why Goatty?

Terminals survived for decades because they keep a remarkably useful boundary:
applications own their behaviour, while the terminal renders a portable stream.

Goatty tries to modernize that boundary without erasing it.

**GO**ATTY - **Goatty is Go without ceremony.** Readable enough to enter. Small
enough to understand. Direct enough to change.

G**OA**TTY - **Goatty is an Open Architecture.** Nothing sacred. Nothing sealed.
Every boundary is an invitation to replace, extend and experiment.

G**OA**TTY - **Goatty is where terminals become Object-Aware.** Applications
do not merely print bytes; they express intent. Diagnostics, progress, tables
and artifacts arrive with their meaning intact.

GOA**TT**Y - **Goatty is a Terminal Toolkit.** Not a monolith you merely run,
but a foundation you can build on.

GOA**TTY** - **Goatty is still a TTY.** It moves forward without abandoning
pipes, shells, SSH or decades of Unix composition.

**GOAT**TY - **Goatty is the GOAT.** Greatest Of All Time is more than a joke;
it is the ambition. Built to climb :goat: beyond what terminals are today
without losing the ground that made them timeless.

> **Go Open Architecture & Toolkit for TTYs.**
>
> Not just another terminal. A place to build what terminals become next.

## A semantic layer, not a browser

ANSI output is a stream of drawing instructions. Goatty's proposed semantic
layer runs beside that stream and describes what output *means*:

```text
program
   |-- ANSI/text         -> works everywhere
   `-- semantic objects  -> richer presentation when supported
```

A compiler could emit a diagnostic object like this:

```text
diagnostic
  severity: error
  file:     internal/parser.go
  line:     42
  message:  undefined: token
  fallback: "internal/parser.go:42: undefined: token"
```

A compatible terminal may render it as a navigable card. A traditional
terminal, a pipe or a log file still receives the textual fallback. The
producer describes the meaning; it does not draw Goatty-specific UI.

This boundary is deliberate. Goatty may understand text, images, progress,
tables, diagnostics and canvas content. It does not aim to become a widget
tree, layout engine, browser or application framework.

Read the full design in [The Semantic Layer](docs/PROTOCOL.md).

## Architecture

```text
                         Application
                              |
                   +----------+----------+
                   |                     |
                  UI               Terminal View
                                         |
                                  Renderer API
                                         |
                              Replaceable GPU backend

                         Terminal Core
        +------------------------------------------------+
        | Screen | Parser | Grid | PTY | Scrollback      |
        | Selection | Input | Images | Semantic objects  |
        +------------------------------------------------+
```

The direction is simple:

- the renderer does not own terminal state;
- the parser does not know about the renderer;
- the PTY does not know about the UI;
- applications depend on small interfaces, not concrete backends.

The result should feel like the **Ebitengine of terminal emulators**: useful as
an application, more valuable as a foundation.

## What works today

- GPU rendering with Ebitengine
- Unicode text, powerline glyphs and font ligatures
- embedded default fonts and installed monospaced fonts
- configurable themes, opacity and cursor images
- Sixel image rendering
- URL, colour, permission and timestamp hints
- screenshots and mouse selection
- headless terminal core for conformance testing
- Linux and macOS builds in CI
- native macOS builds for Apple Silicon and Intel

## Install

Download the binary for your platform from the
[latest release](https://github.com/kyamel/goatty/releases/latest):

| Platform | Release asset |
| --- | --- |
| Linux x86-64 | `goatty-linux-amd64` |
| macOS Apple Silicon | `goatty-darwin-arm64` |
| macOS Intel | `goatty-darwin-amd64` |

Then make it executable and place it somewhere in your `PATH`:

```sh
chmod +x goatty-*
mkdir -p "$HOME/.local/bin"
mv goatty-* "$HOME/.local/bin/goatty"
```

macOS binaries are currently unsigned. Gatekeeper may ask you to approve the
binary in **System Settings > Privacy & Security** the first time it runs.

### Build from source

Goatty requires Go 1.26 or newer.

On Debian or Ubuntu, install the native graphics dependencies first:

```sh
sudo apt install xorg-dev libgl1-mesa-dev
```

On macOS, install the Xcode command-line tools:

```sh
xcode-select --install
```

Then build:

```sh
git clone https://github.com/kyamel/goatty.git
cd goatty
make build
./goatty
```

With Nix:

```sh
nix build
./result/bin/goatty
```

## Configuration

| Platform | Configuration directory |
| --- | --- |
| Linux | `$XDG_CONFIG_HOME/goatty/` or `$HOME/.config/goatty/` |
| macOS | `$HOME/Library/Application Support/goatty/` |

Generate a starting configuration and inspect available fonts:

```sh
goatty --rewrite-config
goatty list-fonts
```

`config.yaml` may contain:

```yaml
opacity: 1.0
font:
  family: ""       # Empty uses the embedded font.
  size: 16
  dpi: 72
  ligatures: true
cursor:
  image: ""
```

Themes live beside the configuration as `theme.yaml`. Missing configuration
and theme files are fine; Goatty starts with built-in defaults.

## Key bindings

| Action | Binding |
| --- | --- |
| Copy | `Ctrl+Shift+C` |
| Paste | `Ctrl+Shift+V` |
| Decrease font size | `Ctrl+-` |
| Increase font size | `Ctrl+=` |
| Take screenshot | `Ctrl+Shift+[` |
| Open URL | `Ctrl+Click` |

## Project direction

The architectural work is intentionally incremental:

1. separate the terminal core from rendering and window management;
2. stabilize the interfaces between parser, grid, PTY and renderer;
3. expose reusable Go packages;
4. introduce semantic scrollback and command blocks;
5. add a small set of interoperable semantic object types;
6. define an extension boundary outside the terminal process;
7. move the protocol into an implementation-independent specification.

The detailed documents are:

- [Vision](docs/VISION.md) - the principles and long-term destination
- [Semantic protocol](docs/PROTOCOL.md) - meaning alongside ANSI and the line
  between terminal and application
- [Roadmap](docs/ROADMAP.md) - the dependency-ordered implementation plan

## Development

Enter the reproducible development environment:

```sh
nix develop
```

Run the normal checks:

```sh
go build ./...
go vet ./...
go test ./...
```

The repository also contains conformance tests, parser fuzzing, benchmarks and
golden renderer tests. See the [Makefile](Makefile) for the available targets.

Contributions should make the system easier to understand, replace or embed.
Performance matters, but an optimization that makes contributors afraid to
touch the code is too expensive.

## Platform support

| Platform | Status |
| --- | --- |
| Linux x86-64 | Supported and tested in CI |
| macOS Apple Silicon | Supported and tested in CI |
| macOS Intel | Supported and tested in CI |
| Windows | Not yet supported; the process layer needs ConPTY |

## From Darktile to Goatty

Goatty evolves the Darktile codebase rather than discarding a working GPU
terminal. Darktile proved that a fast terminal in Go is viable. Goatty's job is
to turn that implementation into a platform whose pieces can stand on their
own.

The migration is visible in package paths, command names and configuration.
That is intentional: architecture changes first, cosmetic renames follow when
they no longer create churn.

## License

[MIT](LICENSE)
