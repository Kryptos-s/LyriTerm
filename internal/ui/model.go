package ui

import (
	"context"
	"fmt"
	"github.com/kryptos-s/lyriterm/internal/config"
	"github.com/kryptos-s/lyriterm/internal/lyrics"
	"github.com/kryptos-s/lyriterm/internal/player"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewState int

const (
	ViewMain ViewState = iota
	ViewSettings
)

type tickMsg time.Time
type lyricsMsg struct {
	data lyrics.Lyrics
	err  error
}

type Model struct {
	configMgr *config.Manager
	Config    config.Config
	Theme     config.Theme
	styles    Styles

	player  player.Controller
	fetcher lyrics.Fetcher

	viewState    ViewState
	state        player.State
	lyrics       lyrics.Lyrics
	lyricsLoaded bool
	loading      bool
	lastFetchKey string

	viewport  viewport.Model
	progress  progress.Model
	width     int
	height    int
	activeIdx int

	settingsIdx int
}

func NewModel() (Model, error) {
	mgr, err := config.NewManager()
	if err != nil {
		return Model{}, err
	}
	cfg, _ := mgr.Load()

	pb := progress.New(progress.WithDefaultGradient())
	pb.Width = 50

	t := config.GetTheme(cfg.ThemeName)

	return Model{
		configMgr: mgr,
		Config:    cfg,
		Theme:     t,
		styles:    InitStyles(t, 0),
		player:    player.NewSpotifyCtl(),
		fetcher:   lyrics.NewFetcher(),
		progress:  pb,
		viewState: ViewMain,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport = viewport.New(msg.Width, msg.Height-5)
		m.progress.Width = msg.Width - 4
		m.styles = InitStyles(m.Theme, m.width)
		return m, nil

	case tickMsg:
		newState, _ := m.player.GetState()
		m.state = newState

		if m.viewState == ViewMain {
			key := fmt.Sprintf("%s|%s", m.state.Artist, m.state.Title)
			if m.state.Status != "STOPPED" && key != m.lastFetchKey && !m.loading && m.state.Artist != "" {
				m.lastFetchKey = key
				m.loading = true
				m.lyricsLoaded = false
				cmds = append(cmds, m.fetchCmd(m.state.Artist, m.state.Title))
			}
			m.updateViewport()
		}

		cmds = append(cmds, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}))

	case lyricsMsg:
		m.loading = false
		if msg.err == nil {
			m.lyrics = msg.data
			m.lyricsLoaded = true
			m.viewport.SetYOffset(0)
		}

	case tea.KeyMsg:
		if m.viewState == ViewSettings {
			newModel, settingsCmd := m.updateSettings(msg)
			return newModel, tea.Batch(append(cmds, settingsCmd)...)
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "s":
				m.viewState = ViewSettings
				m.settingsIdx = 0
			case "o":
				m.Config.OffsetMs += 100
			case "p":
				m.Config.OffsetMs -= 100
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) fetchCmd(artist, title string) tea.Cmd {
	return func() tea.Msg {
		l, err := m.fetcher.Fetch(context.Background(), artist, title)
		return lyricsMsg{data: l, err: err}
	}
}

func (m *Model) updateViewport() {
	if !m.lyricsLoaded {
		msg := "\n\n   Waiting for playback..."
		if m.loading {
			msg = fmt.Sprintf("\n\n   Fetching: %s - %s...", m.state.Artist, m.state.Title)
		}
		m.viewport.SetContent(msg)
		return
	}

	effTime := m.state.Position + float64(m.Config.OffsetMs)/1000.0

	m.activeIdx = -1
	for i, line := range m.lyrics {
		if line.Time <= effTime {
			m.activeIdx = i
		} else {
			break
		}
	}

	if m.activeIdx >= 0 {
		half := m.viewport.Height / 2
		if m.activeIdx > half {
			m.viewport.SetYOffset(m.activeIdx - half)
		} else {
			m.viewport.SetYOffset(0)
		}
	}

	m.viewport.SetContent(m.renderLyrics(effTime))
}

func (m Model) renderLyrics(currTime float64) string {
	var sb strings.Builder

	for i, line := range m.lyrics {
		if i == m.activeIdx {
			if m.Config.KaraokeMode {
				sb.WriteString(m.renderKaraoke(line, currTime) + "\n")
			} else {
				sb.WriteString(m.styles.Active.Render(line.Text) + "\n")
			}
		} else {
			if m.Config.DimInactive {
				sb.WriteString(m.styles.Dim.Render(line.Text) + "\n")
			} else {
				sb.WriteString(m.styles.Inactive.Render(line.Text) + "\n")
			}
		}
	}
	return sb.String()
}

func (m Model) renderKaraoke(line lyrics.Line, currTime float64) string {
	prog := 0.0
	if line.Duration > 0 {
		prog = (currTime - line.Time) / line.Duration
	}
	idx := int(math.Floor(prog * float64(len(line.Words))))

	var sb strings.Builder

	for i, w := range line.Words {
		if i <= idx {
			sb.WriteString(m.styles.Sung.Render(w) + " ")
		} else {
			sb.WriteString(m.styles.Future.Render(w) + " ")
		}
	}
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(sb.String())
}

func (m Model) View() string {
	if m.viewState == ViewSettings {
		return m.settingsView()
	}

	header := m.styles.Header.Render(fmt.Sprintf(" %s - %s ", m.state.Artist, m.state.Title))

	prog := ""
	if m.Config.ShowProgress && m.state.Duration > 0 {
		pct := m.state.Position / m.state.Duration
		if pct > 1.0 {
			pct = 1.0
		}
		prog = "\n" + m.progress.ViewAs(pct)
	}

	status := m.styles.Status.Render(fmt.Sprintf("Offset: %dms | Settings (s) | Quit (q)", m.Config.OffsetMs))

	return fmt.Sprintf("%s\n%s\n%s\n%s", header, m.viewport.View(), prog, status)
}
