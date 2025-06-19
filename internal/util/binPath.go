package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// BinPath returns the full relative path to a binary inside the assets folder.
// It appends .exe automatically on Windows.
// It also verifies existence and prints a friendly error if the binary is missing.
func BinPath(name string) string {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join("assets", name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Missing required binary: %s\n", path)
		fmt.Fprintln(os.Stderr, "💡 Please run: scripts/download-assets.(sh|ps1)")
		os.Exit(1)
	}

	return path
}
