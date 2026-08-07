package render

import "github.com/hajimehoshi/ebiten/v2"

func (r *Render) drawSixels() {
	for _, sixel := range r.frame.Sixels {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			float64(int(sixel.Sixel.X)*r.font.CellSize.X),
			float64(sixel.ViewLineOffset*r.font.CellSize.Y),
		)
		r.canvas.DrawImage(
			ebiten.NewImageFromImage(sixel.Sixel.Image),
			op,
		)
	}
}
