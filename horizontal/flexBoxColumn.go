package horizontal

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// FlexBoxColumn is the container for the cells, this object has the least to do with the ratio
// of the construction as it takes all of the needed ratio information from the cell slice
// columns are stacked vertically
type FlexBoxColumn struct {
	// style of the column
	style         lipgloss.Style
	styleAncestor bool

	cells []*FlexBoxCell

	height int
	width  int

	// recalculateFlag indicates if next render should make calculations regarding
	// the cells objects height/width
	recalculateFlag bool
}

// AddCells appends the cells to the column
// if the cell ID is not set it will default to the index of the cell
func (r *FlexBoxColumn) AddCells(cells ...*FlexBoxCell) *FlexBoxColumn {
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
func (r *FlexBoxColumn) CellsLen() int {
	return len(r.cells)
}

// GetCell returns the FlexBoxCell on the given index if it exists
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *FlexBoxColumn) GetCell(index int) *FlexBoxCell {
	if index >= 0 && index < len(r.cells) {
		r.setRecalculate()
		return r.cells[index]
	}
	return nil
}

// GetCellCopy returns a copy of the FlexBoxCell on the given index, if cell
// does not exist it will return nil. This is useful when you need to get
// cells attribute without triggering a recalculation.
func (r *FlexBoxColumn) GetCellCopy(index int) *FlexBoxCell {
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
func (r *FlexBoxColumn) GetCellWithID(id string) *FlexBoxCell {
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
func (r *FlexBoxColumn) UpdateCellWithIndex(index int, cell *FlexBoxCell) {
	if index >= 0 && len(r.cells) > 0 && index < len(r.cells) {
		r.cells[index] = cell
		r.setRecalculate()
	}
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *FlexBoxColumn) SetStyle(style lipgloss.Style) *FlexBoxColumn {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()

	return r
}

// StylePassing set whether the style should be passed to the cells
func (r *FlexBoxColumn) StylePassing(value bool) *FlexBoxColumn {
	r.styleAncestor = value
	return r
}

func (r *FlexBoxColumn) setHeight(value int) {
	r.height = value
	r.setRecalculate()
}

func (r *FlexBoxColumn) setWidth(value int) {
	r.width = value
	r.setRecalculate()
}

func (r *FlexBoxColumn) render(inherited ...lipgloss.Style) string {
	var inheritedStyle []lipgloss.Style

	for _, style := range inherited {
		r.style = r.style.Inherit(style)
	}

	// intentionally applied after column inherits the box style
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
		//Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedCells...))
		Render(lipgloss.JoinVertical(lipgloss.Left, renderedCells...))
}

func (r *FlexBoxColumn) setRecalculate() {
	r.recalculateFlag = true
}

func (r *FlexBoxColumn) unsetRecalculate() {
	r.recalculateFlag = false
}

// recalculate fetches the cell's height/width distribution slices and sets it on the cells
func (r *FlexBoxColumn) recalculate() {
	if r.recalculateFlag {
		if len(r.cells) > 0 {
			r.distributeCellDimensions(r.calculateCellsDimensions())
		}
		r.unsetRecalculate()
	}
}

// distributeCellDimensions sets width of each column per distribution array
func (r *FlexBoxColumn) distributeCellDimensions(xMatrix, yMatrix []int) {
	for index, y := range yMatrix {
		r.cells[index].width = xMatrix[index]
		r.cells[index].height = y
	}
}

// calculateCellsDimensions calculates the height and width of the each cell
func (r *FlexBoxColumn) calculateCellsDimensions() (xMatrix, yMatrix []int) {
	// calculate the cell height, it uses fixed combined ratio since the height of each cell
	// is individual and does not stack, column height will be calculated using the ratio of the
	// highest cell in the slice
	cellXMatrix, cellXMatrixMax := r.getCellWidthMatrix()

	// reminder not needed here due to how combined ratio is passed
	xMatrix, _ = distributeToMatrix(r.getContentWidth(), cellXMatrixMax, cellXMatrix)

	// get the min width matrix of the cells if any
	withMinHeigth := false
	var minHeigthMatrix []int
	for _, c := range r.cells {
		minHeigthMatrix = append(minHeigthMatrix, c.minHeigth)
		if c.minHeigth > 0 {
			withMinHeigth = true
		}
	}

	// calculate the cell width matrix
	if withMinHeigth {
		yMatrix = calculateRatioWithMinimum(r.getContentHeight(), r.getCellHeightMatrix(), minHeigthMatrix)
	} else {
		yMatrix = calculateRatio(r.getContentHeight(), r.getCellHeightMatrix())
	}

	return xMatrix, yMatrix
}

// getCellHeightMatrix return the matrix of the cell height in cells
func (r *FlexBoxColumn) getCellHeightMatrix() (cellHeightMatrix []int) {
	for _, cell := range r.cells {
		cellHeightMatrix = append(cellHeightMatrix, cell.ratioY)
	}
	return cellHeightMatrix
}

// getCellWidthMatrix return the matrix of the cell width in cells and the max value in it
func (r *FlexBoxColumn) getCellWidthMatrix() (cellWidthMatrix []int, max int) {
	max = 0
	for _, cell := range r.cells {
		if cell.ratioX > max {
			max = cell.ratioX
		}
		cellWidthMatrix = append(cellWidthMatrix, cell.ratioX)
	}
	return cellWidthMatrix, max
}

func (r *FlexBoxColumn) getContentWidth() int {
	return r.getMaxWidth() - r.getExtraWidth()
}

func (r *FlexBoxColumn) getContentHeight() int {
	return r.getMaxHeight() - r.getExtraHeight()
}

func (r *FlexBoxColumn) getMaxWidth() int {
	return r.width
}

func (r *FlexBoxColumn) getMaxHeight() int {
	return r.height
}

func (r *FlexBoxColumn) getExtraWidth() int {
	return r.style.GetHorizontalMargins() + r.style.GetHorizontalBorderSize()
}

func (r *FlexBoxColumn) getExtraHeight() int {
	return r.style.GetVerticalMargins() + r.style.GetVerticalBorderSize()
}

func (r *FlexBoxColumn) copy() FlexBoxColumn {
	var cells []*FlexBoxCell
	for _, cell := range r.cells {
		cellCopy := cell.copy()
		cells = append(cells, &cellCopy)
	}
	columnCopy := *r
	columnCopy.cells = cells
	columnCopy.style = r.style.Copy()

	return columnCopy
}
