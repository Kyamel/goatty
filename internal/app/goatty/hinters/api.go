package hinters

import (
	"image"

	"github.com/kyamel/goatty/internal/app/goatty/termutil"
)

type HintAPI interface {
	ShowMessage(msg string)
	SetCursorToPointer()
	ResetCursor()
	Highlight(start termutil.Position, end termutil.Position, label string, img image.Image)
	ClearHighlight()
	CellSize() image.Point
}
