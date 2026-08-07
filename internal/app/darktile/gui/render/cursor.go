package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kyamel/goatty/internal/app/darktile/termutil"
)

func (r *Render) drawCursor() {
	//draw cursor
	cursor := r.frame.Cursor
	if !cursor.Visible {
		return
	}

	pixelX := float64(int(cursor.Col) * r.font.CellSize.X)
	pixelY := float64(int(cursor.Line) * r.font.CellSize.Y)
	cell := r.frame.Cell(cursor.Col, cursor.Line)

	useFace := r.font.Regular
	if cell != nil {
		if cell.Bold() && cell.Italic() {
			useFace = r.font.BoldItalic
		} else if cell.Bold() {
			useFace = r.font.Bold
		} else if cell.Italic() {
			useFace = r.font.Italic
		}
	}

	pixelW, pixelH := float64(r.font.CellSize.X), float64(r.font.CellSize.Y)

	// empty rect without focus
	if !ebiten.IsFocused() {
		ebitenutil.DrawRect(r.canvas, pixelX, pixelY, pixelW, pixelH, r.frame.Colours.CursorBackground)
		ebitenutil.DrawRect(r.canvas, pixelX+1, pixelY+1, pixelW-2, pixelH-2, r.frame.Colours.CursorForeground)
		return
	}

	// draw the cursor shape
	switch cursor.Shape {
	case termutil.CursorShapeBlinkingBar, termutil.CursorShapeSteadyBar:
		ebitenutil.DrawRect(r.canvas, pixelX, pixelY, 2, pixelH, r.frame.Colours.CursorBackground)
	case termutil.CursorShapeBlinkingUnderline, termutil.CursorShapeSteadyUnderline:
		ebitenutil.DrawRect(r.canvas, pixelX, pixelY+pixelH-2, pixelW, 2, r.frame.Colours.CursorBackground)
	default:
		// draw a custom cursor if we have one and there are no characters in the way
		if r.cursorImage != nil && (cell == nil || cell.Rune().Rune == 0) {
			opt := &ebiten.DrawImageOptions{}
			_, h := r.cursorImage.Size()
			ratio := 1 / (float64(h) / float64(r.font.CellSize.Y))
			actualHeight := float64(h) * ratio
			offsetY := (float64(r.font.CellSize.Y) - actualHeight) / 2
			opt.GeoM.Scale(ratio, ratio)
			opt.GeoM.Translate(pixelX, pixelY+offsetY)
			r.canvas.DrawImage(r.cursorImage, opt)
			return
		}

		ebitenutil.DrawRect(r.canvas, pixelX, pixelY, pixelW, pixelH, r.frame.Colours.CursorBackground)

		// we've drawn over the cell contents, so we need to draw it again in the cursor colours
		if cell != nil && cell.Rune().Rune > 0 {
			text.Draw(r.canvas, string(cell.Rune().Rune), useFace, int(pixelX), int(pixelY)+r.font.DotDepth, r.frame.Colours.CursorForeground)
		}
	}
}
