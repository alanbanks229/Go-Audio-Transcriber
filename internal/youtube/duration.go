// internal/youtube/duration.go
package youtube

import (
	"bufio" // For reading stderr line by line
	"encoding/json"
	"fmt"
	"io" // For io.ReadAll
	"os/exec"
	"strings"

	"fyne.io/fyne/v2" // For fyne.Do() and logging
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
)

// VideoInfo struct to unmarshal relevant parts of yt-dlp's JSON output
type VideoInfo struct {
	Duration float64 `json:"duration"` // Duration in seconds
}

// GetDuration fetches the duration of a YouTube video using yt-dlp.
// It returns the duration in seconds and an error, if any.
func GetDuration(url string, log func(string)) (float64, error) {

	// Command to get video info as JSON, but only fetch 'duration'
	cmd := exec.Command(util.BinPath("yt-dlp"), "--dump-json", "--flat-playlist", "--quiet", "--no-warnings", url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("Error getting stdout pipe for yt-dlp: %v", err))
		})
		return 0, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("Error getting stderr pipe for yt-dlp: %v", err))
		})
		return 0, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("Error starting yt-dlp: %v. Make sure yt-dlp is installed and in your system's PATH.", err))
		})
		return 0, fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	// Read stdout to get the JSON output
	jsonBytes, err := io.ReadAll(stdout)
	if err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("Error reading yt-dlp stdout: %v", err))
		})
		return 0, fmt.Errorf("failed to read yt-dlp output: %w", err)
	}

	// Read stderr in a separate goroutine to prevent deadlock if stderr is large
	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			line := errScanner.Text()
			fyne.Do(func() {
				// Log stderr output, especially if it's an error
				if strings.Contains(line, "error") || strings.Contains(line, "warning") {
					log(fmt.Sprintf("yt-dlp (stderr): %s", line))
				}
			})
		}
	}()

	if err := cmd.Wait(); err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("yt-dlp command failed: %v", err))
		})
		return 0, fmt.Errorf("yt-dlp command failed: %w", err)
	}

	var videoInfo VideoInfo
	if err := json.Unmarshal(jsonBytes, &videoInfo); err != nil {
		fyne.Do(func() {
			log(fmt.Sprintf("Error parsing yt-dlp JSON: %v. JSON output: %s", err, string(jsonBytes)))
		})
		return 0, fmt.Errorf("failed to parse yt-dlp JSON output: %w", err)
	}

	return videoInfo.Duration, nil
}
