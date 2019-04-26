package cli

import "strings"

type TableType struct {
	numRows int
	numCols int
	data    []string
}

func NewTable(numRows, numCols int) *TableType {
	return &TableType{
		numRows: numRows,
		numCols: numCols,
		data:    make([]string, numRows*numCols),
	}
}

func (table *TableType) Format() string {
	fieldLengths := table.getFieldLengths()
	var sb strings.Builder
	for row := 0; row < table.numRows; row++ {
		rowOutput := make([]string, table.numCols)
		for col := 0; col < table.numCols; col++ {
			rowOutput[col] = table.formatField(table.Get(row, col), fieldLengths[col])
		}
		sb.WriteString(strings.Join(rowOutput, " "))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (table *TableType) formatField(value string, fieldLength int) string {
	// Left align
	return value + strings.Repeat(" ", fieldLength-len(value))
}

func (table *TableType) getFieldLengths() []int {
	fieldLengths := make([]int, table.numCols)
	for col := 0; col < table.numCols; col++ {
		fieldLengths[col] = table.getFieldLength(col)
	}
	return fieldLengths
}

func (table *TableType) getFieldLength(col int) int {
	result := 0
	for row := 0; row < table.numRows; row++ {
		l := len(table.Get(row, col))
		if l > result {
			result = l
		}
	}
	return result
}

func (table *TableType) Get(row, col int) string {
	return table.data[table.getIndex(row, col)]
}

func (table *TableType) getIndex(row int, col int) int {
	return table.numCols*row + col
}

func (table *TableType) Set(row, col int, value string) {
	table.data[table.getIndex(row, col)] = value
}
