package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ViewMode represents the current mode of the application
type ViewMode int

const (
	BrowseMode ViewMode = iota
	FileViewMode
)

// FileViewer handles file content viewing
type FileViewer struct {
	FilePath  string
	FileName  string
	Content   []string // Lines of the file
	ScrollPos int      // Current scroll position
	Width     int
	Height    int
	Err       error
}

// NewFileViewer creates a new file viewer for the given file path
func NewFileViewer(filePath, fileName string) FileViewer {
	fv := FileViewer{
		FilePath:  filePath,
		FileName:  fileName,
		ScrollPos: 0,
	}
	fv.loadFile()
	return fv
}

// loadFile reads the file content into memory
func (fv *FileViewer) loadFile() {
	// Read file with size limit to prevent loading huge files
	const maxFileSize = 10 * 1024 * 1024 // 10 MB limit

	fileInfo, err := os.Stat(fv.FilePath)
	if err != nil {
		fv.Err = err
		return
	}

	if fileInfo.Size() > maxFileSize {
		fv.Err = fmt.Errorf("file too large (max 10MB)")
		return
	}

	data, err := os.ReadFile(fv.FilePath)
	if err != nil {
		fv.Err = err
		return
	}

	// Split into lines - handle both Windows (\r\n) and Unix (\n) line endings
	content := string(data)
	// Normalize line endings to \n
	content = strings.ReplaceAll(content, "\r\n", "\n")
	// Remove any remaining \r (carriage return) characters
	content = strings.ReplaceAll(content, "\r", "")
	fv.Content = strings.Split(content, "\n")
}

// Update handles keyboard input for the file viewer
func (fv *FileViewer) Update(msg tea.KeyMsg) {
	maxVisible := fv.Height - 6 // Reserve space for header and footer

	switch msg.String() {
	case "up", "k":
		if fv.ScrollPos > 0 {
			fv.ScrollPos--
		}

	case "down", "j":
		maxScroll := len(fv.Content) - maxVisible
		if maxScroll < 0 {
			maxScroll = 0
		}
		if fv.ScrollPos < maxScroll {
			fv.ScrollPos++
		}

	case "g":
		// Jump to top
		fv.ScrollPos = 0

	case "G":
		// Jump to bottom
		maxScroll := len(fv.Content) - maxVisible
		if maxScroll < 0 {
			maxScroll = 0
		}
		fv.ScrollPos = maxScroll

	case "pageup", "ctrl+u":
		// Scroll up half a page
		fv.ScrollPos -= maxVisible / 2
		if fv.ScrollPos < 0 {
			fv.ScrollPos = 0
		}

	case "pagedown", "ctrl+d":
		// Scroll down half a page
		maxScroll := len(fv.Content) - maxVisible
		if maxScroll < 0 {
			maxScroll = 0
		}
		fv.ScrollPos += maxVisible / 2
		if fv.ScrollPos > maxScroll {
			fv.ScrollPos = maxScroll
		}
	}
}

// View renders the file viewer
func (fv FileViewer) View() string {
	if fv.Err != nil {
		return fmt.Sprintf("Error loading file: %v\n\nPress q or Esc to go back.", fv.Err)
	}

	var b strings.Builder

	// Title
	title := titleStyle.Render(fmt.Sprintf("📄 Viewing: %s", fv.FileName))
	b.WriteString(title + "\n")

	// File info
	info := fmt.Sprintf("Lines: %d | Position: %d", len(fv.Content), fv.ScrollPos+1)
	b.WriteString(info + "\n\n")

	// Calculate visible range
	maxVisible := fv.Height - 6
	visibleStart := fv.ScrollPos
	visibleEnd := visibleStart + maxVisible

	if visibleEnd > len(fv.Content) {
		visibleEnd = len(fv.Content)
	}

	// Display file content with line numbers
	for i := visibleStart; i < visibleEnd; i++ {
		lineNum := fmt.Sprintf("%4d │ ", i+1)
		line := fv.Content[i]

		// Convert tabs to spaces for consistent display
		line = strings.ReplaceAll(line, "\t", "    ")

		// Truncate long lines if needed
		if len(line) > fv.Width-10 {
			line = line[:fv.Width-13] + "..."
		}

		b.WriteString(lineNum + line + "\n")
	}

	// Footer with help
	b.WriteString("\n")
	help := helpStyle.Render("↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u: page up | Ctrl+d: page down | q/Esc: back")
	b.WriteString(help)

	return b.String()
}
