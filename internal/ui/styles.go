package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	// Base colors
	ColorPrimary   = lipgloss.Color("#6C8EEF") // Blue
	ColorSecondary = lipgloss.Color("#50C878") // Green
	ColorAccent    = lipgloss.Color("#FF8C00") // Orange
	ColorError     = lipgloss.Color("#FF5252") // Red
	ColorWarning   = lipgloss.Color("#FFD700") // Yellow
	ColorInfo      = lipgloss.Color("#00BFFF") // Light Blue
	ColorSuccess   = lipgloss.Color("#4CAF50") // Green
	ColorMuted     = lipgloss.Color("#9E9E9E") // Gray

	// Text colors
	ColorText      = lipgloss.Color("#FFFFFF") // White
	ColorTextDim   = lipgloss.Color("#CCCCCC") // Light Gray
	ColorTextMuted = lipgloss.Color("#999999") // Medium Gray

	// Background colors
	ColorBackground = lipgloss.Color("#1E1E1E") // Dark Gray
	ColorSurface    = lipgloss.Color("#2D2D2D") // Medium Dark Gray
)

// Styles for different UI elements
var (
	// Title styles
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			MarginBottom(1).Underline(true)

	// Subtitle styles

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	// Text styles
	StyleText = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleTextDim = lipgloss.NewStyle().
			Foreground(ColorTextDim)

	StyleTextMuted = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Italic(true)

	// Status styles
	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	StyleInfo = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// Box styles
	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	StyleErrorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// List styles
	StyleListItem = lipgloss.NewStyle().
			PaddingLeft(2)

	StyleListItemSelected = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorPrimary).
				Bold(true)

	// Command styles
	StyleCommand = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	StyleCommandArg = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Italic(true)

	// Value style
	StyleValue = lipgloss.NewStyle().
			Foreground(ColorSecondary)
)
