# LyriTerm

Synced lyrics in your terminal for the currently playing song (MPRIS via `playerctl`, lyrics from LRCLIB).

<video src="preview.mp4" autoplay loop muted playsinline></video>


## Features
- synced auto-scrolling lyrics
- karaoke-style highlighting when timing is available
- themes: spotify, nord, dracula, gruvbox, monokai, arch
- offset control (ms) to fix drift

## Requirements
- `playerctl`
- an MPRIS-capable player (Spotify, VLC, etc.)

### Install playerctl
**Arch**
```bash
sudo pacman -S playerctl
```

**Debian/Ubuntu**
```bash
sudo apt install playerctl
```

**Fedora**
```bash
sudo dnf install playerctl
```

## Install

### Arch (AUR)
```bash
yay -S lyriterm-git
```

### Go
```bash
go install github.com/Kryptos-s/LyriTerm/cmd/lyriterm@latest
```

Make sure Go’s bin dir is in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Build from source
```bash
git clone https://github.com/Kryptos-s/LyriTerm.git
cd LyriTerm
go build -trimpath -o lyriterm ./cmd/lyriterm
sudo install -m755 lyriterm /usr/local/bin/lyriterm
```

## Usage
Start music in your player, then run:
```bash
lyriterm
```

## Keybinds

| Key | Action |
|---|---|
| `q` / `Ctrl+C` | quit |
| `s` | open settings |
| `o` | offset +100ms |
| `p` | offset -100ms |
| `Esc` | close settings |

## Config
Config file:
- `~/.config/lyriterm/config.json`

Example:
```json
{
  "theme": "spotify",
  "offset_ms": 0,
  "show_progress": true,
  "karaoke_mode": true,
  "dim_inactive": true
}
```

Options:

| Field | Type | Default | Notes |
|---|---|---:|---|
| `theme` | string | `spotify` | `spotify,nord,dracula,gruvbox,monokai,arch` |
| `offset_ms` | number | `0` | + delays lyrics, - advances |
| `show_progress` | bool | `true` | progress bar |
| `karaoke_mode` | bool | `true` | needs word timing data |
| `dim_inactive` | bool | `true` | dims non-current lines |

## Troubleshooting
- **No player detected**: run `playerctl status`. if it fails, your player is not exposing MPRIS.
- **Lyrics not found**: track may not exist on LRCLIB or has no timed lyrics.
- **Drift**: adjust offset with `o` / `p`, then set `offset_ms` in config.

## License
MIT
