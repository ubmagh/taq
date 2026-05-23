package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorPrefix = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1")).Render("[error]")
	warnPrefix  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3")).Render("[warn]")
	infoPrefix  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Render("[info]")
)

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", errorPrefix, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", warnPrefix, fmt.Sprintf(format, args...))
}

func Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "%s %s\n", infoPrefix, fmt.Sprintf(format, args...))
}
