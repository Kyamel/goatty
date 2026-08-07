package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kyamel/goatty/internal/app/goatty/gui/render"
)

// Draw renders the terminal GUI to the ebtien window. Required to implement the ebiten interface.
func (g *GUI) Draw(screen *ebiten.Image) {
	// Handing the previous frame back lets it reuse the cell storage instead of
	// allocating a viewport on every redraw.
	g.frame = g.terminal.Snapshot(g.frame)

	render.
		New(screen, g.frame, g.fontManager, g.popupMessages, g.opacity, g.enableLigatures, g.cursorImage).
		Draw()

	if g.screenshotRequested {
		g.takeScreenshot(screen)
	}
}
