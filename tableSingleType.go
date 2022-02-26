package stickers

// TableSingleType is Table that is using only 1 type for rows allowing for easier AddRows with fewer errors
type TableSingleType[T Ordered] struct {
	Table
}

// NewTableSingleType initialize TableSingleType object with defaults
func NewTableSingleType[T Ordered](width, height int, columnHeaders []string) *TableSingleType[T] {
	var defaultTypes []any
	var usedType T

	// set type to selected type
	for range columnHeaders {
		defaultTypes = append(defaultTypes, usedType)
	}

	t := &TableSingleType[T]{
		Table: *NewTable(width, height, columnHeaders),
	}

	_, err := t.Table.SetTypes(defaultTypes...)
	if err != nil {
		panic(err)
	}

	return t
}

// SetTypes overridden for TableSimple
func (r *TableSingleType[T]) SetTypes() {
}

func (r *TableSingleType[T]) AddRows(rows [][]T) *TableSingleType[T] {
	for _, row := range rows {
		var _row []any
		for _, cell := range row {
			_row = append(_row, cell)
		}
		r.rows = append(r.rows, _row)
	}
	r.setRowsUpdate()
	return r
}

func (r *TableSingleType[T]) MustAddRows(rows [][]T) *TableSingleType[T] {
	return r.AddRows(rows)
}
