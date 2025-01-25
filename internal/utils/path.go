package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands the tilde in paths to the user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// ResolvePath resolves a path relative to a base directory
func ResolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

// GetXDGConfigHome returns the XDG config home directory
func GetXDGConfigHome() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return xdgConfig
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config")
}

// GetXDGDataHome returns the XDG data home directory
func GetXDGDataHome() string {
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return xdgData
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share")
}

// GetXDGCacheHome returns the XDG cache home directory
func GetXDGCacheHome() string {
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		return xdgCache
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache")
}

// SanitizePath cleans and normalizes the provided path, expanding any relative paths
// against the provided base directory (usually home directory)
func SanitizePath(path, base string) string {
	expanded := ExpandPath(path)
	cleaned := filepath.Clean(expanded)
	if !filepath.IsAbs(cleaned) {
		cleaned = filepath.Join(base, cleaned)
	}
	return cleaned
}
