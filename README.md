# Windows TUI File Explorer

A lightweight, terminal-based file explorer built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Navigate your file system with keyboard shortcuts in a beautiful, responsive TUI (Terminal User Interface).

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=flat&logo=windows)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

## Features

- 📂 Browse directories with an intuitive interface
- 📖 **Read-only file viewer** with vim-style navigation
- ⌨️ **Vim-style command mode** (`:` to enter commands)
- 🔍 **Full-text search** with highlighted matches and navigation
- 🎨 **Syntax highlighting** for 200+ languages (Go, Python, JS, Java, C/C++, Rust, and more)
- ⚡ Vim-style keyboard navigation (`hjkl`) + arrow keys
- 🌈 Color-coded files and folders in browser
- 📊 Human-readable file sizes
- 🔢 Line numbers in file viewer
- 🔄 Optional line wrapping (toggle via command)
- 🚀 Fast and lightweight (single executable, no dependencies)
- 🪟 Native Windows support (handles CRLF line endings)
- 💻 Works in Windows Terminal, PowerShell, and VSCode

## Screenshots

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

**File Viewer (with Syntax Highlighting and Search):**
```
📄 Viewing: main.go
Lines: 17 | Position: 1 | Wrap: OFF

   1 │ package main
   2 │ 
   3 │ import (
   4 │     "fmt"
   5 │     "os"
   6 │ 
   7 │     tea "github.com/charmbracelet/bubbletea"
   8 │     "github.com/HolyStarGazer/windows-tui-go/ui"
   9 │ )
  10 │ 
  11 │ func main() {
  12 │     p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
  13 │     if _, err := p.Run(); err != nil {
  14 │         fmt.Printf("Error: %v\n", err)
  15 │         os.Exit(1)
  16 │     }
  17 │ }

Found 3 match(es) - n: next, N: prev
↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u: page up | Ctrl+d: page down | q/Esc: back
```
*Note: Keywords appear in color, and search terms are highlighted with yellow background*

**Command Mode:**
```
:search func_
```
*Type commands after pressing `:` - try :help for available commands*

## Installation

### Prerequisites

- Go 1.21 or higher ([Download Go](https://go.dev/dl/))
- Windows Terminal or PowerShell (recommended for best experience)

### Quick Start

1. **Clone or download this repository**
   ```bash
   git clone https://github.com/HolyStarGazer/windows-tui-go.git
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
| `n` | Next search match |
| `N` | Previous search match |
| `:` | Enter command mode |
| `q` / `Esc` | Return to file browser |
| `Ctrl+C` | Quit application |

#### Command Mode (press `:` to enter)
| Command | Action |
|---------|--------|
| `:set wrap` | Enable line wrapping |
| `:set nowrap` | Disable line wrapping |
| `:set syntax` | Enable syntax highlighting |
| `:set nosyntax` | Disable syntax highlighting |
| `:wrap` | Toggle line wrapping |
| `:syntax` | Toggle syntax highlighting |
| `:search <term>` | Search for text |
| `:/<pattern>` | Quick search (vim-style) |
| `:n` or `:next` | Jump to next match |
| `:N` or `:prev` | Jump to previous match |
| `:clear` | Clear search highlighting |
| `:help` or `:h` | Show available commands |
| `Esc` | Cancel command |

### Tips

- **Use Windows Terminal** for the best experience with emoji support and syntax highlighting colors
- **VSCode Integration**: Works perfectly in VSCode's integrated terminal
- **Portable**: Copy `file-explorer.exe` to a USB drive and run it on any Windows machine
- **File Viewing**: Press Enter on any text file to read its contents with automatic syntax highlighting - works great for `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.md`, `.json`, `.xml`, and 200+ more file types
- **Large Files**: Files over 10MB cannot be viewed to prevent performance issues
- **Syntax Colors**: The viewer uses the Monokai theme - keywords, strings, comments, and more are automatically colorized
- **Command Mode**: Press `:` to access all viewer options - try `:help` to see available commands
- **Quick Search**: Use `:/pattern` to quickly search for text, then `n` and `N` to navigate through matches
- **Line Wrapping**: Toggle with `:wrap` - useful for long lines of code
- **Persistent Search**: Search highlighting stays active as you scroll - use `:clear` to remove

### Command Mode Examples


**Searching for text:**
```
Press :
Type /main
→ All instances of "main" highlighted in yellow
→ Press n to jump to next match
→ Press N to jump to previous match
```

**Toggling options:**
```
:wrap          → Toggle line wrapping on/off
:syntax        → Toggle syntax highlighting on/off
:set nowrap    → Explicitly disable wrapping
:set syntax    → Explicitly enable syntax highlighting
```


**Getting help:**
```
:help          → Show all available commands
```

**Clearing search:**
```
:clear         → Remove search highlighting
```

**Example session:**
1. Open a Go file: `file-explorer.exe`
2. Navigate to a `.go` file and press Enter
3. Search for functions: `:/func`
4. Navigate results: `n`, `n`, `N`
5. Toggle wrapping: `:wrap`
6. Check available commands: `:help`
7. Return to browser: `Esc` or `q`

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
- **[Chroma](https://github.com/alecthomas/chroma)** - Syntax highlighting for 200+ languages
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
- [x] Syntax highlighting in viewer
- [x] Vim-style command mode
- [x] Full-text search with highlighting
- [ ] Jump to line number (`:goto <line>` or `:<number>`)
- [ ] File operations (copy, delete, rename)
- [ ] File preview pane
- [ ] Bookmarks for quick navigation
- [ ] Dual-pane mode
- [ ] Hidden files toggle
- [ ] Sort options (name, size, date)
- [ ] Custom color themes (`:theme <name>`)
- [ ] Search history
- [ ] Regular expression search

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