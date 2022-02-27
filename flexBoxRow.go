package stickers

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// FlexBoxRow is the container for the cells, this object has the least to do with the ratio
// of the construction as it takes all of the needed ratio information from the cell slice
// rows are stacked vertically
type FlexBoxRow struct {
	// style of the row, it will be passed to the cell when rendering for inheritance
	style lipgloss.Style

	cells []*FlexBoxCell

	height int
	width  int

	// recalculateFlag indicates if next render should make calculations regarding
	// the cells objects height/width
	recalculateFlag bool
}

// AddCells appends the cells to the row
// if the cell ID is not set it will default to the index of the cell
func (r *FlexBoxRow) AddCells(cells []*FlexBoxCell) *FlexBoxRow {
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
func (r *FlexBoxRow) CellsLen() int {
	return len(r.cells)
}

// Cell returns the FlexBoxCell on the given index if it exists
// note: forces the recalculation if found
func (r *FlexBoxRow) Cell(index int) *FlexBoxCell {
	if index >= 0 && len(r.cells) > 0 && index < len(r.cells) {
		r.setRecalculate()
		return r.cells[index]
	}
	return nil
}

// GetCellWithID returns the cell with the given ID if existing
func (r *FlexBoxRow) GetCellWithID(id string) (cell FlexBoxCell, exists bool) {
	for _, c := range r.cells {
		if c.id == id {
			return *c, true
		}
	}
	return FlexBoxCell{}, false
}

// GetCellWithIndex returns the cell with the given index if existing
// note: it does not return a pointer
func (r *FlexBoxRow) GetCellWithIndex(index int) (cell FlexBoxCell, exists bool) {
	if index >= 0 && len(r.cells) > 0 && index < len(r.cells) {
		return *r.cells[index], true
	}
	return FlexBoxCell{}, false
}

// MustGetCellWithIndex returns the cell with the given index if existing, panic if not
// note: it does not return a pointer
func (r *FlexBoxRow) MustGetCellWithIndex(index int) FlexBoxCell {
	return *r.cells[index]
}

// UpdateCellWithIndex replaces the cell on the given index if it exists
// if its not existing no changes will apply
func (r *FlexBoxRow) UpdateCellWithIndex(index int, cell *FlexBoxCell) {
	if index >= 0 && len(r.cells) > 0 && index < len(r.cells) {
		r.cells[index] = cell
		r.setRecalculate()
	}
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *FlexBoxRow) SetStyle(style lipgloss.Style) *FlexBoxRow {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()

	return r
}

func (r *FlexBoxRow) setHeight(value int) {
	r.height = value
	r.setRecalculate()
}

func (r *FlexBoxRow) setWidth(value int) {
	r.width = value
	r.setRecalculate()
}

func (r *FlexBoxRow) render() string {
	r.recalculate()
	var renderedCells []string
	for _, cell := range r.cells {
		renderedCells = append(renderedCells, cell.render(r.style))
	}
	return r.style.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedCells...))
}

func (r *FlexBoxRow) setRecalculate() {
	r.recalculateFlag = true
}

func (r *FlexBoxRow) unsetRecalculate() {
	r.recalculateFlag = false
}

// recalculate fetches the cells height/width distribution slices and sets it on the cells
func (r *FlexBoxRow) recalculate() {
	if r.recalculateFlag {
		if len(r.cells) > 0 {
			r.distributeCellDimensions(r.calculateCellsDimensions())
		}
		r.unsetRecalculate()
	}
}

// distributeCellDimensions sets height of each row per distribution array
func (r *FlexBoxRow) distributeCellDimensions(xMatrix, yMatrix []int) {
	for index, x := range xMatrix {
		r.cells[index].width = x
		r.cells[index].height = yMatrix[index]
	}
}

// calculateCellsDimensions calculates the height and width of the each cell
func (r *FlexBoxRow) calculateCellsDimensions() (xMatrix, yMatrix []int) {
	// calculate the cell height, it uses fixed combined ratio since the height of each cell
	// is indivitual and does not stack, row height will be calculated using the ratio of the
	// highest cell in the slice
	cellYMatrix, cellYMatrixMax := r.getCellHeightMatrix()
	// reminder not needed here due to how combined ratio is passed
	yMatrix, _ = distributeToMatrix(r.height, cellYMatrixMax, cellYMatrix)

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
		xMatrix = calculateRatioWithMinimum(r.width, r.getCellWidthMatrix(), minWidthMatrix)
	} else {
		xMatrix = calculateRatio(r.width, r.getCellWidthMatrix())
	}

	return xMatrix, yMatrix
}

// getCellHeightMatrix return the matrix of the cell height in cells and the max value in it
func (r *FlexBoxRow) getCellHeightMatrix() (cellHeightMatrix []int, max int) {
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
func (r *FlexBoxRow) getCellWidthMatrix() (cellWidthMatrix []int) {
	for _, cell := range r.cells {
		cellWidthMatrix = append(cellWidthMatrix, cell.ratioX)
	}
	return cellWidthMatrix
}
