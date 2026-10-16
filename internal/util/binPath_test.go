package util

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestModelAssetKeepsBinExtension(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("assets", 0755); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("assets", "ggml-small.en.bin")
	if err := os.WriteFile(want, []byte("model fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := ResolveAssetPath("ggml-small.en.bin"); got != want {
		t.Fatalf("model path = %q, want %q", got, want)
	}
}

func TestBinPathUsesBundledExecutable(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("assets", 0755); err != nil {
		t.Fatal(err)
	}
	name := "whisper-cli"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	want, err := filepath.Abs(filepath.Join("assets", name))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(want, nil, 0700); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{"whisper-cli", name} {
		if got := BinPath(input); got != want {
			t.Fatalf("BinPath(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestMissingBinaryDoesNotExit(t *testing.T) {
	t.Chdir(t.TempDir())
	got := BinPath("go-audio-transcriber-nonexistent-tool")
	if !filepath.IsAbs(got) {
		t.Fatalf("missing tool should return an absolute path, got %q", got)
	}
	if _, err := os.Stat(got); !os.IsNotExist(err) {
		t.Fatalf("expected missing tool, got %v", err)
	}
}
