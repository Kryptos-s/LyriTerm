package ui

import (
	"github.com/kryptos-s/lyriterm/internal/config"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Header   lipgloss.Style
	Status   lipgloss.Style
	Active   lipgloss.Style
	Inactive lipgloss.Style
	Dim      lipgloss.Style

	Sung     lipgloss.Style
	Future   lipgloss.Style
	Title    lipgloss.Style
	Selected lipgloss.Style
	Normal   lipgloss.Style
	Help     lipgloss.Style
	Border   lipgloss.Style
}

func InitStyles(t config.Theme, width int) Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Bg)).
			Background(lipgloss.Color(t.Primary)).
			Width(width).
			Align(lipgloss.Center),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Dim)).
			Align(lipgloss.Center).
			Width(width),

		Active: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Primary)).
			Align(lipgloss.Center).
			Width(width),

		Inactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)).
			Align(lipgloss.Center).
			Width(width),

		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Dim)).
			Align(lipgloss.Center).
			Width(width),

		Sung: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Primary)).
			Bold(true),

		Future: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Text)),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Bg)).
			Background(lipgloss.Color(t.Primary)).
			Padding(0, 1).
			Align(lipgloss.Center).
			Width(50),

		Selected: lipgloss.NewStyle().
			Width(50).
			Padding(0, 2).
			Background(lipgloss.Color(t.Dim)).
			Foreground(lipgloss.Color(t.Text)).
			Bold(true),

		Normal: lipgloss.NewStyle().
			Width(50).
			Padding(0, 2),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Dim)).
			Align(lipgloss.Center).
			Width(50).
			MarginTop(2),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.Primary)),
	}
}
