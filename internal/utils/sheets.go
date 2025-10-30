package utils

import "strconv"

// SheetsConvertColumnToLetter example column 1 = 'A'
func SheetsConvertColumnToLetter(col int) string {
	col--
	letters := ""
	for col >= 0 {
		letters = string(rune('A'+(col%26))) + letters
		col = (col / 26) - 1
	}
	return letters
}

// SheetsConvertCoordToPosition [1,2] = 'A2'
func SheetsConvertCoordToPosition(x, y int) string {
	return SheetsConvertColumnToLetter(x) + strconv.Itoa(y)
}
