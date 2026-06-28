// internal/youtube/duration.go
package youtube

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2" // For fyne.Do() and logging
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
)

// GetDuration fetches the duration of a YouTube video using yt-dlp.
// It returns the duration in seconds and an error, if any.
func GetDuration(url string, log func(string)) (float64, error) {
	cmd := exec.Command(
		util.BinPath("yt-dlp"),
		"--print", "duration",
		"--skip-download",
		"--ignore-no-formats-error",
		"--no-warnings",
		"--no-playlist",
		url,
	)

	output, err := cmd.CombinedOutput()
	outputText := strings.TrimSpace(string(output))
	if err != nil {
		msg := fmt.Sprintf("yt-dlp duration probe failed: %v", err)
		if outputText != "" {
			msg = fmt.Sprintf("%s: %s", msg, outputText)
		}
		fyne.Do(func() {
			log(msg)
		})
		return 0, errors.New(msg)
	}

	durationText := firstNonEmptyLine(outputText)
	duration, err := strconv.ParseFloat(durationText, 64)
	if err != nil || duration <= 0 {
		fyne.Do(func() {
			log(fmt.Sprintf("Error parsing yt-dlp duration output %q: %v", outputText, err))
		})
		if err == nil {
			err = fmt.Errorf("duration must be greater than zero")
		}
		return 0, fmt.Errorf("failed to parse yt-dlp duration output %q: %w", outputText, err)
	}

	return duration, nil
}

func firstNonEmptyLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
