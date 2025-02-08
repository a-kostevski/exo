package fs

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands a leading tilde (~) in the provided path to the user's home directory.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path // fallback to the original path if home cannot be determined
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// ResolvePath returns an absolute path by joining base with path if path is not absolute.
func ResolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

// GetXDGConfigHome returns the XDG_CONFIG_HOME directory, or defaults to $HOME/.config.
func GetXDGConfigHome() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config")
}

// GetXDGDataHome returns the XDG_DATA_HOME directory, or defaults to $HOME/.local/share.
func GetXDGDataHome() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share")
}

// GetXDGCacheHome returns the XDG_CACHE_HOME directory, or defaults to $HOME/.cache.
func GetXDGCacheHome() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache")
}

// SanitizePath cleans the provided path after expanding any tilde. If the result is not absolute,
// it is joined with the provided home directory.
func SanitizePath(path, home string) string {
	expanded := ExpandPath(path)
	cleaned := filepath.Clean(expanded)
	if !filepath.IsAbs(cleaned) {
		cleaned = filepath.Join(home, cleaned)
	}
	return cleaned
}
