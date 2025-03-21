package ui

import (
	"fmt"

	"github.com/bayou-brogrammer/mygo/internal/logger"
	"github.com/charmbracelet/lipgloss"
)

// PrintTitle prints a formatted title
func PrintTitle(title string) {
	fmt.Println(StyleTitle.Render(title))
}

// PrintSubtitle prints a formatted subtitle
func PrintSubtitle(subtitle string) {
	fmt.Println(StyleSubtitle.Render(subtitle))
}

// PrintInfo prints formatted information text
func PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(StyleInfo.Render(message))
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(StyleSuccess.Render(message))
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(StyleError.Render(message))
	logger.Error("%s", message)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(StyleWarning.Render(message))
	logger.Warn("%s", message)
}

// PrintCommand prints a command with its arguments
func PrintCommand(cmd string, args ...interface{}) {
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(cmd, args...)
	} else {
		message = cmd
	}
	fmt.Println(StyleCommand.Render(message))
}

// PrintList prints a list of items
func PrintList(items []string, selected int) {
	for i, item := range items {
		if i == selected {
			fmt.Println(StyleListItemSelected.Render(item))
		} else {
			fmt.Println(StyleListItem.Render(item))
		}
	}
}

// PrintBox prints content in a styled box
func PrintBox(content string) {
	fmt.Println(StyleBox.Render(content))
}

// PrintErrorBox prints error content in a styled box
func PrintErrorBox(content string) {
	fmt.Println(StyleErrorBox.Render(content))
}

// FormatKey formats a key string (like a configuration key)
func FormatKey(key string) string {
	return StyleCommand.Render(key)
}

// FormatValue formats a value string (like a configuration value)
func FormatValue(value string) string {
	return StyleTextDim.Render(value)
}

// FormatKeyValue formats a key-value pair
func FormatKeyValue(key, value string) string {
	return FormatKey(key) + ": " + FormatValue(value)
}

// FormatCommand formats a command name with accent color for better visibility
func FormatCommand(cmd string) string {
	return StyleCommand.Render(cmd)
}

// FormatAnyTextWithColor formats any text with a specified color
func FormatTextWithColor(text string, style *lipgloss.Style, color lipgloss.Color) string {
	return style.Foreground(color).Render(text)
}
