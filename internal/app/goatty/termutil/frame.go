package termutil

import "image/color"

// Frame is a snapshot of everything needed to draw the terminal once. It is
// filled under the terminal's lock, so whatever consumes it neither touches
// live buffer state nor needs the lock, and a slow renderer cannot stall the
// parser.
type Frame struct {
	Width  uint16
	Height uint16

	Cursor    Cursor
	Selection *Selection
	Highlight *Highlight
	Sixels    []VisibleSixel
	Colours   Colours

	cells   []Cell
	present []bool
}

type Cursor struct {
	Col     uint16
	Line    uint16
	Shape   CursorShape
	Visible bool
}

// Highlight is the region a hinter has marked, plus the annotation to show
// against it, if it produced one.
type Highlight struct {
	Start      Position
	End        Position
	Annotation *Annotation
}

// Colours are the theme entries a renderer needs, resolved here so that it
// does not have to carry the theme.
type Colours struct {
	Foreground          color.Color
	Background          color.Color
	CursorForeground    color.Color
	CursorBackground    color.Color
	SelectionForeground color.Color
	SelectionBackground color.Color
}

// Cell returns the cell at a viewport position, or nil where the buffer holds
// none: past the end of a short line, or below the last line written.
func (f *Frame) Cell(col, row uint16) *Cell {
	if col >= f.Width || row >= f.Height {
		return nil
	}
	i := int(row)*int(f.Width) + int(col)
	if !f.present[i] {
		return nil
	}
	return &f.cells[i]
}

// Snapshot copies the visible state into frame and returns it, allocating one
// when frame is nil. Handing the previous frame back reuses its cell storage,
// which is what keeps a steady redraw from allocating a viewport per draw.
func (t *Terminal) Snapshot(frame *Frame) *Frame {
	t.mu.Lock()
	defer t.mu.Unlock()

	if frame == nil {
		frame = &Frame{}
	}

	b := t.activeBuffer
	frame.Width, frame.Height = b.ViewWidth(), b.ViewHeight()

	size := int(frame.Width) * int(frame.Height)
	if cap(frame.cells) < size {
		frame.cells = make([]Cell, size)
		frame.present = make([]bool, size)
	}
	frame.cells = frame.cells[:size]
	frame.present = frame.present[:size]

	for row := uint16(0); row < frame.Height; row++ {
		base := int(row) * int(frame.Width)
		for col := uint16(0); col < frame.Width; col++ {
			cell := b.GetCell(col, row)
			if cell == nil {
				frame.present[base+int(col)] = false
				continue
			}
			frame.cells[base+int(col)] = *cell
			frame.present[base+int(col)] = true
		}
	}

	frame.Cursor = Cursor{
		Col:     b.CursorColumn(),
		Line:    b.CursorLine(),
		Shape:   b.GetCursorShape(),
		Visible: b.IsCursorVisible(),
	}

	_, frame.Selection = b.GetSelection()

	frame.Highlight = nil
	if start, end, ok := b.GetViewHighlight(); ok {
		frame.Highlight = &Highlight{
			Start:      start,
			End:        end,
			Annotation: b.GetHighlightAnnotation(),
		}
	}

	frame.Sixels = b.GetVisibleSixels()

	frame.Colours = Colours{
		Foreground:          t.theme.DefaultForeground(),
		Background:          t.theme.DefaultBackground(),
		CursorForeground:    t.theme.CursorForeground(),
		CursorBackground:    t.theme.CursorBackground(),
		SelectionForeground: t.theme.SelectionForeground(),
		SelectionBackground: t.theme.SelectionBackground(),
	}

	return frame
}
