package flexbox

import "github.com/charmbracelet/lipgloss"

// FlexBox responsive box grid insipred by CSS flexbox
type FlexBox struct {
	// style to apply to the gridbox itself
	style         lipgloss.Style
	styleAncestor bool

	// width is fixed width of the box
	width int
	// height is fixed height of the box
	height int
	// fixedRowHeight will lock row height to a number, this disabless responsivness
	fixedRowHeight int

	rows []*Row

	// recalculateFlag indicates if next render should make calculations regarding
	// the rows objects height
	recalculateFlag bool
}

// New initialize FlexBox object with defaults
func New(width, height int) *FlexBox {
	r := &FlexBox{
		width:           width,
		height:          height,
		fixedRowHeight:  -1,
		style:           lipgloss.NewStyle(),
		recalculateFlag: false,
	}
	return r
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *FlexBox) SetStyle(style lipgloss.Style) *FlexBox {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()
	return r
}

// StylePassing set whether the style should be passed to the rows
func (r *FlexBox) StylePassing(value bool) *FlexBox {
	r.styleAncestor = value
	return r
}

// NewRow initialize a new Row with width inherited from the FlexBox
func (r *FlexBox) NewRow() *Row {
	rw := &Row{
		cells: []*Cell{},
		width: r.width,
		style: lipgloss.NewStyle(),
	}
	return rw
}

// AddRows appends additional rows to the FlexBox
func (r *FlexBox) AddRows(rows []*Row) *FlexBox {
	r.rows = append(r.rows, rows...)
	r.setRecalculate()
	return r
}

// SetRows replace rows on the FlexBox
func (r *FlexBox) SetRows(rows []*Row) *FlexBox {
	r.rows = rows
	r.setRecalculate()
	return r
}

// RowsLen returns the len of the rows slice
func (r *FlexBox) RowsLen() int {
	return len(r.rows)
}

// GetRow returns the Row on the given index if it exists
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *FlexBox) GetRow(index int) *Row {
	if index >= 0 && index < len(r.rows) {
		r.setRecalculate()
		return r.rows[index]
	}
	return nil
}

// GetRowCopy returns a copy of the Row on the given index, if row
// does not exist it will return nil. Copied row also gets copies of the
// cells. This is useful when you need to get rows attribute without
// triggering a recalculation.
func (r *FlexBox) GetRowCopy(index int) *Row {
	if index >= 0 && index < len(r.rows) {
		rowCopy := r.rows[index].copy()
		return &rowCopy
	}
	return nil
}

// GetRowCellCopy returns a copy of the FlexBoxCell on the given index x,
// within the given row with index y, if row or cell do not exist it will
// return nil. This is useful when you need to get rows attribute without
// triggering a recalculation.
func (r *FlexBox) GetRowCellCopy(rowIndex, cellIndex int) *Cell {
	if rowIndex >= 0 && rowIndex < len(r.rows) {
		if cellIndex >= 0 && cellIndex < len(r.rows[rowIndex].cells) {
			cellCopy := r.rows[rowIndex].cells[cellIndex].copy()
			return &cellCopy
		}
	}
	return nil
}

// UpdateRow replaces the Row on the given index
func (r *FlexBox) UpdateRow(index int, row *Row) *FlexBox {
	r.rows[index] = row
	r.setRecalculate()
	return r
}

// LockRowHeight sets the fixed height value for all the rows
// this will disable vertical scaling
func (r *FlexBox) LockRowHeight(value int) *FlexBox {
	r.fixedRowHeight = value
	return r
}

// SetHeight sets the FlexBox height
func (r *FlexBox) SetHeight(value int) *FlexBox {
	r.height = value
	r.setRecalculate()
	return r
}

// SetWidth sets the FlexBox width
func (r *FlexBox) SetWidth(value int) *FlexBox {
	r.width = value
	for _, row := range r.rows {
		row.setWidth(value)
	}
	return r
}

// GetHeight yields current FlexBox height
func (r *FlexBox) GetHeight() int {
	return r.getMaxHeight()
}

// GetWidth yields current FlexBox width
func (r *FlexBox) GetWidth() int {
	return r.getMaxWidth()
}

// Render initiates the recalculation of the rows dimensions(height) if the recalculate flag is on,
// and then it renders all the rows and combines them on the vertical axis
func (r *FlexBox) Render() string {
	var inheritedStyle []lipgloss.Style
	if r.styleAncestor {
		inheritedStyle = append(inheritedStyle, r.style)
	}

	r.recalculate()
	var renderedRows []string
	for _, row := range r.rows {
		renderedRows = append(renderedRows, row.render(inheritedStyle...))
	}
	// TODO: allow setting join align value for rows of variable width
	return r.style.
		Width(r.getContentWidth()).MaxWidth(r.getMaxWidth()).
		Height(r.getContentHeight()).MaxHeight(r.getMaxHeight()).
		Render(lipgloss.JoinVertical(lipgloss.Left, renderedRows...))
}

// ForceRecalculate forces the recalculation for the box and all the rows
func (r *FlexBox) ForceRecalculate() {
	r.recalculate()
	for _, rw := range r.rows {
		rw.recalculate()
	}
}

// recalculate fetches the row height distribution slice and sets it on the rows
func (r *FlexBox) recalculate() {
	if r.recalculateFlag {
		if len(r.rows) > 0 {
			r.distributeRowsDimensions(r.calculateRowHeight())
		}
		r.unsetRecalculate()
	}
}

func (r *FlexBox) setRecalculate() {
	r.recalculateFlag = true
}

func (r *FlexBox) unsetRecalculate() {
	r.recalculateFlag = false
}

// calculateRowHeight calculates the height of each row and returns the distribution array
func (r *FlexBox) calculateRowHeight() (distribution []int) {
	if r.fixedRowHeight > 0 {
		var fixedRows []int
		for range r.rows {
			fixedRows = append(fixedRows, r.fixedRowHeight)
		}
		return fixedRows
	}
	return calculateMatrixRatio(r.getContentHeight(), r.getRowMatrix())
}

// distributeRowsDimensions sets height and width of each row per distribution array
func (r *FlexBox) distributeRowsDimensions(ratioDistribution []int) {
	for index, row := range r.rows {
		row.setHeight(ratioDistribution[index])
		row.setWidth(r.getContentWidth())
	}
}

// getRowMatrix return the matrix of the cell hights for all the rows
func (r *FlexBox) getRowMatrix() (rowMatrix [][]int) {
	for _, row := range r.rows {
		var cellValues []int
		for _, cell := range row.cells {
			cellValues = append(cellValues, cell.ratioY)
		}
		rowMatrix = append(rowMatrix, cellValues)
	}
	return rowMatrix
}

func (r *FlexBox) getContentWidth() int {
	return r.getMaxWidth() - r.getExtraWidth()
}

func (r *FlexBox) getContentHeight() int {
	return r.getMaxHeight() - r.getExtraHeight()
}

func (r *FlexBox) getMaxWidth() int {
	return r.width
}

func (r *FlexBox) getMaxHeight() int {
	return r.height
}

func (r *FlexBox) getExtraWidth() int {
	return r.style.GetHorizontalMargins() + r.style.GetHorizontalBorderSize()
}

func (r *FlexBox) getExtraHeight() int {
	return r.style.GetVerticalMargins() + r.style.GetVerticalBorderSize()
}
