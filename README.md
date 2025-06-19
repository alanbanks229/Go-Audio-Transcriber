# 🎙️ Audio Transcription Tool

A desktop app built with [Fyne](https://fyne.io) in Go for converting **YouTube videos** or **MP3 files** into **text transcripts**.

---

# 🚀 Features
- Paste a YouTube URL or select a local MP3
- Choose an output folder for the transcript
- Adjust start and end times via a custom range slider
- Displays real-time transcription logs
- Works fully offline with local MP3s


---

# 📦 Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-username/audio-transcription-tool.git
cd audio-transcription-tool
```

### 2. Build the application
```bash
go build -o transcriber main.go
./transcriber
```
⚠️ Ensure you have access to ffmpeg and yt-dlp on your system's $PATH. These are used for audio trimming and YouTube downloading.

---


### 🔧 Dependencies

- Fyne
- yt-dlp (for downloading YouTube audio)
- ffmpeg (for trimming audio)

### 💡 Ideas for Future Versions

- Subtitle generation (.srt)
- Export to .pdf or .docx
- Drag-and-drop support
- Language selection for transcription
