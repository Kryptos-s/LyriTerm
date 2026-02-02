package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const AppName = "lyriterm"

type Theme struct {
	Name    string
	Primary string
	Bg      string
	Text    string
	Dim     string
}

type Config struct {
	ThemeName    string `json:"theme"`
	OffsetMs     int    `json:"offset_ms"`
	ShowProgress bool   `json:"show_progress"`
	KaraokeMode  bool   `json:"karaoke_mode"`
	DimInactive  bool   `json:"dim_inactive"`
}

var Themes = []Theme{
	{Name: "spotify", Primary: "#1DB954", Bg: "#121212", Text: "#FFFFFF", Dim: "#535353"},
	{Name: "nord", Primary: "#88C0D0", Bg: "#2E3440", Text: "#ECEFF4", Dim: "#4C566A"},
	{Name: "dracula", Primary: "#FF79C6", Bg: "#282A36", Text: "#F8F8F2", Dim: "#6272A4"},
	{Name: "gruvbox", Primary: "#fabd2f", Bg: "#282828", Text: "#ebdbb2", Dim: "#928374"},
	{Name: "monokai", Primary: "#A6E22E", Bg: "#272822", Text: "#F8F8F2", Dim: "#75715E"},
	{Name: "arch", Primary: "#1793D1", Bg: "#0f0f0f", Text: "#eeeeee", Dim: "#4d4d4d"},
}

func GetTheme(name string) Theme {
	for _, t := range Themes {
		if t.Name == name {
			return t
		}
	}
	return Themes[0]
}

func DefaultConfig() Config {
	return Config{
		ThemeName:    "spotify",
		OffsetMs:     0,
		ShowProgress: true,
		KaraokeMode:  true,
		DimInactive:  true,
	}
}

type Manager struct {
	Path string
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Manager{
		Path: filepath.Join(home, ".config", AppName, "config.json"),
	}, nil
}

func (m *Manager) Load() (Config, error) {
	if _, err := os.Stat(m.Path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		_ = m.Save(cfg)
		return cfg, nil
	}
	data, err := os.ReadFile(m.Path)
	if err != nil {
		return DefaultConfig(), err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

func (m *Manager) Save(cfg Config) error {
	dir := filepath.Dir(m.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.Path, data, 0644)
}
