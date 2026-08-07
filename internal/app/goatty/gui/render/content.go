package render

func (r *Render) drawContent() {
	// draw base content for each row
	defBg := r.frame.Colours.Background
	defFg := r.frame.Colours.Foreground
	for viewY := int(r.frame.Height - 1); viewY >= 0; viewY-- {
		r.drawRow(viewY, defBg, defFg)
	}
}
