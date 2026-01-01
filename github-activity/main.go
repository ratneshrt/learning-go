package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ratneshrt/github-activity/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Padding(1, 2).
			Bold(true).
			Render(fmt.Sprintf("Error: %s", err))
		fmt.Println(errorStyle)
	}
}
