package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ResolveAssetPath returns a full path to an asset file
// Works in local dev and packaged builds (e.g., .app on macOS)
func ResolveAssetPath(relPath string) string {
	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)

	// macOS .app bundle layout
	if runtime.GOOS == "darwin" && strings.Contains(execPath, ".app/Contents/MacOS") {
		return filepath.Join(baseDir, "../Resources/assets", relPath)
	}

	// Development mode (go run or go build from root)
	// project root: baseDir is root, or ./dist when go build -o ./dist/...
	devPath := filepath.Join(".", "assets", relPath)
	if _, err := os.Stat(devPath); err == nil {
		return devPath
	}

	// Fallback: assume assets bundled beside binary
	return filepath.Join(baseDir, "assets", relPath)
}
