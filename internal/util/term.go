package util

import (
	"os"
	"strconv"
)

func TermWidth() int {
	termWidth := os.Getenv("COLUMNS")

	width, err := strconv.Atoi(termWidth)
	if err != nil {
		return 80
	}

	return width
}
