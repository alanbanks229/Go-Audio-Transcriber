package transcribe

import (
	"strings"
	"testing"
)

func TestMissingModelReturnsErrorWithoutExiting(t *testing.T) {
	t.Chdir(t.TempDir())
	var messages []string
	got := Run("unused.mp3", ".", nil, func(message string) {
		messages = append(messages, message)
	})
	if got != "" {
		t.Fatalf("missing model returned output path %q", got)
	}
	if len(messages) != 1 || !strings.Contains(messages[0], "model is missing or inaccessible") {
		t.Fatalf("expected actionable model error, got %v", messages)
	}
}
