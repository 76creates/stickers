package stickers

import "github.com/charmbracelet/lipgloss"

// FlexBox responsive box grid insipred by CSS flexbox
type FlexBox struct {
	// style to apply to the gridbox itself
	style lipgloss.Style
	// width is fixed width of the box
	width int
	// height is fixed height of the box
	height int
	// fixedRowHeight will lock row height to a number, this disabless responsivness
	fixedRowHeight int

	rows []*FlexBoxRow

	// recalculateFlag indicates if next render should make calculations regarding
	// the rows objects height
	recalculateFlag bool
}

// NewFlexBox initialize FlexBox object with defaults
func NewFlexBox(width, height int) *FlexBox {
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

// NewRow initialize a new FlexBoxRow with width inherited from the FlexBox
func (r *FlexBox) NewRow() *FlexBoxRow {
	rw := &FlexBoxRow{
		cells: []*FlexBoxCell{},
		width: r.width,
		style: lipgloss.NewStyle(),
	}
	return rw
}

// AddRows appends additional rows to the FlexBox
func (r *FlexBox) AddRows(rows []*FlexBoxRow) *FlexBox {
	r.rows = append(r.rows, rows...)
	r.setRecalculate()
	return r
}

// SetRows replace rows on the FlexBox
func (r *FlexBox) SetRows(rows []*FlexBoxRow) *FlexBox {
	r.rows = rows
	r.setRecalculate()
	return r
}

// RowsLen returns the len of the rows slice
func (r *FlexBox) RowsLen() int {
	return len(r.rows)
}

// Row returns the FlexBoxRow on the given index if it exists
// note: forces the recalculation if found
func (r *FlexBox) Row(index int) *FlexBoxRow {
	if index >= 0 && len(r.rows) > 0 && index < len(r.rows) {
		r.setRecalculate()
		return r.rows[index]
		r.setRecalculate()
	}
	return nil
}

// GetRow returns the FlexBoxRow on the given index if it exists
// note: it does not return a pointer
func (r *FlexBox) GetRow(index int) (row FlexBoxRow, exists bool) {
	if index >= 0 && len(r.rows) > 0 && index < len(r.rows) {
		return *r.rows[index], true
	}
	return FlexBoxRow{}, false
}

// UpdateRow replaces the FlexBoxRow on the given index
func (r *FlexBox) UpdateRow(index int, row *FlexBoxRow) *FlexBox {
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
	return r.height
}

// GetWidth yields current FlexBox width
func (r *FlexBox) GetWidth() int {
	return r.width
}

// Render initiates the recalculation of the rows dimensions(height) if the recalculate flag is on,
// and then it renders all the rows and combines them on the vertical axis
func (r *FlexBox) Render() string {
	r.recalculate()
	var renderedRows []string
	for _, row := range r.rows {
		renderedRows = append(renderedRows, row.render())
	}
	// TODO: allow setting join align value for rows of variable width
	return r.style.Width(r.width).Height(r.height).Render(lipgloss.JoinVertical(lipgloss.Left, renderedRows...))
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
			r.distributeRowsHeight(r.calculateRowHeight())
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
	return calculateMatrixRatio(r.height, r.getRowMatrix())
}

// distributeRowsHeight sets height of each row per distribution array
func (r *FlexBox) distributeRowsHeight(ratioDistribution []int) {
	for index, row := range r.rows {
		row.setHeight(ratioDistribution[index])
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
