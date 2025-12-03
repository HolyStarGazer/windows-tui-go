package ui

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
	FilePath           string
	FileName           string
	Content            []string // Lines of the file
	HighlightedContent []string // Lines with syntax highlighting
	ScrollPos          int      // Current scroll position
	Width              int
	Height             int
	Err                error
	UseSyntaxHighlight bool
}

// NewFileViewer creates a new file viewer for the given file path
func NewFileViewer(filePath, fileName string) FileViewer {
	fv := FileViewer{
		FilePath:           filePath,
		FileName:           fileName,
		ScrollPos:          0,
		UseSyntaxHighlight: true,
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
	// Convert tabs to spaces BEFORE highlighting for consistent display
	content = strings.ReplaceAll(content, "\t", "    ")
	fv.Content = strings.Split(content, "\n")

	// Optionally apply syntax highlighting
	if fv.UseSyntaxHighlight {
		fv.applySyntaxHighlighting(content)
	}
}

// applySyntaxHighlighting applies syntax highlighting to the file content
func (fv *FileViewer) applySyntaxHighlighting(content string) {
	// Get lexer based on file extension
	lexer := lexers.Match(fv.FileName)
	if lexer == nil {
		// Fallback to analzing content
		lexer = lexers.Analyse(content)
	}
	if lexer == nil {
		// If still no lexer found, use plaintext
		lexer = lexers.Fallback
	}

	// Use a terminal-friendly style
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	// Create a terminal formatter with 16 colors for better compatibility
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Tokenize and format
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		// If highlighting fails, just use plain content
		fv.HighlightedContent = fv.Content
		return
	}

	// Format to ANSI colors
	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		// If formatting fails, just use plain content
		fv.HighlightedContent = fv.Content
		return
	}

	// Split highlighted content into lines
	highlightedContent := buf.String()
	// Normalize line endings to match how we handled the plain content
	highlightedContent = strings.ReplaceAll(highlightedContent, "\r\n", "\n")
	highlightedContent = strings.ReplaceAll(highlightedContent, "\r", "")
	fv.HighlightedContent = strings.Split(highlightedContent, "\n")
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
	// Use highlighted content if available, otherwise use plain content
	contentToDisplay := fv.Content
	if len(fv.HighlightedContent) > 0 && fv.UseSyntaxHighlight {
		contentToDisplay = fv.HighlightedContent
	}

	for i := visibleStart; i < visibleEnd; i++ {
		lineNum := fmt.Sprintf("%4d │ ", i+1)
		var line string

		// Make sure we don't go out of bounds
		if i < len(contentToDisplay) {
			line = contentToDisplay[i]
		}

		// Truncate long lines if needed
		// if fv.Width > 15 && len(line) > fv.Width-10 {
		// 	line = line[:fv.Width-13] + "..."
		// }

		b.WriteString(lineNum + line + "\n")
	}

	// Footer with help
	b.WriteString("\n")
	help := helpStyle.Render("↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u: page up | Ctrl+d: page down | q/Esc: back")
	b.WriteString(help)

	return b.String()
}
