package flexbox

import (
	"github.com/charmbracelet/lipgloss"
)

// Cell is a building block object of the FlexBox, it represents a single cell within a box
// A FlexBox stacks cells horizontally.
// A HorizontalFlexBox stacks cells vertically. (controverse, isn't it?)
type Cell struct {
	// style of the cell, when rendering it will inherit the style of the parent row
	style lipgloss.Style
	// id of the cell, if not set it will default to the index in the row
	id string

	// TODO: all ratios and sizes should be uint
	// ratioX width ratio of the cell
	ratioX int
	// ratioY height ratio of the cell
	ratioY int
	// minWidth minimal width of the cell
	minWidth int
	// minHeight minimal height of the cell
	minHeight int

	width  int
	height int
	// contentGenerator is a function that generates the content of the cell depending on the
	// size of the cell, this can be useful for wrapping text or generating dynamic content.
	contentGenerator func(maxX, maxY int) string
}

// NewCell initialize FlexBoxCell object with defaults
func NewCell(ratioX, ratioY int) *Cell {
	return &Cell{
		style:    lipgloss.NewStyle(),
		ratioX:   ratioX,
		ratioY:   ratioY,
		minWidth: 0,
		width:    0,
		height:   0,
	}
}

// SetID sets the cells ID
func (r *Cell) SetID(id string) *Cell {
	r.id = id
	return r
}

// SetContent sets the cells content
func (r *Cell) SetContent(content string) *Cell {
	r.contentGenerator = func(_, _ int) string {
		return content
	}
	return r
}

// SetContentGenerator sets the cells content generator function
func (r *Cell) SetContentGenerator(generator func(maxX, maxY int) string) *Cell {
	r.contentGenerator = generator
	return r
}

// GetContent returns the cells raw content
func (r *Cell) GetContent() string {
	return r.contentGenerator(r.getMaxWidth(), r.getMaxHeight())
}

// SetMinWidth sets the cells minimum width, this will not disable responsivness.
// This has only an effect to cells of a normal FlexBox, not a HorizontalFlexBox.
func (r *Cell) SetMinWidth(value int) *Cell {
	r.minWidth = value
	return r
}

// Deprecated: use [*Cell.SetMinHeight]
func (r *Cell) SetMinHeigth(value int) *Cell {
	return r.SetMinHeight(value)
}

// SetMinHeight sets the cells minimum height, this will not disable responsivness.
// This has only an effect to cells of a HorizontalFlexBox.
func (r *Cell) SetMinHeight(value int) *Cell {
	r.minHeight = value
	return r
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *Cell) SetStyle(style lipgloss.Style) *Cell {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()
	return r
}

// GetStyle returns the copy of the cells current style
func (r *Cell) GetStyle() lipgloss.Style {
	return r.style
}

// GetWidth returns real width of the cell
func (r *Cell) GetWidth() int {
	return r.getMaxWidth()
}

// GetHeight returns real height of the cell
func (r *Cell) GetHeight() int {
	return r.getMaxHeight()
}

// render the cell into string
func (r *Cell) render(inherited ...lipgloss.Style) string {
	for _, style := range inherited {
		r.style = r.style.Inherit(style)
	}

	s := r.GetStyle().
		Width(r.getContentWidth()).MaxWidth(r.getMaxWidth()).
		Height(r.getContentHeight()).MaxHeight(r.getMaxHeight())
	return s.Render(r.GetContent())
}

func (r *Cell) getContentWidth() int {
	return r.getMaxWidth() - r.getExtraWidth()
}

func (r *Cell) getContentHeight() int {
	return r.getMaxHeight() - r.getExtraHeight()
}

func (r *Cell) getMaxWidth() int {
	return r.width
}

func (r *Cell) getMaxHeight() int {
	return r.height
}

func (r *Cell) getExtraWidth() int {
	return r.style.GetHorizontalMargins() + r.style.GetHorizontalBorderSize()
}

func (r *Cell) getExtraHeight() int {
	return r.style.GetVerticalMargins() + r.style.GetVerticalBorderSize()
}

func (r *Cell) copy() Cell {
	cellCopy := *r
	cellCopy.style = r.GetStyle()
	return cellCopy
}
