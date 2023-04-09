package flexbox

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// Row is the container for the cells, this object has the least to do with the ratio
// of the construction as it takes all of the needed ratio information from the cell slice
// rows are stacked vertically
type Row struct {
	// style of the row
	style         lipgloss.Style
	styleAncestor bool

	cells []*Cell

	height int
	width  int

	// recalculateFlag indicates if next render should make calculations regarding
	// the cells objects height/width
	recalculateFlag bool
}

// AddCells appends the cells to the row
// if the cell ID is not set it will default to the index of the cell
func (r *Row) AddCells(cells ...*Cell) *Row {
	r.cells = append(r.cells, cells...)
	for i, cell := range r.cells {
		if cell.id == "" {
			cell.id = strconv.Itoa(i)
		}
	}
	r.setRecalculate()
	return r
}

// CellsLen returns the len of the cells slice
func (r *Row) CellsLen() int {
	return len(r.cells)
}

// GetCell returns the FlexBoxCell on the given index if it exists
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *Row) GetCell(index int) *Cell {
	if index >= 0 && index < len(r.cells) {
		r.setRecalculate()
		return r.cells[index]
	}
	return nil
}

// GetCellCopy returns a copy of the FlexBoxCell on the given index, if cell
// does not exist it will return nil. This is useful when you need to get
// cells attribute without triggering a recalculation.
func (r *Row) GetCellCopy(index int) *Cell {
	if index >= 0 && index < len(r.cells) {
		c := r.cells[index].copy()
		return &c
	}
	return nil
}

// GetCellWithID returns the cell with the given ID if existing
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *Row) GetCellWithID(id string) *Cell {
	for _, c := range r.cells {
		if c.id == id {
			r.setRecalculate()
			return c
		}
	}
	return nil
}

// UpdateCellWithIndex replaces the cell on the given index if it exists
// if its not existing no changes will apply
func (r *Row) UpdateCellWithIndex(index int, cell *Cell) {
	if index >= 0 && len(r.cells) > 0 && index < len(r.cells) {
		r.cells[index] = cell
		r.setRecalculate()
	}
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *Row) SetStyle(style lipgloss.Style) *Row {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()

	return r
}

// StylePassing set whether the style should be passed to the cells
func (r *Row) StylePassing(value bool) *Row {
	r.styleAncestor = value
	return r
}

func (r *Row) setHeight(value int) {
	r.height = value
	r.setRecalculate()
}

func (r *Row) setWidth(value int) {
	r.width = value
	r.setRecalculate()
}

func (r *Row) render(inherited ...lipgloss.Style) string {
	var inheritedStyle []lipgloss.Style

	for _, style := range inherited {
		r.style = r.style.Inherit(style)
	}

	// intentionally applied after row inherits the box style
	if r.styleAncestor {
		inheritedStyle = append(inheritedStyle, r.style)
	}

	r.recalculate()
	var renderedCells []string
	for _, cell := range r.cells {
		renderedCells = append(renderedCells, cell.render(inheritedStyle...))
	}
	return r.style.
		Width(r.getContentWidth()).MaxWidth(r.getMaxWidth()).
		Height(r.getContentHeight()).MaxHeight(r.getMaxHeight()).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedCells...))
}

func (r *Row) setRecalculate() {
	r.recalculateFlag = true
}

func (r *Row) unsetRecalculate() {
	r.recalculateFlag = false
}

// recalculate fetches the cell's height/width distribution slices and sets it on the cells
func (r *Row) recalculate() {
	if r.recalculateFlag {
		if len(r.cells) > 0 {
			r.distributeCellDimensions(r.calculateCellsDimensions())
		}
		r.unsetRecalculate()
	}
}

// distributeCellDimensions sets height of each row per distribution array
func (r *Row) distributeCellDimensions(xMatrix, yMatrix []int) {
	for index, x := range xMatrix {
		r.cells[index].width = x
		r.cells[index].height = yMatrix[index]
	}
}

// calculateCellsDimensions calculates the height and width of the each cell
func (r *Row) calculateCellsDimensions() (xMatrix, yMatrix []int) {
	// calculate the cell height, it uses fixed combined ratio since the height of each cell
	// is individual and does not stack, row height will be calculated using the ratio of the
	// highest cell in the slice
	cellYMatrix, cellYMatrixMax := r.getCellHeightMatrix()
	// reminder not needed here due to how combined ratio is passed
	yMatrix, _ = distributeToMatrix(r.getContentHeight(), cellYMatrixMax, cellYMatrix)

	// get the min width matrix of the cells if any
	withMinWidth := false
	var minWidthMatrix []int
	for _, c := range r.cells {
		minWidthMatrix = append(minWidthMatrix, c.minWidth)
		if c.minWidth > 0 {
			withMinWidth = true
		}
	}

	// calculate the cell width matrix
	if withMinWidth {
		xMatrix = calculateRatioWithMinimum(r.getContentWidth(), r.getCellWidthMatrix(), minWidthMatrix)
	} else {
		xMatrix = calculateRatio(r.getContentWidth(), r.getCellWidthMatrix())
	}

	return xMatrix, yMatrix
}

// getCellHeightMatrix return the matrix of the cell height in cells and the max value in it
func (r *Row) getCellHeightMatrix() (cellHeightMatrix []int, max int) {
	max = 0
	for _, cell := range r.cells {
		if cell.ratioY > max {
			max = cell.ratioY
		}
		cellHeightMatrix = append(cellHeightMatrix, cell.ratioY)
	}
	return cellHeightMatrix, max
}

// getCellWidthMatrix return the matrix of the cell width in cells
func (r *Row) getCellWidthMatrix() (cellWidthMatrix []int) {
	for _, cell := range r.cells {
		cellWidthMatrix = append(cellWidthMatrix, cell.ratioX)
	}
	return cellWidthMatrix
}

func (r *Row) getContentWidth() int {
	return r.getMaxWidth() - r.getExtraWidth()
}

func (r *Row) getContentHeight() int {
	return r.getMaxHeight() - r.getExtraHeight()
}

func (r *Row) getMaxWidth() int {
	return r.width
}

func (r *Row) getMaxHeight() int {
	return r.height
}

func (r *Row) getExtraWidth() int {
	return r.style.GetHorizontalMargins() + r.style.GetHorizontalBorderSize()
}

func (r *Row) getExtraHeight() int {
	return r.style.GetVerticalMargins() + r.style.GetVerticalBorderSize()
}

func (r *Row) copy() Row {
	var cells []*Cell
	for _, cell := range r.cells {
		cellCopy := cell.copy()
		cells = append(cells, &cellCopy)
	}
	rowCopy := *r
	rowCopy.cells = cells
	rowCopy.style = r.style.Copy()

	return rowCopy
}
