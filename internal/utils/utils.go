package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func PrintTable(header []string, rows [][]string) {
	t := table.New()

	t.Headers(header...)
	
	for row := range rows {
		t.Row(rows[row]...)
	}
	fmt.Println(t.Render())
}