package utils

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/google/uuid"
)

const DEFAULT_TIMEOUT = 30 * time.Second

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
