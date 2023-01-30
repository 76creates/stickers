package flexbox

import "github.com/charmbracelet/lipgloss"

// HorizontalFlexBox responsive box grid insipred by CSS flexbox
type HorizontalFlexBox struct {
	// style to apply to the gridbox itself
	style         lipgloss.Style
	styleAncestor bool

	// width is fixed width of the box
	width int
	// height is fixed height of the box
	height int

	// fixedColumnWidth will lock column width to a number, this disabless responsivness
	fixedColumnWidth int

	columns []*Column

	// recalculateFlag indicates if next render should make calculations regarding
	// the columns objects height
	recalculateFlag bool
}

// NewHorizontal initialize a HorizontalFlexBox object with defaults
func NewHorizontal(width, height int) *HorizontalFlexBox {
	r := &HorizontalFlexBox{
		width:            width,
		height:           height,
		fixedColumnWidth: -1,
		style:            lipgloss.NewStyle(),
		recalculateFlag:  false,
	}
	return r
}

// SetStyle replaces the style, it unsets width/height related keys
func (r *HorizontalFlexBox) SetStyle(style lipgloss.Style) *HorizontalFlexBox {
	r.style = style.
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight()
	return r
}

// StylePassing set whether the style should be passed to the columns
func (r *HorizontalFlexBox) StylePassing(value bool) *HorizontalFlexBox {
	r.styleAncestor = value
	return r
}

// NewColumn initialize a new FlexBoxColumn with width inherited from the FlexBox
func (r *HorizontalFlexBox) NewColumn() *Column {
	rw := &Column{
		cells: []*Cell{},
		width: r.width,
		style: lipgloss.NewStyle(),
	}
	return rw
}

// AddColumns appends additional columns to the FlexBox
func (r *HorizontalFlexBox) AddColumns(columns []*Column) *HorizontalFlexBox {
	r.columns = append(r.columns, columns...)
	r.setRecalculate()
	return r
}

// SetColumns replace columns on the FlexBox
func (r *HorizontalFlexBox) SetColumns(columns []*Column) *HorizontalFlexBox {
	r.columns = columns
	r.setRecalculate()
	return r
}

// ColumnsLen returns the len of the columns slice
func (r *HorizontalFlexBox) ColumnsLen() int {
	return len(r.columns)
}

// GetColumn returns the FlexBoxColumn on the given index if it exists
// note: forces the recalculation if found
//
//	returns nil if not found
func (r *HorizontalFlexBox) GetColumn(index int) *Column {
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
func (r *HorizontalFlexBox) GetColumnCopy(index int) *Column {
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
func (r *HorizontalFlexBox) GetColumnCellCopy(columnIndex, cellIndex int) *Cell {
	if columnIndex >= 0 && columnIndex < len(r.columns) {
		if cellIndex >= 0 && cellIndex < len(r.columns[columnIndex].cells) {
			cellCopy := r.columns[columnIndex].cells[cellIndex].copy()
			return &cellCopy
		}
	}
	return nil
}

// UpdateColumn replaces the FlexBoxColumn on the given index
func (r *HorizontalFlexBox) UpdateColumn(index int, column *Column) *HorizontalFlexBox {
	r.columns[index] = column
	r.setRecalculate()
	return r
}

// LockColumnWidth sets the fixed width value for all the columns
// this will disable horizontal scaling
func (r *HorizontalFlexBox) LockColumnWidth(value int) *HorizontalFlexBox {
	r.fixedColumnWidth = value
	return r
}

// SetHeight sets the FlexBox height
func (r *HorizontalFlexBox) SetHeight(value int) *HorizontalFlexBox {
	r.height = value
	for _, column := range r.columns {
		column.setHeight(value)
	}
	return r
}

// SetWidth sets the FlexBox width
func (r *HorizontalFlexBox) SetWidth(value int) *HorizontalFlexBox {
	r.width = value
	r.setRecalculate()
	return r
}

// GetHeight yields current FlexBox height
func (r *HorizontalFlexBox) GetHeight() int {
	return r.getMaxHeight()
}

// GetWidth yields current FlexBox width
func (r *HorizontalFlexBox) GetWidth() int {
	return r.getMaxWidth()
}

// Render initiates the recalculation of the columns dimensions(width) if the recalculate flag is on,
// and then it renders all the columns and combines them on the horizontal axis
func (r *HorizontalFlexBox) Render() string {
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
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...))
}

// ForceRecalculate forces the recalculation for the box and all the columns
func (r *HorizontalFlexBox) ForceRecalculate() {
	r.recalculate()
	for _, rw := range r.columns {
		rw.recalculate()
	}
}

// recalculate fetches the column height distribution slice and sets it on the columns
func (r *HorizontalFlexBox) recalculate() {
	if r.recalculateFlag {
		if len(r.columns) > 0 {
			r.distributeColumnsDimensions(r.calculateColumnWidth())
		}
		r.unsetRecalculate()
	}
}

func (r *HorizontalFlexBox) setRecalculate() {
	r.recalculateFlag = true
}

func (r *HorizontalFlexBox) unsetRecalculate() {
	r.recalculateFlag = false
}

// calculateColumnWidth calculates the width of each column and returns the distribution array
func (r *HorizontalFlexBox) calculateColumnWidth() (distribution []int) {
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
func (r *HorizontalFlexBox) distributeColumnsDimensions(ratioDistribution []int) {
	for index, column := range r.columns {
		column.setHeight(r.getContentHeight())
		column.setWidth(ratioDistribution[index])
	}
}

// getColumnMatrix return the matrix of the cell widths for all the columns
func (r *HorizontalFlexBox) getColumnMatrix() (columnMatrix [][]int) {
	for _, column := range r.columns {
		var cellValues []int
		for _, cell := range column.cells {
			cellValues = append(cellValues, cell.ratioX)
		}
		columnMatrix = append(columnMatrix, cellValues)
	}
	return columnMatrix
}

func (r *HorizontalFlexBox) getContentWidth() int {
	return r.getMaxWidth() - r.getExtraWidth()
}

func (r *HorizontalFlexBox) getContentHeight() int {
	return r.getMaxHeight() - r.getExtraHeight()
}

func (r *HorizontalFlexBox) getMaxWidth() int {
	return r.width
}

func (r *HorizontalFlexBox) getMaxHeight() int {
	return r.height
}

func (r *HorizontalFlexBox) getExtraWidth() int {
	return r.style.GetHorizontalMargins() + r.style.GetHorizontalBorderSize()
}

func (r *HorizontalFlexBox) getExtraHeight() int {
	return r.style.GetVerticalMargins() + r.style.GetVerticalBorderSize()
}
