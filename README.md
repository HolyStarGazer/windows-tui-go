# Windows TUI File Explorer

A lightweight, terminal-based file explorer built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Navigate your file system with keyboard shortcuts in a beautiful, responsive TUI (Terminal User Interface).

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

## Features

- 📂 Browse directories with an intuitive interface
- 📖 **Read-only file viewer** with vim-style navigation
- ⌨️ Vim-style keyboard navigation (`hjkl`) + arrow keys
- 🎨 Syntax-highlighted files and folders
- 📊 Human-readable file sizes
- 🔢 Line numbers in file viewer
- 🚀 Fast and lightweight (single executable, no dependencies)
- 🪟 Native Windows support (handles CRLF line endings)
- 💻 Works in Windows Terminal, PowerShell, and VSCode

## Examples

**File Browser:**
```
📁 File Explorer
Current: C:\Users\YourName\Documents

  📂 projects/
  📂 notes/
> 📄 report.docx (2.4 MB)
  📄 budget.xlsx (156.3 KB)
  📄 readme.md (4.2 KB)

4/5 items
↑/k: Up  ↓/j: Down  Enter/l: Open  h/Backspace: Back | g: Top | G: Bottom | q: Quit
```

**File Viewer:**
```
📄 Viewing: main.go
Lines: 142 | Position: 1

   1 │ package main
   2 │ 
   3 │ import (
   4 │     "fmt"
   5 │     "os"
   6 │ 
   7 │     "github.com/HolyStarGazer/windows-tui-go/ui"
   8 │     tea "github.com/charmbracelet/bubbletea"
   9 │ )
  10 │ 

↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u: page up | Ctrl+d: page down | q/Esc: back
```

## Installation

### Prerequisites

- Go 1.21 or higher ([Download Go](https://go.dev/download/))
- Windows Terminal or PowerShell (recommended for best experience)

### Quick Start

1. **Clone or download this repository**
   ```bash
   git clone https://github.com/yourusername/windows-tui-go.git
   cd windows-tui-go
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run .
   ```

### Build Executable

Create a standalone `.exe` file you can run anywhere:

```bash
go build -o file-explorer.exe
```

Now you can run it directly:
```bash
.\file-explorer.exe
```

Or move it to a directory in your PATH to run it from anywhere!

## Usage

### Keyboard Shortcuts

#### File Browser Mode
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `l` / `→` | Open directory or view file |
| `h` / `←` / `Backspace` | Go to parent directory |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `q` / `Ctrl+C` | Quit |

#### File Viewer Mode
| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up one line |
| `↓` / `j` | Scroll down one line |
| `g` | Jump to top of file |
| `G` | Jump to bottom of file |
| `Ctrl+u` | Page up (half screen) |
| `Ctrl+d` | Page down (half screen) |
| `q` / `Esc` | Return to file browser |
| `Ctrl+C` | Quit application |

### Tips

- **Use Windows Terminal** for the best experience with emoji support and better colors
- **VSCode Integration**: Works perfectly in VSCode's integrated terminal
- **Portable**: Copy `file-explorer.exe` to a USB drive and run it on any Windows machine
- **File Viewing**: Press Enter on any text file to read its contents - works great for `.go`, `.md`, `.txt`, `.json`, `.xml`, and other text files
- **Large Files**: Files over 10MB cannot be viewed to prevent performance issues

## Project Structure

```
windows-tui-go/
├── main.go              # Application entry point
├── ui/
│   ├── model.go         # TUI state management and logic
│   ├── viewer.go        # File viewer component
│   ├── styles.go        # Lipgloss styling definitions
│   └── utils.go         # Utility functions
├── types/
│   └── types.go         # Data structures
├── go.mod               # Go module definition
└── README.md            # This file
```

## How It Works

This project uses the **Elm Architecture** pattern via [Bubble Tea](https://github.com/charmbracelet/bubbletea):

1. **Model** - Application state (current directory, selected item, etc.)
2. **Update** - Handles events (keyboard input) and updates the model
3. **View** - Renders the UI based on the current model state

### Key Technologies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Style and layout library
- **Go Standard Library** - File system operations

## Development

### Adding Features

The modular structure makes it easy to add features:

- Add new UI components in `ui/`
- Add data structures in `types/`
- Add new key bindings in `ui/model.go` → `Update()` method
- Customize colors in `ui/styles.go`

### Planned Features

- [x] File viewer (read-only)
- [ ] File search/filter
- [ ] File operations (copy, delete, rename)
- [ ] File preview pane
- [ ] Bookmarks for quick navigation
- [ ] Dual-pane mode
- [ ] Hidden files toggle
- [ ] Sort options (name, size, date)
- [ ] Syntax highlighting in viewer

## Building for Distribution

### Single Executable

```bash
go build -ldflags="-s -w" -o file-explorer.exe
```

The `-ldflags="-s -w"` flag reduces the executable size by stripping debug information.

### Cross-Compilation

Build for other platforms (if desired):

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o file-explorer

# For macOS
GOOS=darwin GOARCH=amd64 go build -o file-explorer
```

## Troubleshooting

### "go: command not found"
- Make sure Go is installed: [Download Go](https://go.dev/download/)
- Restart VSCode/terminal after installing Go

### Emojis not showing up
- Use Windows Terminal instead of old cmd.exe
- Install Windows Terminal from Microsoft Store

### Colors look wrong
- Windows Terminal or modern PowerShell recommended
- Check that your terminal supports 256 colors

## Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests

## License

MIT License - feel free to use this project however you'd like!

## Acknowledgments

- [Charm](https://charm.sh/) for the excellent TUI libraries
- The Go community for amazing tools and resources

## Learn More

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Go by Example](https://gobyexample.com/)

---

**Happy exploring! 🚀**