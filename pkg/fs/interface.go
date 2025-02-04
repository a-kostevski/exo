package fs

type FileSystem interface {
	EnsureDirectoryExists(path string) error
	WriteFile(path string, content []byte) error
	ReadFile(path string) ([]byte, error)
	FileExists(path string) bool
	DeleteFile(path string) error
	OpenInEditor(path, editor string) error
}
