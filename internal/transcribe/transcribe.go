package transcribe

import (
	"bufio" // Import io for MultiReader
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
)

// Run uses whisper.cpp to transcribe an mp3 and returns the output .txt path.
func Run(mp3, dir string, pb *widget.ProgressBar, log func(string)) string {
	baseName := "transcript_" + util.RandString(6)
	outBase := filepath.Join(dir, baseName)

	cmd := exec.Command(
		util.BinPath("whisper-cli"),
		"-m", util.BinPath("ggml-small.en.bin"),
		"-f", mp3,
		"-otxt",
		"-of", outBase,
		"-pp", // Enable progress reporting (to stderr)
		// Consider adding -l 0 to reduce logging if it's too verbose even on stderr
		// "-l", "0",
	)

	// Capture both stdout and stderr for processing
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fyne.Do(func() {
			log("Error: Failed to get stdout pipe: " + err.Error())
		})
		return ""
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fyne.Do(func() {
			log("Error: Failed to get stderr pipe: " + err.Error())
		})
		return ""
	}

	if err := cmd.Start(); err != nil {
		fyne.Do(func() {
			log("Error: Failed to start whisper.cpp: " + err.Error())
		})
		return ""
	}

	// Regex to find progress percentage in stderr output
	progressRegex := regexp.MustCompile(`progress = (\s*)(\d+)%`) // Added \s* for optional space

	// --- Goroutine for reading stderr (progress updates and diagnostic logs) ---
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fyne.Do(func() {
				// Attempt to find progress percentage
				matches := progressRegex.FindStringSubmatch(line)
				if len(matches) > 2 { // We now expect 3 matches: full string, spaces, digits
					if p, err := strconv.Atoi(matches[2]); err == nil { // Use matches[2] for the digits
						pb.SetValue(float64(p) / 100.0)
						// Log progress only if you want to see it explicitly in the debug log
						// log(fmt.Sprintf("Progress: %d%%", p))
					}
				} else {
					// Log other stderr output for debugging purposes, but consider filtering
					// For example, you might only want to log lines that indicate errors or important warnings.
					// If you want minimal logs, you could completely remove this 'else' block for general stderr.
					if strings.Contains(line, "error") || strings.Contains(line, "fail") || strings.Contains(line, "warning") {
						log("Whisper.cpp (stderr): " + strings.TrimSpace(line))
					}
					// Otherwise, if it's just informational, don't log it to the user.
				}
			})
		}
		if err := scanner.Err(); err != nil {
			fyne.Do(func() {
				log("Error reading stderr pipe: " + err.Error())
			})
		}
	}()

	// --- Goroutine for reading stdout (transcription output) ---
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fyne.Do(func() {
				// Only log the actual transcription lines to the user-facing log
				// You might want to filter out empty lines or purely timestamp lines if they appear without text
				if strings.Contains(line, "-->") || strings.TrimSpace(line) != "" { // Basic check for transcription line
					log("Transcription: " + strings.TrimSpace(line))
				}
			})
		}
		if err := scanner.Err(); err != nil {
			fyne.Do(func() {
				log("Error reading stdout pipe: " + err.Error())
			})
		}
	}()

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fyne.Do(func() {
			log("Error: whisper.cpp command finished with error: " + err.Error())
		})
		// Depending on your error handling, you might return early here
	}

	fyne.Do(func() {
		pb.SetValue(1.0) // Ensure it's 100% when done
		log("Transcription process finished.")
	})

	return outBase + ".txt"
}
