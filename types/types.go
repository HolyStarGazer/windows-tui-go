package types

// FileItem represents a file or directory in the file system
type FileItem struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}
