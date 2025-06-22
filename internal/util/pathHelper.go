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
	assetsPath := filepath.Join("assets/binaries", binName)
	if _, err := os.Stat(assetsPath); err == nil {
		fmt.Println("Path is --> ", assetsPath)
		return assetsPath
	}

	// Try system-installed binary (e.g., from Homebrew or PATH)
	if fullPath, err := exec.LookPath(name); err == nil {
		fmt.Println("Trying system installed --> ", assetsPath)
		return fullPath
	}

	// Not found
	fmt.Fprintf(os.Stderr, "Missing required binary: %s\n", binName)
	fmt.Fprintln(os.Stderr, "Please run: scripts/download-assets.(sh|ps1) or install via Homebrew.")
	return ""
}

// ModelPath returns the full path to a model file in assets/models.
// It returns an empty string and logs an error if the file does not exist.
func ModelPath(name string) string {

	modelPath := filepath.Join("assets", "models", name)
	if _, err := os.Stat(modelPath); err == nil {
		return modelPath
	} else if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Model not found: %s\n", modelPath)
	} else {
		fmt.Fprintf(os.Stderr, "Error checking model path: %s (%v)\n", modelPath, err)
	}

	return ""
}
