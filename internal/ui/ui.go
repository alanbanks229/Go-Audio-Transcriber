package ui

import (
	"fmt" // Import fmt for error messages
	"image/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/alanbanks229/Go-Audio-Transcriber/internal/audio"
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/transcribe"
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/util"
	"github.com/alanbanks229/Go-Audio-Transcriber/internal/youtube" // Ensure this is imported
)

func CreateUserInterface(fyneWindow fyne.Window) {

	fyneWindow.CenterOnScreen()

	// Header
	darkIndigo := color.NRGBA{R: 18, G: 24, B: 41, A: 255}
	title := canvas.NewText("Audio Transcription Tool", darkIndigo)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 24

	subtitle := canvas.NewText("Convert YouTube videos or MP3 files to text transcripts", color.NRGBA{R: 95, G: 102, B: 115, A: 255})
	subtitle.TextSize = 14

	// Input Section
	inputHeader := widget.NewLabelWithStyle("YouTube Link or MP3 File Path", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("YouTube URL or select file...")
	selectMP3Btn := widget.NewButton(
		"                     Select MP3 File                     ",
		// onClick MP3 function
		func() {
			fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, _ error) {
				if r != nil {
					inputEntry.SetText(r.URI().Path())
				}
			}, fyneWindow)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
			fd.Show()
		},
	)

	// Create a VBox for the input entry and button (to put them in different rows)
	inputAndButtonVBox := container.NewVBox(
		inputEntry,
		util.VSpacer(4),
		selectMP3Btn,
		util.VSpacer(4),
	)

	inputSection := util.BorderedContainer(container.NewVBox(
		inputHeader,
		container.NewHBox(
			util.HSpacer(12),
			container.NewStack(inputAndButtonVBox),
			util.HSpacer(12),
		),
	))

	// Output Section
	outPath := binding.NewString()
	home, err := os.UserHomeDir() // cross-platform helper
	if err != nil {               // very unlikely, but fall back to cwd
		home = "."
	}
	downloads := filepath.Join(home, "Downloads")
	_ = outPath.Set(downloads)
	outputHeader := widget.NewLabelWithStyle("Output Directory", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder(downloads)
	selectOutputBtn := widget.NewButton("                 Select Output Folder                 ", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, _ error) {
			if uri != nil {
				_ = outPath.Set(uri.Path())
				outputEntry.SetText(uri.Path())
			}
		}, fyneWindow).Show()
	})

	outputAndButtonVBox := container.NewVBox(
		outputEntry,
		util.VSpacer(4),
		selectOutputBtn,
		util.VSpacer(4),
	)

	outputSection := util.BorderedContainer(container.NewVBox(
		outputHeader,
		container.NewHBox(
			util.HSpacer(12),
			container.NewStack(outputAndButtonVBox),
			util.HSpacer(12),
		),
	))

	inputFixed := container.NewStack(inputSection)
	outputFixed := container.NewStack(outputSection)

	// Side-by-side row
	inputOutputRow := container.NewHBox(
		inputFixed,
		util.HSpacer(16),
		outputFixed,
	)

	// Time Range Section
	const min = 0
	const max = 60
	const start = 0
	const end = 60
	rangeSlider := NewRangeSlider(min, max, start, end) // replaced once we know duration

	timeRangeHeader := widget.NewLabelWithStyle("Time Range Selection", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	startLbl := widget.NewLabel("Start: 00:00")
	endLbl := widget.NewLabel("End: 01:00")

	rangeSlider.OnChanged = func() {
		startLbl.SetText("Start: " + util.FormatSec(rangeSlider.Start))
		endLbl.SetText("End: " + util.FormatSec(rangeSlider.End))
		rangeSlider.renderer.Layout(rangeSlider.Size()) // This is really scuffed
	}
	// rangeSlider.Refresh()
	timeRangeSection := util.BorderedContainer(container.NewVBox(
		timeRangeHeader,
		container.New(layout.NewBorderLayout(nil, nil, startLbl, endLbl), startLbl, endLbl),
		container.New(layout.NewCenterLayout(), rangeSlider),
	))

	// Progress Section
	progressHeader := widget.NewLabelWithStyle("Transcription Progress", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	// Multi-line Transcription Section
	logBox := widget.NewMultiLineEntry()
	logBox.SetMinRowsVisible(8)
	logBox.Wrapping = fyne.TextWrapWord
	logBox.SetPlaceHolder("Transcription will appear here...")

	// Light background behind log box
	logBoxBG := canvas.NewRectangle(color.NRGBA{R: 245, G: 245, B: 245, A: 255}) // soft light gray

	// Set fixed height using a layout
	logBoxContainer := container.NewMax(logBoxBG, logBox)
	logBoxWithHeight := container.NewVBox(logBoxContainer)

	// Final Progress Section
	transcriptionProgressSection := util.BorderedContainer(container.NewVBox(
		progressHeader,
		progressBar,
		logBoxWithHeight,
	))

	// Transcript Button
	generateTranscriptBtn := widget.NewButtonWithIcon("Generate Transcript", theme.MediaPlayIcon(), func() {})
	generateTranscriptBtn.Disable()
	youtubeRX := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/`)
	generateTranscriptBtnContainer := util.BorderedContainer(generateTranscriptBtn)

	validateInputField := func() {
		txt := strings.TrimSpace(inputEntry.Text)
		if txt == "" {
			generateTranscriptBtn.Disable()
			// Reset importance when disabled
			generateTranscriptBtn.Importance = widget.LowImportance
			generateTranscriptBtn.Refresh() // Refresh the button to update its appearance
			return
		}

		if util.FileExists(txt) || youtubeRX.MatchString(txt) {
			generateTranscriptBtn.Enable()
			// Set importance to HighImportance for primary color
			generateTranscriptBtn.Importance = widget.HighImportance
			generateTranscriptBtn.Refresh() // Refresh the button to update its appearance
		} else {
			generateTranscriptBtn.Disable()
			// Reset importance when disabled
			generateTranscriptBtn.Importance = widget.LowImportance
			generateTranscriptBtn.Refresh() // Refresh the button to update its appearance
		}
	}

	// appendLog function for unified logging
	appendLog := func(s string) {
		fyne.Do(func() {
			logBox.SetText(logBox.Text + s + "\n")
		})
	}

	// This is the crucial part that needs modification
	inputEntry.OnChanged = func(text string) {
		validateInputField() // Re-validate button state based on new input

		src := strings.TrimSpace(text)
		if src == "" {
			// If input is empty, reset range slider to default and return
			fyne.Do(func() {
				// rangeSlider.SetRange(0, 60)
				// rangeSlider.SetStartEnd(0, 60)
				startLbl.SetText("Start: 00:00")
				endLbl.SetText("End: 01:00")
			})
			return
		}

		go func() { // Run duration check in a goroutine to keep UI responsive
			// Disable the button before starting the background task
			fyne.Do(func() {
				generateTranscriptBtn.Disable()
				// Optionally, reset importance to a neutral state (e.g., LowImportance)
				generateTranscriptBtn.Importance = widget.LowImportance
				generateTranscriptBtn.Refresh() // Make sure the UI updates
			})

			var dur float64
			var err error
			src := inputEntry.Text // Get the source URL/path

			if youtubeRX.MatchString(src) {
				fyne.Do(func() {
					logBox.SetText("")
					appendLog("Probing YouTube video for duration (this may take a moment)...")
				})
				// Call the new GetDuration function from the youtube package
				dur, err = youtube.GetDuration(src, appendLog)
				if err != nil {
					fyne.Do(func() {
						appendLog("Error getting YouTube duration: " + err.Error())
						dialog.ShowError(fmt.Errorf("could not get YouTube duration: %w", err), fyneWindow)
					})
					return // Stop if we can't get duration
				}
				fyne.Do(func() {
					appendLog(fmt.Sprintf("YouTube video duration: %s", util.FormatSec(dur)))
				})
			} else if util.FileExists(src) {
				fyne.Do(func() {
					appendLog("Probing local MP3 for duration...")
				})
				dur = audio.DurationSeconds(src)
				if dur == 0 {
					fyne.Do(func() {
						appendLog("Error: Could not get duration for local MP3. Is it a valid audio file?")
						dialog.ShowError(fmt.Errorf("could not get duration for local MP3"), fyneWindow)
					})
					return // Stop if duration is 0
				}
				fyne.Do(func() {
					appendLog(fmt.Sprintf("Local MP3 duration: %s", util.FormatSec(dur)))
				})
			} else {
				// Input is neither a valid YouTube link nor an existing file.
				// The validateInputField() already handled disabling the button.
				// No need to reset range slider here, as it was already reset
				// if the input was empty initially.
				return
			}

			// Update UI elements on the main goroutine, only if a valid duration (>0) was found
			if dur > 0 {
				fyne.Do(func() {
					// rangeSlider.SetRange(0, dur)
					// rangeSlider.SetStartEnd(0, dur) // Set full range as default selection
					startLbl.SetText("Start: " + util.FormatSec(rangeSlider.Start))
					endLbl.SetText("End: " + util.FormatSec(rangeSlider.End))
					// The renderer.Layout() call might still be needed if NewRangeSlider isn't
					// automatically laying out its children on range changes.
					// rangeSlider.renderer.Layout(rangeSlider.Size())

					// Re-enable the button in primary color if successful
					generateTranscriptBtn.Enable()
					generateTranscriptBtn.Importance = widget.HighImportance
					generateTranscriptBtn.Refresh()
				})
			}
		}()
	}

	// main action (remains largely the same, but benefits from mp3 path being available)
	generateTranscriptBtn.OnTapped = func() {
		generateTranscriptBtn.Disable()
		progressBar.SetValue(0) // Reset progress bar
		logBox.SetText("")      // Clear previous logs
		appendLog("Starting transcription process...")

		go func() {
			src := strings.TrimSpace(inputEntry.Text)
			targetDir, _ := outPath.Get() // outPath should be safely retrieved here

			var mp3Path string // Variable to hold the path to the local MP3 file

			if youtubeRX.MatchString(src) {
				appendLog("Downloading audio from YouTube…")
				downloadedMP3 := youtube.Download(src, targetDir, appendLog)
				if downloadedMP3 == "" {
					fyne.Do(func() {
						dialog.ShowError(fmt.Errorf("failed to download YouTube audio"), fyneWindow)
						generateTranscriptBtn.Enable()
						progressBar.SetValue(0)
					})
					return
				}
				mp3Path = downloadedMP3
			} else {
				// Assume it's already a local MP3 file
				mp3Path = src
			}

			// Ensure we have a valid MP3 path before proceeding
			if mp3Path == "" || !util.FileExists(mp3Path) {
				fyne.Do(func() {
					dialog.ShowError(fmt.Errorf("no valid audio file found for transcription"), fyneWindow)
					generateTranscriptBtn.Enable()
					progressBar.SetValue(0)
				})
				return
			}

			appendLog("Trimming audio…")
			trimmed := audio.Trim(mp3Path, rangeSlider.Start, rangeSlider.End, targetDir)
			if trimmed == "" {
				fyne.Do(func() {
					dialog.ShowError(fmt.Errorf("failed to trim audio"), fyneWindow)
					generateTranscriptBtn.Enable()
					progressBar.SetValue(0)
				})
				return
			}
			appendLog(fmt.Sprintf("Trimmed audio saved to: %s", trimmed))

			appendLog("Transcribing…")
			outFile := transcribe.Run(trimmed, targetDir, progressBar, appendLog)

			fyne.Do(func() {
				if outFile != "" && util.FileExists(outFile) {
					dialog.ShowInformation("Finished", "Transcript saved to:\n"+outFile, fyneWindow)
				} else {
					dialog.ShowError(fmt.Errorf("transcription failed or output file not found"), fyneWindow)
				}
				generateTranscriptBtn.Enable()
				progressBar.SetValue(0) // Reset progress bar
			})
		}()
	}

	// Final layout with defaultSpacing
	centerContent := container.NewVBox(
		title,
		subtitle,
		util.VSpacer(16),
		inputOutputRow,
		util.VSpacer(16),
		timeRangeSection,
		util.VSpacer(16),
		transcriptionProgressSection,
		util.VSpacer(16),
		generateTranscriptBtnContainer,
	)

	paddedAppContent := container.NewVBox(
		util.VSpacer(32), // Top padding
		container.NewHBox(
			util.HSpacer(32), // Left padding
			centerContent,
			util.HSpacer(32), // Right padding
		),
		util.VSpacer(16), // Bottom padding
	)

	// Add soft background for full app
	appBackground := canvas.NewRectangle(color.NRGBA{R: 249, G: 250, B: 251, A: 255}) // light gray
	appBackground.SetMinSize(fyne.NewSize(1, 1))

	fyneWindow.SetContent(container.NewMax(
		appBackground,
		paddedAppContent,
	))
}
