package ui

import (
	"fmt"
	"lyriterm/internal/config"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingItem int

const (
	SettingTheme SettingItem = iota
	SettingOffset
	SettingKaraoke
	SettingDim
	SettingProgress
	SettingBack
)

var settingItems = []SettingItem{
	SettingTheme,
	SettingOffset,
	SettingKaraoke,
	SettingDim,
	SettingProgress,
	SettingBack,
}

func (s SettingItem) String() string {
	switch s {
	case SettingTheme:
		return "Theme"
	case SettingOffset:
		return "Global Offset (ms)"
	case SettingKaraoke:
		return "Karaoke Mode"
	case SettingDim:
		return "Dim Inactive Lines"
	case SettingProgress:
		return "Show Progress Bar"
	case SettingBack:
		return "Save & Back"
	}
	return "?"
}

func (m *Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.viewState = ViewMain
			return m, nil
		case "up", "k":
			m.settingsIdx--
			if m.settingsIdx < 0 {
				m.settingsIdx = len(settingItems) - 1
			}
		case "down", "j":
			m.settingsIdx++
			if m.settingsIdx >= len(settingItems) {
				m.settingsIdx = 0
			}
		case "enter", "space":
			if settingItems[m.settingsIdx] == SettingBack {
				m.viewState = ViewMain
				m.configMgr.Save(m.Config)
			} else {
				m.toggleSetting(settingItems[m.settingsIdx])
			}
		case "left", "h":
			m.adjustSetting(settingItems[m.settingsIdx], -1)
		case "right", "l":
			m.adjustSetting(settingItems[m.settingsIdx], 1)
		}
	}
	return m, nil
}

func (m *Model) toggleSetting(item SettingItem) {
	switch item {
	case SettingKaraoke:
		m.Config.KaraokeMode = !m.Config.KaraokeMode
	case SettingDim:
		m.Config.DimInactive = !m.Config.DimInactive
	case SettingProgress:
		m.Config.ShowProgress = !m.Config.ShowProgress
	}
}

func (m *Model) adjustSetting(item SettingItem, dir int) {
	switch item {
	case SettingTheme:
		curr := 0
		for i, t := range config.Themes {
			if t.Name == m.Config.ThemeName {
				curr = i
				break
			}
		}

		curr += dir
		if curr < 0 {
			curr = len(config.Themes) - 1
		} else if curr >= len(config.Themes) {
			curr = 0
		}

		m.Config.ThemeName = config.Themes[curr].Name
		m.Theme = config.Themes[curr]

		m.styles = InitStyles(m.Theme, m.width)

	case SettingOffset:
		m.Config.OffsetMs += (dir * 50)
	}
}

func (m *Model) settingsView() string {
	var sb strings.Builder

	sb.WriteString("\n" + m.styles.Title.Render("SETTINGS") + "\n\n")

	for i, item := range settingItems {
		label := item.String()
		value := ""

		switch item {
		case SettingTheme:
			value = fmt.Sprintf("< %s >", strings.ToUpper(m.Config.ThemeName))
		case SettingOffset:
			value = fmt.Sprintf("< %d >", m.Config.OffsetMs)
		case SettingKaraoke:
			value = boolStr(m.Config.KaraokeMode)
		case SettingDim:
			value = boolStr(m.Config.DimInactive)
		case SettingProgress:
			value = boolStr(m.Config.ShowProgress)
		case SettingBack:
			value = ">>"
		}

		line := fmt.Sprintf("%-30s %s", label, value)

		if i == m.settingsIdx {
			sb.WriteString(m.styles.Selected.Render(line) + "\n")
		} else {
			sb.WriteString(m.styles.Normal.Render(line) + "\n")
		}
	}

	sb.WriteString(m.styles.Help.Render("Use Arrows/Enter to change. ESC to close."))

	dialog := m.styles.Border.Render(sb.String())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}

func boolStr(b bool) string {
	if b {
		return "[ON]"
	}
	return "[OFF]"
}
