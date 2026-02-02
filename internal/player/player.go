package player

import (
	"os/exec"
	"strconv"
	"strings"
)

type State struct {
	Artist   string
	Title    string
	Status   string
	Position float64
	Duration float64
}

type Controller interface {
	GetState() (State, error)
}

type SpotifyCtl struct{}

func NewSpotifyCtl() *SpotifyCtl {
	return &SpotifyCtl{}
}

func (s *SpotifyCtl) GetState() (State, error) {
	// 1. Check Status
	out, err := exec.Command("playerctl", "-p", "spotify", "status").Output()
	if err != nil {
		return State{Status: "STOPPED"}, nil
	}
	status := strings.TrimSpace(string(out))

	metaOut, err := exec.Command("playerctl", "-p", "spotify", "metadata", "--format", "{{artist}}:::{{title}}:::{{mpris:length}}").Output()
	if err != nil {
		return State{Status: status}, nil
	}
	parts := strings.Split(strings.TrimSpace(string(metaOut)), ":::")
	if len(parts) < 3 {
		return State{Status: status}, nil
	}

	artist := parts[0]
	title := parts[1]

	durMicro, _ := strconv.ParseFloat(parts[2], 64)
	duration := durMicro / 1000000.0

	posOut, err := exec.Command("playerctl", "-p", "spotify", "position").Output()
	pos := 0.0
	if err == nil {
		pos, _ = strconv.ParseFloat(strings.TrimSpace(string(posOut)), 64)
	}

	return State{
		Artist:   artist,
		Title:    title,
		Status:   status,
		Position: pos,
		Duration: duration,
	}, nil
}
