package stickers

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
)

var (
	tableDefaultHeaderStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7158e2")).
				Foreground(lipgloss.Color("#ffffff"))
	tableDefaultFooterStyle = tableDefaultHeaderStyle.Copy().Align(lipgloss.Right).Height(1)
	tableDefaultRowsStyle   = lipgloss.NewStyle().
				Background(lipgloss.Color("#4b4b4b")).
				Foreground(lipgloss.Color("#ffffff"))
	tableDefaultRowsSubsequentStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#3d3d3d")).
					Foreground(lipgloss.Color("#ffffff"))
	tableDefaultRowsCursorStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#f7b731")).
					Foreground(lipgloss.Color("#000000")).
					Bold(true)
	tableDefaultCellCursorStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("#f6e58d")).
					Foreground(lipgloss.Color("#000000"))
)

type tableStyleKey int

const (
	TableHeaderStyleKey tableStyleKey = iota
	TableFooterStyleKey
	TableRowsStyleKey
	TableRowsSubsequentStyleKey
	TableRowsCursorStyleKey
	TableCellCursorStyleKey
)

// Table responsive, x/y scrollable table that uses magic of FlexBox
type Table struct {
	// columnRatio ratio of the columns, is applied to rows as well
	columnRatio []int
	// columnMinWidth minimal width of the column
	columnMinWidth []int
	// columnHeaders column text headers
	// TODO: make this optional, as well as footer
	columnHeaders []string
	rows          [][]string

	// rowsTopIndex top visible index
	rowsTopIndex int
	cursorIndexY int
	cursorIndexX int

	height int
	width  int

	rowsBoxHeight int
	// rowHeight fixed row height value, maybe this should be optional?
	rowHeight int

	styles map[tableStyleKey]lipgloss.Style

	headerBox *FlexBox
	rowsBox   *FlexBox

	// these flags indicate weather we should update rows and headers flex boxes
	updateRowsFlag    bool
	updateHeadersFlag bool
}

// NewTable initialize Table object with defaults
func NewTable(width, height int, columnHeaders []string) *Table {
	var columnRatio, columnMinWidth []int
	for _ = range columnHeaders {
		columnRatio = append(columnRatio, 1)
		columnMinWidth = append(columnMinWidth, 0)
	}

	styles := map[tableStyleKey]lipgloss.Style{
		TableHeaderStyleKey:         tableDefaultHeaderStyle,
		TableFooterStyleKey:         tableDefaultFooterStyle,
		TableRowsStyleKey:           tableDefaultRowsStyle,
		TableRowsSubsequentStyleKey: tableDefaultRowsSubsequentStyle,
		TableRowsCursorStyleKey:     tableDefaultRowsCursorStyle,
		TableCellCursorStyleKey:     tableDefaultCellCursorStyle,
	}

	r := &Table{
		columnHeaders:  columnHeaders,
		columnRatio:    columnRatio,
		columnMinWidth: columnMinWidth,
		cursorIndexX:   0,
		cursorIndexY:   0,

		height: height,
		// when optional header/footer is set rework this
		rowsBoxHeight: height - 2,
		width:         width,

		rowsTopIndex: 0,
		rowHeight:    1,

		headerBox: NewFlexBox(width, 1).SetStyle(tableDefaultHeaderStyle),
		rowsBox:   NewFlexBox(width, height-1),

		styles: styles,
	}
	r.updateHeader()
	return r
}

// SetRatio replaces the ratio slice, it has to be exactly the len of the headers/rows slices
// if it's not matching len it will trigger fatal error
func (r *Table) SetRatio(values []int) *Table {
	if len(values) != len(r.columnHeaders) {
		log.Fatalf("ratio list[%d] not of proper length[%d]\n", len(values), len(r.columnHeaders))
	}
	r.columnRatio = values
	r.setHeadersUpdate()
	r.setRowsUpdate()
	return r
}

// SetMinWidth replaces the minimum width slice, it has to be exactly the len of the headers/rows slices
// if it's not matching len it will trigger fatal error
func (r *Table) SetMinWidth(values []int) *Table {
	if len(values) != len(r.columnHeaders) {
		log.Fatalf("min width list[%d] not of proper length[%d]\n", len(values), len(r.columnHeaders))
	}
	r.columnMinWidth = values
	r.setHeadersUpdate()
	r.setRowsUpdate()
	return r
}

// SetHeight sets the height of the table including the header and footer
func (r *Table) SetHeight(value int) *Table {
	r.height = value
	// we deduct two to take header/footer into the account
	r.rowsBoxHeight = value - 2
	r.rowsBox.SetHeight(r.rowsBoxHeight)
	r.setRowsUpdate()
	return r
}

// SetWidth sets the width of the table
func (r *Table) SetWidth(value int) *Table {
	r.width = value
	r.rowsBox.SetWidth(value)
	r.headerBox.SetWidth(value)
	return r
}

// CursorDown move table cursor down
func (r *Table) CursorDown() *Table {
	if r.cursorIndexY+1 < len(r.rows) {
		r.cursorIndexY++
	}
	r.setTopRow()
	r.setRowsUpdate()
	return r
}

// CursorUp move table cursor up
func (r *Table) CursorUp() *Table {
	if r.cursorIndexY-1 > -1 {
		r.cursorIndexY--
	}
	r.setTopRow()
	r.setRowsUpdate()
	return r
}

// CursorLeft move table cursor left
func (r *Table) CursorLeft() *Table {
	if r.cursorIndexX-1 > -1 {
		r.cursorIndexX--
	}
	// TODO: update row only
	r.setRowsUpdate()
	return r
}

// CursorRight move table cursor right
func (r *Table) CursorRight() *Table {
	if r.cursorIndexX+1 < len(r.columnHeaders) {
		r.cursorIndexX++
	}
	// TODO: update row only
	r.setRowsUpdate()
	return r
}

// GetCursorValue returns the string of the cell under the cursor
func (r *Table) GetCursorValue() string {
	return r.rows[r.cursorIndexY][r.cursorIndexX]
}

func (r *Table) AddRows(rows [][]string) *Table {
	for _, row := range rows {
		if len(row) != len(r.columnHeaders) {
			log.Fatal("row length doesn't match column length")
		}
		r.rows = append(r.rows, row)
	}
	r.updateRows()
	return r
}

// Render renders the table into the string
func (r *Table) Render() string {
	r.updateRows()
	r.updateHeader()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		r.headerBox.Render(),
		r.rowsBox.Render(),
		r.styles[TableFooterStyleKey].
			Width(r.width).
			Render(
				fmt.Sprintf(
					"%d:%d / %d:%d ",
					r.cursorIndexX,
					r.cursorIndexY,
					r.rowsBox.GetWidth(),
					r.rowsBox.GetHeight(),
				),
			),
	)
}

func (r *Table) setRowsUpdate() {
	r.updateRowsFlag = true
}

func (r *Table) unsetRowsUpdate() {
	r.updateRowsFlag = false
}

func (r *Table) setHeadersUpdate() {
	r.updateHeadersFlag = true
}

func (r *Table) unsetHeadersUpdate() {
	r.updateHeadersFlag = false
}

// updateHeader recomputes the header of the table
func (r *Table) updateHeader() *Table {
	if !r.updateHeadersFlag {
		return r
	}
	var cells []*FlexBoxCell
	r.headerBox.SetStyle(r.styles[TableHeaderStyleKey])
	for i, title := range r.columnHeaders {
		cells = append(
			cells,
			NewFlexBoxCell(r.columnRatio[i], 1).SetMinWidth(r.columnMinWidth[i]).SetContent(title),
		)
	}
	r.headerBox.SetRows(
		[]*FlexBoxRow{
			r.headerBox.NewRow().AddCells(cells),
		},
	)
	r.unsetHeadersUpdate()
	return r
}

// updateRows recomputes the rows of the table
// calculate the visible rows top/bottom indexes
// create rows and their cells with styles depending on state
func (r *Table) updateRows() *Table {
	if !r.updateRowsFlag {
		return r
	}
	if r.rowsBoxHeight < 0 {
		r.unsetRowsUpdate()
		return r
	}

	// calculate the bottom most visible row index
	rowsBottomIndex := r.rowsTopIndex + r.rowsBoxHeight
	if rowsBottomIndex > len(r.rows) {
		rowsBottomIndex = len(r.rows)
	}

	var rows []*FlexBoxRow
	for ir, columns := range r.rows[r.rowsTopIndex:rowsBottomIndex] {
		// irCorrected is corrected row index since we iterate only visible rows
		irCorrected := ir + r.rowsTopIndex

		var cells []*FlexBoxCell
		for ic, column := range columns {
			// initialize column cell
			c := NewFlexBoxCell(r.columnRatio[ic], r.rowHeight).
				SetMinWidth(r.columnMinWidth[ic]).
				SetContent(column)
			// update style if cursor is on the cell, otherwise it's inherited from the row
			if irCorrected == r.cursorIndexY && ic == r.cursorIndexX {
				c.SetStyle(r.styles[TableCellCursorStyleKey])
			}
			cells = append(cells, c)
		}
		// initialize new row from the rows box and add generated cells
		rw := r.rowsBox.NewRow().AddCells(cells)

		// rows have three styles, normal, subsequent and selected
		// normal and subsequent rows should differ for readability
		// TODO: make this ^ optional
		if irCorrected == r.cursorIndexY {
			rw.SetStyle(r.styles[TableRowsCursorStyleKey])
		} else if irCorrected%2 == 0 || irCorrected == 0 {
			rw.SetStyle(r.styles[TableRowsSubsequentStyleKey])
		} else {
			rw.SetStyle(r.styles[TableRowsStyleKey])
		}

		rows = append(rows, rw)
	}

	// lock row height, this might get optional at some point
	r.rowsBox.LockRowHeight(r.rowHeight)
	r.rowsBox.SetRows(rows)
	r.unsetRowsUpdate()
	return r
}

// setTopRow calculates the row top index used when deciding whats visable
func (r *Table) setTopRow() {
	// if rows are empty set x and y to 0
	// will be useful for filtering
	if len(r.rows) == 0 {
		r.cursorIndexY = 0
		r.cursorIndexX = 0
	} else if r.cursorIndexY > len(r.rows) {
		// when filtering if cursor is higher than row length
		// set it to the bottom of the list
		r.cursorIndexY = len(r.rows) - 1
	}

	// case when cursor is in between top or bottom visible row
	if r.cursorIndexY >= r.rowsTopIndex && r.cursorIndexY < r.rowsTopIndex+r.rowsBoxHeight {
		// if cursor is on the last item in row, adjust the row top
		if r.cursorIndexY == len(r.rows)-1 {
			// if all rows can fit on screen
			if len(r.rows) <= r.rowsBoxHeight {
				r.rowsTopIndex = 0
				return
			}
			// fit max rows on the table
			r.rowsTopIndex = r.cursorIndexY - (r.rowsBoxHeight - 1)
		}
		return
	}

	// if cursor is above the top
	if r.cursorIndexY < r.rowsTopIndex {
		r.rowsTopIndex = r.cursorIndexY
		return
	}

	if r.cursorIndexY > r.rowsTopIndex {
		//log.Fatal(fmt.Sprintf("[%d][%d][%d]", r.rowsTopIndex, r.cursorIndexY, r.rowsBoxHeight))
		r.rowsTopIndex = r.cursorIndexY - r.rowsBoxHeight + 1
		return
	}
}
