package coding

import (
	"testing"

	"github.com/MichaelKinsy/PiG/coding/extension"
)

// The first session_start of a Runtime carries initial; the Sessions its
// replacements create do not (SPEC-0007 D2).
func TestRuntimeMarksOnlyItsFirstSessionStartInitial(t *testing.T) {
	h := newRuntimeTestHarness(t, runtimeTestOptions{})
	if got := h.runtime.Session().StartEvent(); got != (extension.SessionStartEvent{Type: "session_start", Reason: "startup", Initial: true}) {
		t.Fatalf("first session_start = %+v, want reason startup with initial", got)
	}
	runtimePrompt(t, h.runtime, "hello")
	previous := h.runtime.Session().Path()
	if result, err := h.runtime.NewSession(t.Context(), nil); err != nil || result.Cancelled {
		t.Fatalf("new session = %v, %v", result, err)
	}
	if got := h.runtime.Session().StartEvent(); got != (extension.SessionStartEvent{Type: "session_start", Reason: "new", PreviousSessionFile: previous}) {
		t.Fatalf("replacement session_start = %+v, want reason new without initial", got)
	}
}
