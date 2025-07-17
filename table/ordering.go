package table

import (
	"fmt"
	"reflect"
	"strconv"
)

type Ordered interface {
	int | int8 | int32 | int16 | int64 | float32 | float64 | string
}

type SortingOrderKey int

const (
	SortingOrderAscending SortingOrderKey = iota
	SortingOrderDescending
)

// GetOrder returns the current order column index and phase
func (r *Table) GetOrder() (int, SortingOrderKey) { return r.orderedColumnIndex, r.orderedColumnPhase }

// OrderByAsc orders rows by a column with index n, in ascending order
func (r *Table) OrderByAsc(index int) *Table {
	// sanity check first, we won't return errors here, simply ignore if the user sends non-existing index
	if index < len(r.columnHeaders) && len(r.filteredRows) > 1 {
		r.orderedColumnPhase = SortingOrderAscending
		r.rows = sortRows(r.rows, index, r.orderedColumnPhase)
		r.setRowsUpdate()
		r.orderedColumnIndex = index
		r.setHeadersUpdate()
	}
	return r
}

// OrderByDesc orders rows by a column with index n, in descending order
func (r *Table) OrderByDesc(index int) *Table {
	// sanity check first, we won't return errors here, simply ignore if the user sends non existing index
	if index < len(r.columnHeaders) && len(r.filteredRows) > 1 {
		r.orderedColumnPhase = SortingOrderDescending
		r.rows = sortRows(r.rows, index, r.orderedColumnPhase)
		r.setRowsUpdate()
		r.orderedColumnIndex = index
		r.setHeadersUpdate()
	}
	return r
}

// updateOrderedVars updates bits and pieces revolving around ordering
// toggling between asc and desc
// updating ordering vars on TableOrdered
func (r *Table) updateOrderedVars(index int) {
	// toggle between ascending and descending and set default first sort to ascending
	if r.orderedColumnIndex == index {
		switch r.orderedColumnPhase {
		case SortingOrderAscending:
			r.orderedColumnPhase = SortingOrderDescending

		case SortingOrderDescending:
			r.orderedColumnPhase = SortingOrderAscending
		}
	} else {
		r.orderedColumnPhase = SortingOrderDescending
	}
	r.orderedColumnIndex = index

	r.setHeadersUpdate()
}

func sortRows(rows [][]any, index int, orderKey SortingOrderKey) [][]any {
	// sorted rows
	var sorted [][]any
	// list of column values used for ordering
	var orderingCol []any
	for _, rw := range rows {
		orderingCol = append(orderingCol, rw[index])
	}
	// get sorting index
	sortingIndex := sortIndexByOrderedColumn(orderingCol, orderKey)
	// update rows
	for _, i := range sortingIndex {
		sorted = append(sorted, rows[i])
	}
	return sorted
}

// isOrdered check if type is one of valid Ordered types
func isOrdered(e any) bool {
	switch e.(type) {
	case string, int, int8, int16, int32, float32, float64:
		return true
	default:
		return false
	}
}

// getStringFromOrdered returns string from interface that was produced with one of Ordered types
func getStringFromOrdered(i any) string {
	switch i := i.(type) {
	case string:
		return i
	case int:
		return strconv.Itoa(i)
	case int8:
		return strconv.Itoa(int(i))
	case int16:
		return strconv.Itoa(int(i))
	case int32:
		return strconv.Itoa(int(i))
	case int64:
		return strconv.Itoa(int(i))
	case float32:
		// default precision of 24
		return strconv.FormatFloat(float64(i), 'G', 0, 32)
	case float64:
		// default precision of 24
		return strconv.FormatFloat(i, 'G', 0, 64)
	default:
		return ""
	}
}

// sortIndexByOrderedColumn casts to the one of Ordered type that is used on the column and sends to sorting
// returns sorted index of elements rather than elements themselves
func sortIndexByOrderedColumn(i []any, order SortingOrderKey) (sortedIndex []int) {
	// if len of slice is 0 return empty sort order
	if len(i) == 0 {
		return sortedIndex
	}

	switch i[0].(type) {
	case string:
		var s []string
		for _, el := range i {
			s = append(s, el.(string))
		}
		return sortIndex(s, order)
	case int:
		var s []int
		for _, el := range i {
			s = append(s, el.(int))
		}
		return sortIndex(s, order)
	case int8:
		var s []int8
		for _, el := range i {
			s = append(s, el.(int8))
		}
		return sortIndex(s, order)
	case int16:
		var s []int16
		for _, el := range i {
			s = append(s, el.(int16))
		}
		return sortIndex(s, order)
	case int32:
		var s []int32
		for _, el := range i {
			s = append(s, el.(int32))
		}
		return sortIndex(s, order)
	case int64:
		var s []int64
		for _, el := range i {
			s = append(s, el.(int64))
		}
		return sortIndex(s, order)
	case float32:
		var s []float32
		for _, el := range i {
			s = append(s, el.(float32))
		}
		return sortIndex(s, order)
	case float64:
		var s []float64
		for _, el := range i {
			s = append(s, el.(float64))
		}
		return sortIndex(s, order)

	default:
		panic(fmt.Sprintf("type %s not subtype of Ordered", reflect.TypeOf(i[0]).String()))
	}
}

// sortIndex is simple generic bubble sort, returns sorted index slice
// bubble sort implemented for simplicity, if you need faster alg feel free to open a PR for it :zap:
func sortIndex[T Ordered](slice []T, order SortingOrderKey) []int {
	// could do this in sortIndexByOrderedColumn where we cycle through the slice anyhow
	// tho I think this is cheap op and makes code a bit cleaner, worthy trade for now
	var index []int
	for i := 0; i < len(slice); i++ {
		index = append(index, i)
	}

	// bubble sort slice and update index in a process
	for i := len(slice); i > 0; i-- {
		for j := 1; j < i; j++ {
			if order == SortingOrderDescending && slice[j] < slice[j-1] {
				slice[j], slice[j-1] = slice[j-1], slice[j]
				index[j], index[j-1] = index[j-1], index[j]
			} else if order == SortingOrderAscending && slice[j] > slice[j-1] {
				slice[j], slice[j-1] = slice[j-1], slice[j]
				index[j], index[j-1] = index[j-1], index[j]
			}
		}
	}
	return index
}
