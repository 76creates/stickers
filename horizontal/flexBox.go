package horizontal

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

	// fixedColumnWidth will lock column height to a number, this disabless responsivness
	fixedColumnWidth int

	columns []*FlexBoxColumn

	// recalculateFlag indicates if next render should make calculations regarding
	// the columns objects height
	recalculateFlag bool
}

// NewFlexBox initialize FlexBox object with defaults
func NewFlexBox(width, height int) *FlexBox {
	r := &FlexBox{
		width:            width,
		height:           height,
		fixedColumnWidth: -1,
		style:            lipgloss.NewStyle(),
		recalculateFlag:  false,
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

// StylePassing set whether the style should be passed to the columns
func (r *FlexBox) StylePassing(value bool) *FlexBox {
	r.styleAncestor = value
	return r
}

// NewColumn initialize a new FlexBoxColumn with width inherited from the FlexBox
func (r *FlexBox) NewColumn() *FlexBoxColumn {
	rw := &FlexBoxColumn{
		cells: []*FlexBoxCell{},
		width: r.width,
		style: lipgloss.NewStyle(),
	}
	return rw
}

// AddColumns appends additional columns to the FlexBox
func (r *FlexBox) AddColumns(columns []*FlexBoxColumn) *FlexBox {
	r.columns = append(r.columns, columns...)
	r.setRecalculate()
	return r
}

// SetColumns replace columns on the FlexBox
func (r *FlexBox) SetColumns(columns []*FlexBoxColumn) *FlexBox {
	r.columns = columns
	r.setRecalculate()
	return r
}

// ColumnsLen returns the len of the columns slice
func (r *FlexBox) ColumnsLen() int {
	return len(r.columns)
}

// GetColumn returns the FlexBoxColumn on the given index if it exists
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *FlexBox) GetColumn(index int) *FlexBoxColumn {
	if index >= 0 && index < len(r.columns) {
		r.setRecalculate()
		return r.columns[index]
	}
	return nil
}

// GetColumnCopy returns a copy of the FlexBoxColumn on the given index, if column
// does not exist it will return nil. Copied column also gets copies of the
// cells. This is useful when you need to get columns attribute without
// triggering a recalculation.
func (r *FlexBox) GetColumnCopy(index int) *FlexBoxColumn {
	if index >= 0 && index < len(r.columns) {
		columnCopy := r.columns[index].copy()
		return &columnCopy
	}
	return nil
}

// GetColumnCellCopy returns a copy of the FlexBoxCell on the given index x,
// within the given column with index y, if column or cell do not exist it will
// return nil. This is useful when you need to get columns attribute without
// triggering a recalculation.
func (r *FlexBox) GetColumnCellCopy(columnIndex, cellIndex int) *FlexBoxCell {
	if columnIndex >= 0 && columnIndex < len(r.columns) {
		if cellIndex >= 0 && cellIndex < len(r.columns[columnIndex].cells) {
			cellCopy := r.columns[columnIndex].cells[cellIndex].copy()
			return &cellCopy
		}
	}
	return nil
}

// UpdateColumn replaces the FlexBoxColumn on the given index
func (r *FlexBox) UpdateColumn(index int, column *FlexBoxColumn) *FlexBox {
	r.columns[index] = column
	r.setRecalculate()
	return r
}

// LockColumnHeight sets the fixed height value for all the columns
// this will disable horizontal scaling
func (r *FlexBox) LockColumnHeight(value int) *FlexBox {
	r.fixedColumnWidth = value
	return r
}

// SetHeight sets the FlexBox height
func (r *FlexBox) SetHeight(value int) *FlexBox {
	r.height = value
	for _, column := range r.columns {
		column.setHeight(value)
	}
	return r
}

// SetWidth sets the FlexBox width
func (r *FlexBox) SetWidth(value int) *FlexBox {
	r.width = value
	r.setRecalculate()
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

// Render initiates the recalculation of the columns dimensions(height) if the recalculate flag is on,
// and then it renders all the columns and combines them on the vertical axis
func (r *FlexBox) Render() string {
	var inheritedStyle []lipgloss.Style
	if r.styleAncestor {
		inheritedStyle = append(inheritedStyle, r.style)
	}

	r.recalculate()
	var renderedColumns []string
	for _, column := range r.columns {
		renderedColumns = append(renderedColumns, column.render(inheritedStyle...))
	}
	// TODO: allow setting join align value for columns of variable width
	return r.style.
		Width(r.getContentWidth()).MaxWidth(r.getMaxWidth()).
		Height(r.getContentHeight()).MaxHeight(r.getMaxHeight()).
		//Render(lipgloss.JoinVertical(lipgloss.Left, renderedColumns...))
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...))
}

// ForceRecalculate forces the recalculation for the box and all the columns
func (r *FlexBox) ForceRecalculate() {
	r.recalculate()
	for _, rw := range r.columns {
		rw.recalculate()
	}
}

// recalculate fetches the column height distribution slice and sets it on the columns
func (r *FlexBox) recalculate() {
	if r.recalculateFlag {
		if len(r.columns) > 0 {
			r.distributeColumnsDimensions(r.calculateColumnWidth())
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

// calculateColumnWidth calculates the width of each column and returns the distribution array
func (r *FlexBox) calculateColumnWidth() (distribution []int) {
	if r.fixedColumnWidth > 0 {
		var fixedColumns []int
		for range r.columns {
			fixedColumns = append(fixedColumns, r.fixedColumnWidth)
		}
		return fixedColumns
	}
	return calculateMatrixRatio(r.getContentWidth(), r.getColumnMatrix())
}

// distributeColumnsDimensions sets height and width of each column per distribution array
func (r *FlexBox) distributeColumnsDimensions(ratioDistribution []int) {
	for index, column := range r.columns {
		column.setHeight(r.getContentHeight())
		column.setWidth(ratioDistribution[index])
	}
}

// getColumnMatrix return the matrix of the cell widths for all the columns
func (r *FlexBox) getColumnMatrix() (columnMatrix [][]int) {
	for _, column := range r.columns {
		var cellValues []int
		for _, cell := range column.cells {
			cellValues = append(cellValues, cell.ratioX)
		}
		columnMatrix = append(columnMatrix, cellValues)
	}
	return columnMatrix
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
