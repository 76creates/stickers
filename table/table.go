package table

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableDefaultHeaderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#7158e2")).
		Foreground(lipgloss.Color("#ffffff"))
	tableDefaultFooterStyle = tableDefaultHeaderStyle.Align(lipgloss.Right).Height(1)
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

	tableDefaultSortAscChar  = "▲"
	tableDefaultSortDescChar = "▼"
	tableDefaultFilterChar   = "⑂"

	tableDefaultStyles = map[StyleKey]lipgloss.Style{
		StyleKeyHeader:         tableDefaultHeaderStyle,
		StyleKeyFooter:         tableDefaultFooterStyle,
		StyleKeyRows:           tableDefaultRowsStyle,
		StyleKeyRowsSubsequent: tableDefaultRowsSubsequentStyle,
		StyleKeyRowsCursor:     tableDefaultRowsCursorStyle,
		StyleKeyCellCursor:     tableDefaultCellCursorStyle,
	}
)

type StyleKey int

const (
	StyleKeyHeader StyleKey = iota
	StyleKeyFooter
	StyleKeyRows
	StyleKeyRowsSubsequent
	StyleKeyRowsCursor
	StyleKeyCellCursor
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
	columnType    []any
	rows          [][]any

	filteredRows   [][]any
	filteredColumn int
	filterString   string

	// orderColumnIndex notes which column is used for sorting
	// -1 means that no column is sorted
	orderedColumnIndex int
	// orderedColumnPhase remarks if the sort is asc or desc, basically works like a toggle
	// 0 indicates desc sorting, 1 indicates
	orderedColumnPhase SortingOrderKey

	// rowsTopIndex top visible index
	rowsTopIndex int
	cursorIndexY int
	cursorIndexX int

	height int
	width  int

	rowsBoxHeight int
	// rowHeight fixed row height value, maybe this should be optional?
	rowHeight int

	styles map[StyleKey]lipgloss.Style
	// stylePassing if true, styles are passed all the way down from box to cell
	stylePassing bool

	headerBox *flexbox.FlexBox
	rowsBox   *flexbox.FlexBox

	// these flags indicate weather we should update rows and headers flex boxes
	updateRowsFlag    bool
	updateHeadersFlag bool
}

// NewTable initialize Table object with defaults
func NewTable(width, height int, columnHeaders []string) *Table {
	var columnRatio, columnMinWidth []int
	for range columnHeaders {
		columnRatio = append(columnRatio, 1)
		columnMinWidth = append(columnMinWidth, 0)
	}

	// by default all columns are of type string
	var defaultType string
	var defaultTypes []any
	for range columnHeaders {
		defaultTypes = append(defaultTypes, defaultType)
	}

	styles := tableDefaultStyles

	r := &Table{
		columnHeaders:  columnHeaders,
		columnRatio:    columnRatio,
		columnMinWidth: columnMinWidth,
		cursorIndexX:   0,
		cursorIndexY:   0,

		columnType:         defaultTypes,
		orderedColumnIndex: -1,
		orderedColumnPhase: SortingOrderDescending,

		filteredColumn: -1,
		filterString:   "",

		height: height,
		width:  width,
		// when optional header/footer is set rework this
		rowsBoxHeight: height - 2,

		rowsTopIndex: 0,
		rowHeight:    1,

		headerBox: flexbox.New(width, 1).SetStyle(tableDefaultHeaderStyle),
		rowsBox:   flexbox.New(width, height-1),

		styles:       styles,
		stylePassing: false,
	}
	r.setHeadersUpdate()
	return r
}

// SetRatio replaces the ratio slice, it has to be exactly the len of the headers/rows slices
// also each value have to be greater than 0, if either fails we panic
func (r *Table) SetRatio(values []int) *Table {
	if len(values) != len(r.columnHeaders) {
		log.Fatalf("ratio list[%d] not of proper length[%d]\n", len(values), len(r.columnHeaders))
	}
	for _, val := range values {
		if val < 1 {
			log.Fatalf("ratio value must be greater than 0")
		}
	}
	r.columnRatio = values
	r.setHeadersUpdate()
	r.setRowsUpdate()
	return r
}

// SetTypes sets the column type, setting this will remove all the rows so make sure you do it when instantiating
// Table object or add new rows after this, types have to be one of Ordered interface types
func (r *Table) SetTypes(columnTypes ...any) (*Table, error) {
	if len(columnTypes) != len(r.columnHeaders) {
		return r, errors.New("column types not the same len as headers")
	}
	for i, t := range columnTypes {
		if !isOrdered(t) {
			message := fmt.Sprintf(
				"column of type %s on index %d is not of type Ordered", reflect.TypeOf(t).String(), i,
			)
			return r, ErrorBadType{msg: message}
		}
	}
	r.cursorIndexY, r.cursorIndexX = 0, 0
	r.rows = [][]any{}
	r.columnType = columnTypes
	r.setRowsUpdate()
	return r, nil
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
	r.setTopRow()
	return r
}

// SetWidth sets the width of the table
func (r *Table) SetWidth(value int) *Table {
	r.width = value
	r.rowsBox.SetWidth(value)
	r.headerBox.SetWidth(value)
	return r
}

// SetStyles allows overrides of styling elements of the table
// When only a partial set of overrides are provided, the default styling will be used
func (r *Table) SetStyles(styles map[StyleKey]lipgloss.Style) *Table {
	mergedStyles := tableDefaultStyles
	for key, style := range styles {
		mergedStyles[key] = style
	}
	r.styles = mergedStyles
	r.setRowsUpdate()
	r.setHeadersUpdate()
	return r
}

// SetStylePassing sets the style passing flag, if true, styles are passed all the way down from box to cell
func (r *Table) SetStylePassing(value bool) *Table {
	r.stylePassing = value
	r.headerBox.StylePassing(value)
	r.rowsBox.StylePassing(value)
	r.setRowsUpdate()
	r.setHeadersUpdate()
	return r
}

// UnsetFilter resets filtering
func (r *Table) UnsetFilter() *Table {
	r.filterString = ""
	r.filteredColumn = -1
	r.setTopRow()
	r.setRowsUpdate()
	r.setHeadersUpdate()
	return r
}

// SetFilter sets filtering string on a column
func (r *Table) SetFilter(columnIndex int, s string) *Table {
	if columnIndex < len(r.columnHeaders) {
		r.filterString = s
		r.filteredColumn = columnIndex

		r.setRowsUpdate()
	}
	return r
}

// GetFilter returns string used for filtering and the column index
// TODO: enable multi column filtering
func (r *Table) GetFilter() (columnIndex int, s string) {
	return r.filteredColumn, r.filterString
}

// CursorDown move table cursor down
func (r *Table) CursorDown() *Table {
	if r.cursorIndexY+1 < len(r.filteredRows) {
		r.cursorIndexY++
		r.setTopRow()
		r.setRowsUpdate()
	}
	return r
}

// CursorUp move table cursor up
func (r *Table) CursorUp() *Table {
	if r.cursorIndexY-1 > -1 {
		r.cursorIndexY--
		r.setTopRow()
		r.setRowsUpdate()
	}
	return r
}

// CursorLeft move table cursor left
func (r *Table) CursorLeft() *Table {
	if r.cursorIndexX-1 > -1 {
		r.cursorIndexX--
		// TODO: update row only
		r.setRowsUpdate()
	}
	return r
}

// CursorRight move table cursor right
func (r *Table) CursorRight() *Table {
	if r.cursorIndexX+1 < len(r.columnHeaders) {
		r.cursorIndexX++
		// TODO: update row only
		r.setRowsUpdate()
	}
	return r
}

// GetCursorLocation returns the current x,y position of the cursor
func (r *Table) GetCursorLocation() (int, int) {
	return r.cursorIndexX, r.cursorIndexY
}

// GetCursorValue returns the string of the cell under the cursor
func (r *Table) GetCursorValue() string {
	// handle 0 rows situation and when table is not active
	if len(r.filteredRows) == 0 || r.cursorIndexX < 0 || r.cursorIndexY < 0 {
		return ""
	}
	return getStringFromOrdered(r.filteredRows[r.cursorIndexY][r.cursorIndexX])
}

// AddRows add multiple rows, will return error on the first instance of a row that does not match the type set on table
// will update rows only when there are no errors
func (r *Table) AddRows(rows [][]any) (*Table, error) {
	// check for errors
	for _, row := range rows {
		if err := r.validateRow(row...); err != nil {
			return r, err
		}
	}
	// append rows
	r.rows = append(r.rows, rows...)

	r.applyFilter()
	r.setRowsUpdate()
	return r, nil
}

// MustAddRows executes AddRows and panics if there is an error
func (r *Table) MustAddRows(rows [][]any) *Table {
	if _, err := r.AddRows(rows); err != nil {
		panic(err)
	}
	return r
}

// ClearRows removes all previously added rows, can be used as part of an update loop
func (r *Table) ClearRows() *Table {
	r.rows = make([][]any, 0, 10)
	r.setRowsUpdate()
	return r
}

// Render renders the table into the string
func (r *Table) Render() string {
	r.updateRows()
	r.updateHeader()

	statusMessage := fmt.Sprintf(
		"%d:%d / %d:%d ",
		r.cursorIndexX,
		r.cursorIndexY,
		r.rowsBox.GetWidth(),
		r.rowsBox.GetHeight(),
	)
	if r.cursorIndexX == r.filteredColumn {
		statusMessage = fmt.Sprintf("filtered by: %q / %s", r.filterString, statusMessage)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		r.headerBox.Render(),
		r.rowsBox.Render(),
		r.styles[StyleKeyFooter].Width(r.width).Render(statusMessage),
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

// validateRow checks the row for validity, number of cells must match table header length
// and header types per cell as well
func (r *Table) validateRow(cells ...any) error {
	var message string
	// check row len
	if len(cells) != len(r.columnType) {
		message = fmt.Sprintf(
			"len of row[%d] does not equal number of columns[%d]", len(cells), len(r.columnType),
		)
		return ErrorRowLen{msg: message}
	}
	// check cell type
	for i, c := range cells {
		switch c.(type) {
		case string, int, int8, int16, int32, float32, float64:
			// check if the cell matches the type of the column
			if reflect.TypeOf(c) != reflect.TypeOf(r.columnType[i]) {
				message = fmt.Sprintf(
					"type of the cell[%v] on index %d not matching type of the column[%v]",
					reflect.TypeOf(c), i, reflect.TypeOf(r.columnType[i]),
				)
				return ErrorBadCellType{msg: message}
			}
		default:
			message = fmt.Sprintf(
				"type[%v] on index %d not matching Ordered interface types", reflect.TypeOf(c), i,
			)
			return ErrorBadType{msg: message}
		}
	}
	return nil
}

// updateHeader recomputes the header of the table
func (r *Table) updateHeader() *Table {
	if !r.updateHeadersFlag {
		return r
	}
	var cells []*flexbox.Cell
	r.headerBox.SetStyle(r.styles[StyleKeyHeader])
	for i, title := range r.columnHeaders {
		// titleSuffix at the moment can be sort and filter characters
		// filtering symbol should be visible always, if possible of course, and as far right as possible
		// there should be a minimum of space bar between two symbols and symbol and row to the right
		var titleSuffix string
		_h := r.headerBox.GetRow(0)
		// skip the case when we initialize table
		if _h != nil {
			_c := _h.GetCellCopy(i)
			if _c == nil {
				panic("cell with index " + strconv.Itoa(i) + " is nil")
			}
			_w := _c.GetWidth()

			// add sorting symbol if the sorting is active on the column
			if r.orderedColumnIndex == i {
				if r.orderedColumnPhase == SortingOrderDescending {
					titleSuffix = " " + tableDefaultSortDescChar
				} else if r.orderedColumnPhase == SortingOrderAscending {
					titleSuffix = " " + tableDefaultSortAscChar
				}
			}

			// add filtering symbol if the filtering is active on the column
			if r.filteredColumn == i && r.filterString != "" {
				// add at least one space bar between char to the left, and one to the right
				titleSuffix = titleSuffix + strings.Repeat(
					" ", int(math.Max(
						1,
						float64(
							_w-utf8.RuneCountInString(title+titleSuffix)-2,
						),
					)),
				) + tableDefaultFilterChar + " "
			}

			// if title and suffix exceed width trim the title
			if _w-utf8.RuneCountInString(title+titleSuffix) < 0 {
				// this will be the cae only when sort is on and filter is off
				// add one space bar between sort and column to the right
				if utf8.RuneCountInString(titleSuffix) == 2 {
					titleSuffix = titleSuffix + " "
				}
				// trim the title
				title = title[0:int(math.Max(0, float64(_w-utf8.RuneCountInString(titleSuffix))))]
			}
		}
		cells = append(
			cells,
			flexbox.NewCell(r.columnRatio[i], 1).SetMinWidth(r.columnMinWidth[i]).SetContent(title+titleSuffix),
		)
	}
	r.headerBox.SetRows(
		[]*flexbox.Row{
			r.headerBox.NewRow().StylePassing(r.stylePassing).AddCells(cells...),
		},
	)
	r.unsetHeadersUpdate()
	return r
}

// updateRows recomputes the rows of the table
// calculate the visible rows top/bottom indexes
// create rows and their cells with styles depending on state
func (r *Table) updateRows() {
	if !r.updateRowsFlag {
		return
	}
	if r.rowsBoxHeight < 0 {
		r.unsetRowsUpdate()
		return
	}
	r.applyFilter()

	// calculate the bottom most visible row index
	rowsBottomIndex := r.rowsTopIndex + r.rowsBoxHeight
	if rowsBottomIndex > len(r.filteredRows) {
		rowsBottomIndex = len(r.filteredRows)
	}

	var rows []*flexbox.Row
	for ir, columns := range r.filteredRows[r.rowsTopIndex:rowsBottomIndex] {
		// irCorrected is corrected row index since we iterate only visible rows
		irCorrected := ir + r.rowsTopIndex

		var cells []*flexbox.Cell
		for ic, column := range columns {
			// initialize column cell
			c := flexbox.NewCell(r.columnRatio[ic], r.rowHeight).
				SetMinWidth(r.columnMinWidth[ic]).
				SetContent(getStringFromOrdered(column))
			// update style if cursor is on the cell, otherwise it's inherited from the row
			if irCorrected == r.cursorIndexY && ic == r.cursorIndexX {
				c.SetStyle(r.styles[StyleKeyCellCursor])
			}
			cells = append(cells, c)
		}
		// initialize new row from the rows box and add generated cells
		rw := r.rowsBox.NewRow().StylePassing(r.stylePassing).AddCells(cells...)

		// rows have three styles, normal, subsequent and selected
		// normal and subsequent rows should differ for readability
		// TODO: make this ^ optional
		if irCorrected == r.cursorIndexY {
			rw.SetStyle(r.styles[StyleKeyRowsCursor])
		} else if irCorrected%2 == 0 || irCorrected == 0 {
			rw.SetStyle(r.styles[StyleKeyRowsSubsequent])
		} else {
			rw.SetStyle(r.styles[StyleKeyRows])
		}

		rows = append(rows, rw)
	}

	// lock row height, this might get optional at some point
	r.rowsBox.LockRowHeight(r.rowHeight)
	r.rowsBox.SetRows(rows)
	r.unsetRowsUpdate()
}

// applyFilter filters column n by a value s
func (r *Table) applyFilter() *Table {
	// sending empty string should reset the filtering
	if r.filterString == "" {
		r.filteredRows = r.rows
		return r
	}
	var filteredRows [][]any
	for _, row := range r.rows {
		cellValue := getStringFromOrdered(row[r.filteredColumn])
		// convert to lower, not sure if anybody needs case-sensitive filtering
		// if you are reading this and need it, open up an issue :zap:
		if strings.Contains(strings.ToLower(cellValue), strings.ToLower(r.filterString)) {
			filteredRows = append(filteredRows, row)
		}
	}
	r.filteredRows = filteredRows
	r.setTopRow()
	r.setHeadersUpdate()
	return r
}

// setTopRow calculates the row top index used when deciding what is visible
func (r *Table) setTopRow() {
	// if rows are empty set y to 0, retain x pos
	// will be useful for filtering
	if len(r.filteredRows) == 0 {
		r.cursorIndexY = 0
	} else if r.cursorIndexY > len(r.filteredRows) {
		// when filtering if cursor is higher than row length
		// set it to the bottom of the list
		r.cursorIndexY = len(r.filteredRows) - 1
	}

	// case when cursor is in between top or bottom visible row
	if r.cursorIndexY >= r.rowsTopIndex && r.cursorIndexY < r.rowsTopIndex+r.rowsBoxHeight {
		// if cursor is on the last item in row, adjust the row top
		if r.cursorIndexY == len(r.filteredRows)-1 {
			// if all rows can fit on screen
			if len(r.filteredRows) <= r.rowsBoxHeight {
				r.rowsTopIndex = 0
				return
			}
			// fit max rows on the table
			r.rowsTopIndex = r.cursorIndexY - (r.rowsBoxHeight - 1)
		} else if r.cursorIndexY > len(r.filteredRows)-1 && len(r.filteredRows) != 0 {
			r.cursorIndexY = len(r.filteredRows) - 1
		}
		return
	}

	// if cursor is above the top
	if r.cursorIndexY < r.rowsTopIndex {
		if r.cursorIndexY == len(r.filteredRows)-1 {
			// if all rows can fit on screen
			if len(r.filteredRows) <= r.rowsBoxHeight {
				r.rowsTopIndex = 0
				return
			}
			// fit max rows on the table
			r.rowsTopIndex = r.cursorIndexY - (r.rowsBoxHeight - 1)
			return
		}
		r.rowsTopIndex = r.cursorIndexY
		return
	}

	if r.cursorIndexY > r.rowsTopIndex {
		r.rowsTopIndex = r.cursorIndexY - r.rowsBoxHeight + 1
		return
	}
}
