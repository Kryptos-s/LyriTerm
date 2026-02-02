LyriTerm
=======

LyriTerm is a terminal app that shows synced lyrics for the song you’re playing.
It reads playback state via MPRIS (through playerctl) and fetches lyrics from LRCLIB.

Features
--------
- synced, auto-scrolling lyrics
- karaoke-style word highlighting when the lyric data supports it
- themes: spotify, nord, dracula, gruvbox, monokai, arch
- runtime offset control (+/- ms) for syncing
- built with Bubble Tea

Requirements
------------
- playerctl (required)
- a player that exposes MPRIS (example: Spotify, VLC)

Install playerctl (examples)
----------------------------
Arch Linux:
  pacman -S playerctl

Debian/Ubuntu:
  sudo apt install playerctl

Fedora:
  sudo dnf install playerctl


Install
-------

Arch (AUR)
~~~~~~~~~~
If you published lyriterm-git:
  yay -S lyriterm-git


Go install
~~~~~~~~~~
Requires Go.

  go install github.com/Kryptos-s/LyriTerm/cmd/lyriterm@latest

Make sure $(go env GOPATH)/bin is in your PATH.


Build from source
~~~~~~~~~~~~~~~~~
  git clone https://github.com/Kryptos-s/LyriTerm.git
  cd LyriTerm
  go build -trimpath -o lyriterm ./cmd/lyriterm
  sudo install -m755 lyriterm /usr/local/bin/lyriterm


Usage
-----
Start playback in your player, then run:
  lyriterm


Keybinds
--------
- q / Ctrl+C: quit
- s: settings
- o: offset +100ms
- p: offset -100ms
- Esc: close settings


Configuration
-------------
On first run, LyriTerm creates:
  ~/.config/lyriterm/config.json

Example config:
{
  "theme": "spotify",
  "offset_ms": 0,
  "show_progress": true,
  "karaoke_mode": true,
  "dim_inactive": true
}

Options
-------
- theme: spotify | nord | dracula | gruvbox | monokai | arch
- offset_ms: positive delays lyrics, negative advances them
- karaoke_mode: word highlighting (only if the lyric source provides word timing)
- dim_inactive: dims non-current lines
- show_progress: progress bar


Troubleshooting
---------------
- “Lyrics not found”: the track may not exist on LRCLIB, or it has no timed lyrics.
- No player detected: ensure `playerctl status` works and your player exposes MPRIS.
- Lyrics drift: adjust offset with o / p, then save it in config.


License
-------
MIT
