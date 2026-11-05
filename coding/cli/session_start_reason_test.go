//go:build linux

package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A process that opens an existing session used to report session_start reason
// "startup", so an extension could not tell a continued session from a fresh
// one, and no session_start said "the process just started" (SPEC-0007 D2).
// The first session_start of a Runtime carries initial; a replacement does not.
func TestSessionStartReportsResumeAndInitialAtProcessStart(t *testing.T) {
	if testing.Short() {
		t.Skip("starts the pig binary and an extension process")
	}
	bin := buildPigBinaryForSignalTest(t)
	fixture, err := filepath.Abs(filepath.Join("testdata", "session-start-reason.mjs"))
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	sessionDir := filepath.Join(home, "sessions")
	sessionFile := filepath.Join(home, "named.jsonl")
	work := t.TempDir()
	run := func(t *testing.T, args ...string) []string {
		t.Helper()
		log := filepath.Join(t.TempDir(), "start.log")
		cmd := exec.Command(bin, append([]string{"--model", "test-faux/faux-1", "--print", "--session-dir", sessionDir, "-e", fixture}, args...)...)
		cmd.Dir = work
		cmd.Env = append(os.Environ(), "PIG_HOME="+home, "PIG_CODING_AGENT_DIR="+filepath.Join(home, "agent"),
			"PIG_TEST_FAUX=1", "PIG_TEST_REPLACE_LOG="+log)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("pig %v: %v\n%s", args, err, out)
		}
		return replaceLogLines(t, log)
	}
	expect := func(t *testing.T, got []string, want ...string) {
		t.Helper()
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Fatalf("session_start events:\n got %q\nwant %q", got, want)
		}
	}

	expect(t, run(t, "What is 20+22?"), "session_start reason=startup initial=true")
	// A --session path that does not exist yet selects a new Session.
	expect(t, run(t, "--session", sessionFile, "What is 20+22?"), "session_start reason=startup initial=true")
	expect(t, run(t, "--session", sessionFile, "What is 20+22?"), "session_start reason=resume initial=true")
	expect(t, run(t, "--continue", "What is 20+22?"), "session_start reason=resume initial=true")
	expect(t, run(t, "--continue", "/replace-new"),
		"session_start reason=resume initial=true",
		"session_start reason=new initial=false")

	// RPC mode binds through its own host; it reports the same event.
	log := filepath.Join(t.TempDir(), "rpc.log")
	p := startRPCProcessAt(t, work, []string{"PIG_HOME=" + home, "PIG_CODING_AGENT_DIR=" + filepath.Join(home, "agent"), "PIG_TEST_FAUX=1", "PIG_TEST_REPLACE_LOG=" + log},
		"--model", "test-faux/faux-1", "--session-dir", sessionDir, "-e", fixture, "--continue")
	p.sendJSON(map[string]any{"id": "state", "type": "get_state"})
	p.await("get_state response", func(r rpcRecord) bool { return isSuccessResponse(r, "state") })
	p.closeAndWait("after get_state")
	expect(t, replaceLogLines(t, log), "session_start reason=resume initial=true")
}

// Interactive mode binds its first Session through its own rebind path, which fell back to a fixed startup event; it reports the Runtime's event too (SPEC-0007 D2).
func TestInteractiveSessionStartReportsResumeAndInitial(t *testing.T) {
	if testing.Short() {
		t.Skip("starts the pig binary and an extension process")
	}
	t.Parallel()
	fixture, err := filepath.Abs(filepath.Join("testdata", "session-start-reason.mjs"))
	if err != nil {
		t.Fatal(err)
	}
	expect := func(t *testing.T, p *interactivePig, want ...string) {
		t.Helper()
		if got := replaceLogLines(t, p.log); strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Fatalf("session_start events:\n got %q\nwant %q\nscreen tail: %q", got, want, tail(p.output.since(0), 600))
		}
	}

	t.Run("fresh", func(t *testing.T) {
		p := startInteractivePigWith(t, fixture)
		p.command("/new")
		waitReplaceLog(t, p.log, 2, "/new")
		p.quit()
		expect(t, p, "session_start reason=startup initial=true", "session_start reason=new initial=false")
	})

	t.Run("continue", func(t *testing.T) {
		p := &interactivePig{t: t, log: filepath.Join(t.TempDir(), "replace.log"), output: &ptyOutput{}, home: shortTempDir(t), dir: t.TempDir()}
		seed := exec.Command(buildPigBinaryForSignalTest(t), "--model", "test-faux/faux-1", "--print", "--session-dir", filepath.Join(p.home, "sessions"), "What is 20+22?")
		seed.Dir = p.dir
		seed.Env = append(os.Environ(), "PIG_HOME="+p.home, "PIG_CODING_AGENT_DIR="+filepath.Join(p.home, "agent"), "PIG_TEST_FAUX=1")
		if out, err := seed.CombinedOutput(); err != nil {
			t.Fatalf("seed session: %v\n%s", err, out)
		}
		startInteractivePigIn(t, p, fixture, "--continue")
		p.quit()
		expect(t, p, "session_start reason=resume initial=true")
	})
}
