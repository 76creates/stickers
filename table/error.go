package table

// ErrorBadType type does not match Ordered interface types
type ErrorBadType struct {
	msg string
}

func (e ErrorBadType) Error() string {
	return e.msg
}

// ErrorRowLen row length is not matching headers len
type ErrorRowLen struct {
	msg string
}

func (e ErrorRowLen) Error() string {
	return e.msg
}

// ErrorBadCellType type of cell does not match type of column
type ErrorBadCellType struct {
	msg string
}

func (e ErrorBadCellType) Error() string {
	return e.msg
}
