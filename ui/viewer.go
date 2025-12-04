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
	UseSyntaxHighlight bool   // Toggle for syntax highlighting
	WrapLines          bool   // Toggle for line wrapping
	CommandMode        bool   // Whether in command mode
	CommandBuffer      string // Buffer for command input
	StatusMessage      string // Status or error messages
	SearchTerm         string // Current search term
	SearchMatches      []int  // Line numbers with matches
	CurrentMatchIndex  int    // Index of the current match
}

// NewFileViewer creates a new file viewer for the given file path
func NewFileViewer(filePath, fileName string) FileViewer {
	fv := FileViewer{
		FilePath:           filePath,
		FileName:           fileName,
		ScrollPos:          0,
		UseSyntaxHighlight: true,
		WrapLines:          false,
		CommandMode:        false,
		CommandBuffer:      "",
		StatusMessage:      "",
		SearchTerm:         "",
		SearchMatches:      []int{},
		CurrentMatchIndex:  -1,
	}
	fv.loadFile()
	return fv
}

// executeCommand parses and executes a command
func (fv *FileViewer) executeCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return
	}

	// Handle vim-style /search shortcut
	if strings.HasPrefix(cmd, "/") {
		searchTerm := strings.TrimPrefix(cmd, "/")
		fv.performSearch(searchTerm)
		return
	}

	// Split command into parts
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "/", "search":
		// Search command
		if len(parts) < 2 {
			fv.StatusMessage = "Usage: :search <term> or :/<term>"
			return
		}
		searchTerm := strings.Join(parts[1:], " ")
		fv.performSearch(searchTerm)

	case "set":
		// Set options
		if len(parts) < 2 {
			fv.StatusMessage = "Error: :set requires an argument"
			return
		}
		option := parts[1]

		switch option {
		case "wrap":
			fv.WrapLines = true
			fv.StatusMessage = "Line wrapping enabled"
		case "nowrap":
			fv.WrapLines = false
			fv.StatusMessage = "Line wrapping disabled"
		case "syntax":
			fv.UseSyntaxHighlight = true
			fv.StatusMessage = "Syntax highlighting enabled"
		case "nosyntax":
			fv.UseSyntaxHighlight = false
			fv.StatusMessage = "Syntax highlighting disabled"
		default:
			fv.StatusMessage = fmt.Sprintf("Unknown option '%s'", option)
		}

	case "wrap":
		fv.WrapLines = !fv.WrapLines
		if fv.WrapLines {
			fv.StatusMessage = "Line wrapping enabled"
		} else {
			fv.StatusMessage = "Line wrapping disabled"
		}

	case "syntax":
		fv.UseSyntaxHighlight = !fv.UseSyntaxHighlight
		if fv.UseSyntaxHighlight {
			fv.StatusMessage = "Syntax highlighting enabled"
		} else {
			fv.StatusMessage = "Syntax highlighting disabled"
		}

	case "help", "h":
		fv.StatusMessage = "Commands: :set [wrap|nowrap] | :set [syntax|nosyntax] | :help"

	case "n", "next":
		fv.nextMatch()

	case "N", "prev", "previous":
		fv.prevMatch()

	case "clear", "clearsearch":
		fv.performSearch("")

	default:
		fv.StatusMessage = fmt.Sprintf("Unknown command '%s' (try :help)", command)
	}
}

// performSearch searches for a term in the file content
func (fv *FileViewer) performSearch(term string) {
	if term == "" {
		fv.SearchTerm = ""
		fv.SearchMatches = []int{}
		fv.CurrentMatchIndex = -1
		fv.StatusMessage = "Search cleared"
		return
	}

	fv.SearchTerm = strings.ToLower(term)
	fv.SearchMatches = []int{}

	// Search through content (case-insensitive)
	for i, line := range fv.Content {
		if strings.Contains(strings.ToLower(line), fv.SearchTerm) {
			fv.SearchMatches = append(fv.SearchMatches, i)
		}
	}

	if len(fv.SearchMatches) > 0 {
		fv.CurrentMatchIndex = 0
		fv.ScrollPos = fv.SearchMatches[0]
		fv.StatusMessage = fmt.Sprintf("Found %d match(es) - n: next, N: prev", len(fv.SearchMatches))
	} else {
		fv.CurrentMatchIndex = -1
		fv.StatusMessage = fmt.Sprintf("Pattern not found: %s", term)
	}
}

// nextMatch jumps to the next search match
func (fv *FileViewer) nextMatch() {
	if len(fv.SearchMatches) == 0 {
		fv.StatusMessage = "No active search"
		return
	}

	fv.CurrentMatchIndex = (fv.CurrentMatchIndex + 1) % len(fv.SearchMatches)
	fv.ScrollPos = fv.SearchMatches[fv.CurrentMatchIndex]
	fv.StatusMessage = fmt.Sprintf("Match %d of %d", fv.CurrentMatchIndex+1, len(fv.SearchMatches))
}

// prevMatch jumps to the previous search match
func (fv *FileViewer) prevMatch() {
	if len(fv.SearchMatches) == 0 {
		fv.StatusMessage = "No active search"
		return
	}

	fv.CurrentMatchIndex--
	if fv.CurrentMatchIndex < 0 {
		fv.CurrentMatchIndex = len(fv.SearchMatches) - 1
	}
	fv.ScrollPos = fv.SearchMatches[fv.CurrentMatchIndex]
	fv.StatusMessage = fmt.Sprintf("Match %d of %d", fv.CurrentMatchIndex+1, len(fv.SearchMatches))
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
	// Handle command mode
	if fv.CommandMode {
		switch msg.String() {
		case "enter":
			// Execute command
			fv.executeCommand(fv.CommandBuffer)
			fv.CommandMode = false
			fv.CommandBuffer = ""

		case "esc", "ctrl+c":
			// Cancel command
			fv.CommandMode = false
			fv.CommandBuffer = ""
			fv.StatusMessage = ""

		case "backspace":
			// Delete last character
			if len(fv.CommandBuffer) > 0 {
				fv.CommandBuffer = fv.CommandBuffer[:len(fv.CommandBuffer)-1]
			}

		default:
			// Add character to command buffer (only printable characters)
			if len(msg.String()) == 1 {
				fv.CommandBuffer += msg.String()
			}
		}

		return
	}

	// Normal navigation mode
	maxVisible := fv.Height - 6 // Reserve space for header and footer

	switch msg.String() {
	case ":":
		// Enter command mode
		fv.CommandMode = true
		fv.CommandBuffer = ""
		fv.StatusMessage = ""

	case "n":
		// Next search match
		fv.nextMatch()

	case "N":
		// Previous search match
		fv.prevMatch()

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

// wrapLine wraps a line to fit within the given width, preserving ANSI color codes
func wrapLine(line string, width int, lineNum int) []string {
	if width <= 0 {
		return []string{line}
	}

	// Calculate available width (accounting for line number column: "1234 | ")
	lineNumWidth := 8 // "    1 | " = 8 characters
	availableWidth := width - lineNumWidth

	// If line is short enough, return as is
	visualLen := visualLength(line)
	if visualLen <= availableWidth {
		return []string{line}
	}

	// For long lines, we need to wrap them
	var wrapped []string
	remaining := line

	for len(remaining) > 0 && visualLength(remaining) > availableWidth {
		// Find a good breaking point
		breakPoint := findBreakPoint(remaining, availableWidth)
		if breakPoint <= 0 {
			breakPoint = availableWidth
		}

		// Split at the break point
		part := remaining[:breakPoint]
		remaining = remaining[breakPoint:]

		wrapped = append(wrapped, part)
	}

	// Add the remaining part
	if len(remaining) > 0 {
		wrapped = append(wrapped, remaining)
	}

	return wrapped
}

// visualLength calculates the visible length of a string, ignoring ANSI escape codes
func visualLength(s string) int {
	length := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' { // ESC character starts ANSI sequence
			inEscape = true
		} else if inEscape {
			if s[i] == 'm' { // 'm' ends ANSI color sequence
				inEscape = false
			}
		} else {
			length++
		}
	}

	return length
}

// truncateAtVisualWidth truncates a string at a visual width, preserving ANSI codes
func truncateAtVisualWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	visualPos := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		} else if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		// Count visible character
		visualPos++

		// If we've reached max width, truncate here
		if visualPos >= maxWidth {
			return s[:i+1]
		}
	}

	return s
}

// findBreakPoint finds a good place to break a line (at spaces or punctuation)
func findBreakPoint(s string, maxWidth int) int {
	if maxWidth <= 0 {
		return 0
	}

	visualPos := 0
	lastSpace := -1
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		} else if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		// Count visible character
		visualPos++

		// Track spaces as potential break points
		if s[i] == ' ' || s[i] == '\t' || s[i] == '-' || s[i] == ',' || s[i] == '.' {
			lastSpace = i
		}

		// If we've reached max width
		if visualPos >= maxWidth {
			// Break at last space if possible
			if lastSpace > 0 && lastSpace > i-20 { // Within last 20 chars
				return lastSpace + 1
			}
			// Otherwise break here
			return i
		}
	}

	return len(s)
}

// highlightSearchMatches highlights search terms occurrences in a line
func highlightSearchMatches(line, searchTerm string) string {
	if searchTerm == "" {
		return line
	}

	// Case-insensitive search
	lowerLine := strings.ToLower(line)
	lowerTerm := strings.ToLower(searchTerm)

	if !strings.Contains(lowerLine, lowerTerm) {
		return line
	}

	// Build hightlighted version
	var result strings.Builder
	seearchLen := len(searchTerm)
	pos := 0

	for {
		idx := strings.Index(lowerLine[pos:], lowerTerm)
		if idx == -1 {
			result.WriteString(line[pos:])
			break
		}

		actualIdx := pos + idx
		// Write text before match
		result.WriteString(line[pos:actualIdx])
		// Write highlighted match (yellow background, black text)
		result.WriteString("\x1b[43m\x1b[30m")
		result.WriteString(line[actualIdx : actualIdx+seearchLen])
		result.WriteString("\x1b[0m")

		pos = actualIdx + seearchLen
	}

	return result.String()
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
	wrapStatus := "Wrap: OFF"
	if fv.WrapLines {
		wrapStatus = "Wrap: ON"
	}
	info := fmt.Sprintf("Lines: %d | Position: %d | %s", len(fv.Content), fv.ScrollPos+1, wrapStatus)
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

	linesRendered := 0
	for i := visibleStart; i < visibleEnd && linesRendered < maxVisible; i++ {
		if i >= len(contentToDisplay) {
			break
		}

		line := contentToDisplay[i]

		// Apply search highlighting if active
		if fv.SearchTerm != "" {
			line = highlightSearchMatches(line, fv.SearchTerm)
		}

		lineNum := fmt.Sprintf("%4d │ ", i+1)

		if fv.WrapLines {
			// Wrap the line if wrapping is enabled
			wrappedLines := wrapLine(line, fv.Width, i+1)

			// Render first line with line number
			if len(wrappedLines) > 0 {
				b.WriteString(lineNum + wrappedLines[0] + "\n")
				linesRendered++
			}

			// Render continuation lines with indentation
			for j := 1; j < len(wrappedLines) && linesRendered < maxVisible; j++ {
				b.WriteString("     ╎ " + wrappedLines[j] + "\n")
				linesRendered++
			}
		} else {
			// No wrapping - truncate long lines with indicator
			visualLen := visualLength(line)
			availableWidth := fv.Width - 10 // Account for line numbers and margin

			if availableWidth > 0 && visualLen > availableWidth {
				// Truncate at visual width (accounting for ANSI codes)
				truncated := truncateAtVisualWidth(line, availableWidth-3)
				b.WriteString(lineNum + truncated + "...\n")
			} else {
				b.WriteString(lineNum + line + "\n")
			}
			linesRendered++
		}
	}

	// Footer with help
	b.WriteString("\n")

	if fv.CommandMode {
		// Show command prompt
		commandPrompt := fmt.Sprintf(":%s", fv.CommandBuffer)
		b.WriteString(commandPrompt)
	} else if fv.StatusMessage != "" {
		// Show status message
		status := statusStyle.Render(fv.StatusMessage)
		b.WriteString(status + "\n")
		help := helpStyle.Render("↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u/d: page | :: command | q/Esc: back")
		b.WriteString(help)
	} else {
		// Show normal help
		help := helpStyle.Render("↑/k: up | ↓/j: down | g: top | G: bottom | Ctrl+u/d: page | :: command | q/Esc: back")
		b.WriteString(help)
	}

	return b.String()
}
