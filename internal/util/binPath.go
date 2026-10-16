package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// BinPath resolves executables only. Use ResolveAssetPath for model/data files.
// Missing tools are returned as paths so callers can report command errors
// without terminating the GUI.
func BinPath(name string) string {
	// Add .exe for Windows
	binName := name
	if runtime.GOOS == "windows" && !strings.EqualFold(filepath.Ext(name), ".exe") {
		binName += ".exe"
	}

	// Check ./assets/{name}
	assetsPath := ResolveAssetPath(binName)
	if absolute, err := filepath.Abs(assetsPath); err == nil {
		assetsPath = absolute
	}
	if info, err := os.Stat(assetsPath); err == nil && !info.IsDir() {
		return assetsPath
	}

	// Try system-installed binary (e.g., from Homebrew or PATH)
	if fullPath, err := exec.LookPath(name); err == nil {
		return fullPath
	}

	return assetsPath
}
