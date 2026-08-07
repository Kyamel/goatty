package render

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (r *Render) drawSelection() {
	selection := r.frame.Selection
	if selection == nil {
		// nothing selected
		return
	}

	bg, fg := r.frame.Colours.SelectionBackground, r.frame.Colours.SelectionForeground

	for y := selection.Start.Line; y <= selection.End.Line; y++ {
		xStart, xEnd := 0, int(r.frame.Width)
		if y == selection.Start.Line {
			xStart = int(selection.Start.Col)
		}
		if y == selection.End.Line {
			xEnd = int(selection.End.Col)
		}
		for x := xStart; x <= xEnd; x++ {
			pX, pY := float64(x*r.font.CellSize.X), float64(y*uint64(r.font.CellSize.Y))
			ebitenutil.DrawRect(r.canvas, pX, pY, float64(r.font.CellSize.X), float64(r.font.CellSize.Y), bg)
			cell := r.frame.Cell(uint16(x), uint16(y))
			if cell == nil || cell.Rune().Rune == 0 {
				continue
			}
			text.Draw(r.canvas, string(cell.Rune().Rune), r.font.Regular, int(pX), int(pY)+r.font.DotDepth, fg)
		}
	}
}
