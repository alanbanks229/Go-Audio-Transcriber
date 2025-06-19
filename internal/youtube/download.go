package youtube

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
)

// Download downloads audio (mp3) via yt-dlp and returns the local path.
func Download(url, outDir string, log func(string)) string {
	out := filepath.Join(outDir, "downloaded_"+util.RandString(6)+".%(ext)s")
	cmd := exec.Command(
		util.BinPath("yt-dlp"),
		"-x", "--audio-format", "mp3", "-o", out, url,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log(fmt.Sprintf("Error running yt-dlp: %v\nOutput:\n%s", err, string(output)))

		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 127 {
			log("yt-dlp is not installed or not available in PATH")
		}
		return ""
	}

	final := strings.TrimSuffix(out, ".%(ext)s") + ".mp3"
	log("Saved: " + final)
	return final
}
