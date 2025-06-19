package audio

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
)

// DurationSeconds returns length of file or 0 on error (URL returns 0 – handled by caller)
// DurationSeconds returns the duration of an audio file in seconds as a float64.
func DurationSeconds(filepath string) float64 {
	// Example using ffprobe to get duration. Adjust command if yours is different.
	// You might need to make sure ffprobe is installed on the user's system.
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filepath,
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running ffprobe: %v\n", err)
		return 0 // Return 0 on error
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		fmt.Printf("Error parsing duration from ffprobe output '%s': %v\n", durationStr, err)
		return 0 // Return 0 on parsing error
	}

	return duration
}

// Trim copies sub-range into new file using ffmpeg (copy codec → fast)
func Trim(src string, start, end float64, dir string) string {
	out := filepath.Join(dir, "trimmed_"+util.RandString(6)+".mp3")
	cmd := exec.Command("ffmpeg",
		"-y", "-i", src, "-ss", fmt.Sprintf("%f", start),
		"-to", fmt.Sprintf("%f", end), "-c", "copy", out)
	_ = cmd.Run()
	return out
}
