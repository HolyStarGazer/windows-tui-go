package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HolyStarGazer/windows-tui-go/types"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	CurrentPath string
	Items       []types.FileItem
	Cursor      int
	Width       int
	Height      int
	Err         error
	Mode        ViewMode
	FileViewer  *FileViewer
}

// NewModel creates and returns the initial model state
func NewModel() Model {
	// Start in the current directory
	currentPath, err := os.Getwd()
	if err != nil {
		currentPath = "."
	}

	m := Model{
		CurrentPath: currentPath,
		Cursor:      0,
		Mode:        BrowseMode,
	}
	m.loadDirectory()
	return m
}

// loadDirectory reads teh contents of the current directory
func (m *Model) loadDirectory() {
	m.Items = []types.FileItem{}
	m.Cursor = 0
	m.Err = nil

	// Add parent directory entry if not at root
	if m.CurrentPath != filepath.VolumeName(m.CurrentPath)+string(filepath.Separator) {
		m.Items = append(m.Items, types.FileItem{
			Name:  "..",
			Path:  filepath.Dir(m.CurrentPath),
			IsDir: true,
		})
	}

	entries, err := os.ReadDir(m.CurrentPath)
	if err != nil {
		m.Err = err
		return
	}

	// Separate directories and files
	var dirs []types.FileItem
	var files []types.FileItem

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := types.FileItem{
			Name:  entry.Name(),
			Path:  filepath.Join(m.CurrentPath, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		}

		if entry.IsDir() {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}

	// Add directories first, then files
	m.Items = append(m.Items, dirs...)
	m.Items = append(m.Items, files...)
}

// Init initializes the model (called once at startup)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
		if m.FileViewer != nil {
			m.FileViewer.Height = msg.Height
			m.FileViewer.Width = msg.Width
		}
		return m, nil

	case tea.KeyMsg:
		// Handle file viewer mode
		if m.Mode == FileViewMode {
			switch msg.String() {
			case "q", "esc":
				// Return to browse mode
				m.Mode = BrowseMode
				m.FileViewer = nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				// Pass other keys to the file viewer
				if m.FileViewer != nil {
					m.FileViewer.Update(msg)
				}
			}
			return m, nil
		}

		// Handle browse mode
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}

		case "enter", "l", "right":
			if len(m.Items) > 0 {
				selected := m.Items[m.Cursor]
				if selected.IsDir {
					m.CurrentPath = selected.Path
					m.loadDirectory()
				} else {
					// Open file viewer
					viewer := NewFileViewer(selected.Path, selected.Name)
					viewer.Height = m.Height
					viewer.Width = m.Width
					m.FileViewer = &viewer
					m.Mode = FileViewMode
				}
			}

		case "h", "left", "backspace":
			// Go to parent directory
			parent := filepath.Dir(m.CurrentPath)
			if parent != m.CurrentPath {
				m.CurrentPath = parent
				m.loadDirectory()
			}

		case "g":
			// Go to top
			m.Cursor = 0

		case "G":
			// Go to bottom
			if len(m.Items) > 0 {
				m.Cursor = len(m.Items) - 1
			}
		}
	}

	return m, nil
}

// View renders the current state of the model
func (m Model) View() string {
	// If in file viewer mode, show the file viewer
	if m.Mode == FileViewMode && m.FileViewer != nil {
		return m.FileViewer.View()
	}

	// Otherwise show the file browser
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.Err)
	}

	var b strings.Builder

	// Title
	title := titleStyle.Render("📁 File Explorer")
	b.WriteString(title + "\n")

	// Current Path
	pathDisplay := fmt.Sprintf("Current Path: %s", m.CurrentPath)
	b.WriteString(pathDisplay + "\n\n")

	// File list
	visibleStart := 0
	visibleEnd := len(m.Items)
	maxVisible := m.Height - 8 // Reserve space for header and footer

	if maxVisible > 0 && len(m.Items) > maxVisible {
		// Calculate visible windows
		if m.Cursor >= maxVisible/2 {
			visibleStart = m.Cursor - maxVisible/2
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(m.Items) {
			visibleEnd = len(m.Items)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		item := m.Items[i]
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		// Format the item
		var itemStr string
		if item.IsDir {
			itemStr = directoryStyle.Render("📁 " + item.Name + "/")
		} else {
			sizeStr := FormatSize(item.Size)
			itemStr = fileStyle.Render(fmt.Sprintf("📄 %s (%s)", item.Name, sizeStr))
		}

		// Apply selection style if this is the cursor position
		line := fmt.Sprintf("%s %s", cursor, itemStr)
		if m.Cursor == i {
			line = selectedStyle.Render(line)
		}

		b.WriteString(line + "\n")
	}

	// Status bar
	if len(m.Items) > 0 {
		status := statusStyle.Render(fmt.Sprintf("\n%d/%d items", m.Cursor+1, len(m.Items)))
		b.WriteString(status + "\n")
	}

	// Help text
	help := helpStyle.Render("↑/k: Up  ↓/j: Down  Enter/l: Open  h/Backspace: Back | g: Top | G: Bottom | q: Quit")
	b.WriteString(help)

	return b.String()
}
