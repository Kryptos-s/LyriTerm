


# LyriTerm

LyriTerm is a terminal-based utility written in Go that displays synchronized lyrics for the currently playing song. It interfaces with MPRIS-compatible media players (Spotify, VLC, etc.) via `playerctl` and fetches lyrics from LRCLIB.

## Features

- **synchronized Lyrics**: Displays auto-scrolling lyrics synchronized to playback.
- **Karaoke Highlighting**: Highlights individual words in real-time when supported by the lyric data.
- **Multiple Themes**: Includes color schemes for Spotify, Nord, Dracula, Gruvbox, Monokai, and Arch Linux.
- **Offset Adjustment**: Allows runtime adjustment of lyric timing to correct synchronization issues.
- **Resource Efficient**: Built with Bubble Tea for minimal system resource usage.

## Requirements

LyriTerm requires `playerctl` to communicate with media players.

- **Arch Linux**: `pacman -S playerctl`
- **Debian/Ubuntu**: `apt install playerctl`
- **Fedora**: `dnf install playerctl`

## Installation

### Using Go Install

If you have a Go environment set up, you can install the binary directly:

```bash
go install https://github.com/kryptos-s/lyriterm/

```

Ensure your `$GOPATH/bin` is added to your system `$PATH`.

### Building from Source

To build the binary manually:

```bash
git clone https://github.com/kryptos-s/lyriterm.git
cd lyriterm
go build -o lyriterm lyriterm/main.go

```

You can then move the `lyriterm` binary to a directory in your path, such as `/usr/local/bin`.

## Usage

Start the application while a supported media player is active:

```bash
lyriterm

```

### Key Bindings

| Key | Action |
| --- | --- |
| `q` / `Ctrl+C` | Quit application |
| `s` | Open Settings menu |
| `o` | Increase global offset (+100ms) |
| `p` | Decrease global offset (-100ms) |
| `Esc` | Close Settings menu |

## Configuration

On the first run, a configuration file is generated at `~/.config/lyriterm/config.json`.

**Default Configuration:**

```json
{
  "theme": "spotify",
  "offset_ms": 0,
  "show_progress": true,
  "karaoke_mode": true,
  "dim_inactive": true
}

```

### Options

* **theme**: The color scheme to use. Options: `spotify`, `nord`, `dracula`, `gruvbox`, `monokai`, `arch`.
* **offset_ms**: Global sync offset in milliseconds. Positive values delay the lyrics; negative values speed them up.
* **karaoke_mode**: Enables word-by-word highlighting (requires synced word data from the API).
* **dim_inactive**: Lowers the opacity of lines that are not currently being sung.
* **show_progress**: Toggles the playback progress bar at the bottom of the window.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

```

```