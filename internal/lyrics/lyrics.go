package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Line struct {
	Time     float64
	Text     string
	Duration float64
	Words    []string
}

type Lyrics []Line

type Fetcher interface {
	Fetch(ctx context.Context, artist, title string) (Lyrics, error)
}

type LrcLibFetcher struct {
	Client   *http.Client
	CacheDir string
}

func NewFetcher() *LrcLibFetcher {

	home, err := os.UserHomeDir()
	var cacheDir string
	if err == nil {
		cacheDir = filepath.Join(home, ".config", "lyriterm", "cache")
	}

	return &LrcLibFetcher{
		Client:   &http.Client{Timeout: 5 * time.Second},
		CacheDir: cacheDir,
	}
}

type apiResponse struct {
	SyncedLyrics string `json:"syncedLyrics"`
}

func (f *LrcLibFetcher) Fetch(ctx context.Context, artist, title string) (Lyrics, error) {

	artist = strings.TrimSpace(artist)
	title = strings.TrimSpace(title)

	if cachedLrc, ok := f.loadFromCache(artist, title); ok {
		return parseLRC(cachedLrc), nil
	}

	cleanTitle := strings.Split(title, " - ")[0]
	cleanTitle = strings.Split(cleanTitle, " (")[0]

	params := url.Values{}
	params.Add("artist_name", artist)
	params.Add("track_name", cleanTitle)

	lrc, err := f.doRequest(ctx, "https://lrclib.net/api/get", params)
	if err == nil && lrc != "" {
		f.saveToCache(artist, title, lrc)
		return parseLRC(lrc), nil
	}

	lrc, err = f.doSearch(ctx, params)
	if err == nil && lrc != "" {
		f.saveToCache(artist, title, lrc)
		return parseLRC(lrc), nil
	}

	return nil, fmt.Errorf("no lyrics found")
}

func (f *LrcLibFetcher) getCachePath(artist, title string) string {
	if f.CacheDir == "" {
		return ""
	}

	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	safeArtist := reg.ReplaceAllString(strings.ToLower(artist), "_")
	safeTitle := reg.ReplaceAllString(strings.ToLower(title), "_")

	safeArtist = strings.Trim(safeArtist, "_")
	safeTitle = strings.Trim(safeTitle, "_")

	filename := fmt.Sprintf("%s-%s.lrc", safeArtist, safeTitle)
	return filepath.Join(f.CacheDir, filename)
}

func (f *LrcLibFetcher) loadFromCache(artist, title string) (string, bool) {
	path := f.getCachePath(artist, title)
	if path == "" {
		return "", false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (f *LrcLibFetcher) saveToCache(artist, title, content string) {
	path := f.getCachePath(artist, title)
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {

		return
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {

	}
}

func (f *LrcLibFetcher) doRequest(ctx context.Context, endpoint string, params url.Values) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
	resp, err := f.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var res apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.SyncedLyrics, nil
}

func (f *LrcLibFetcher) doSearch(ctx context.Context, params url.Values) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://lrclib.net/api/search"+"?"+params.Encode(), nil)
	resp, err := f.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res []apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil || len(res) == 0 {
		return "", fmt.Errorf("not found")
	}
	return res[0].SyncedLyrics, nil
}

func parseLRC(text string) Lyrics {
	var lines Lyrics
	for _, l := range strings.Split(text, "\n") {
		if !strings.Contains(l, "]") {
			continue
		}
		parts := strings.SplitN(l, "]", 2)
		timePart := strings.TrimPrefix(parts[0], "[")
		t := strings.Split(timePart, ":")
		if len(t) != 2 {
			continue
		}
		m, _ := strconv.ParseFloat(t[0], 64)
		s, _ := strconv.ParseFloat(t[1], 64)
		lines = append(lines, Line{
			Time:  (m * 60) + s,
			Text:  strings.TrimSpace(parts[1]),
			Words: strings.Fields(strings.TrimSpace(parts[1])),
		})
	}
	sort.Slice(lines, func(i, j int) bool { return lines[i].Time < lines[j].Time })
	for i := 0; i < len(lines)-1; i++ {
		lines[i].Duration = lines[i+1].Time - lines[i].Time
	}
	if len(lines) > 0 {
		lines[len(lines)-1].Duration = 5.0
	}
	return lines
}
