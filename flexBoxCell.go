package stickers

import (
	"github.com/charmbracelet/lipgloss"
)

// FlexBoxCell is a building block object of the FlexBox, it represents a single cell within a box
// cells are stacked horizontally
type FlexBoxCell struct {
	// style of the cell, when rendering it will inherit the style of the parent row
	style lipgloss.Style
	// id of the cell, if not set it will default to the index in the row
	id string

	// ratioX width ratio of the cell
	ratioX int
	// ratioY height ratio of the cell
	ratioY int
	// minWidth minimal width of the cell
	minWidth int
	// TODO: implement minimal height

	width   int
	height  int
	content string
}

// NewFlexBoxCell initialize FlexBoxCell object with defaults
func NewFlexBoxCell(ratioX, ratioY int) *FlexBoxCell {
	return &FlexBoxCell{
		style:    lipgloss.NewStyle(),
		ratioX:   ratioX,
		ratioY:   ratioY,
		minWidth: 0,
		width:    0,
		height:   0,
	}
}

// SetID sets the cells ID
func (r *FlexBoxCell) SetID(id string) *FlexBoxCell {
	r.id = id
	return r
}

// SetContent sets the cells content
func (r *FlexBoxCell) SetContent(content string) *FlexBoxCell {
	r.content = content
	return r
}

// GetContent returns the cells raw content
func (r *FlexBoxCell) GetContent() string {
	return r.content
}

// SetMinWidth sets the cells minimum width, this will not disable responsivness
func (r *FlexBoxCell) SetMinWidth(value int) *FlexBoxCell {
	r.minWidth = value
	return r
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *FlexBoxCell) SetStyle(style lipgloss.Style) *FlexBoxCell {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()
	return r
}

// GetStyle returns the copy of the cells current style
func (r *FlexBoxCell) GetStyle() lipgloss.Style {
	return r.style.Copy()
}

// GetWidth returns real width of the cell
func (r *FlexBoxCell) GetWidth() int {
	return r.width - r.style.GetHorizontalMargins()
}

// GetHeight returns real height of the cell
func (r *FlexBoxCell) GetHeight() int {
	return r.height - r.style.GetVerticalMargins()
}

// render renders the cell into string
func (r *FlexBoxCell) render(rowStyle lipgloss.Style) string {
	// width and height include margins so we subtract them
	w := r.GetWidth()
	h := r.GetHeight()

	s := r.GetStyle().
		Inherit(rowStyle).
		Width(w).MaxWidth(w).
		Height(h).MaxHeight(h)
	renderLen := 0

	if len(r.content) < w {
		renderLen = len(r.content)
	} else {
		renderLen = w
	}
	return s.Render(r.content[:renderLen])
}
