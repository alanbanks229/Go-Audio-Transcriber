package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BinPath returns a usable path to a binary in ./assets, falling back to system-installed binary.
func BinPath(name string) string {
	// Add .exe for Windows
	binName := name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// Check ./assets/{name}
	assetsPath := filepath.Join("assets", binName)
	if _, err := os.Stat(assetsPath); err == nil {
		return assetsPath
	}

	// Try system-installed binary (e.g., from Homebrew or PATH)
	if fullPath, err := exec.LookPath(name); err == nil {
		return fullPath
	}

	// Not found
	fmt.Fprintf(os.Stderr, "Missing required binary: %s\n", binName)
	fmt.Fprintln(os.Stderr, "Please run: scripts/download-assets.(sh|ps1) or install via Homebrew.")
	os.Exit(1)
	return "" // unreachable
}
