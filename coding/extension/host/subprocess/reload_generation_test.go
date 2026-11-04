package subprocess

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestReloadDeadlineDoesNotSelectAnotherRuntime(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{"register deadline", newLoadError("fixture", "register", "register_failed", context.DeadlineExceeded), false},
		{"register cancellation", newLoadError("fixture", "register", "register_failed", context.Canceled), false},
		{"admission deadline", &admissionWriteError{err: context.DeadlineExceeded}, false},
		{"packed deadline", &packedAcceptError{perMember: map[string]error{"fixture": newLoadError("fixture", "register", "register_failed", context.DeadlineExceeded)}}, false},
		{"runtime registration failure", newLoadError("fixture", "register", "register_failed", errors.New("invalid frame")), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := isStartFailure(tc.err); got != tc.want {
				t.Fatalf("runtime fallback classification = %v, want %v: %v", got, tc.want, tc.err)
			}
		})
	}
}

func TestReloadHandshakeStopsStalledProcess(t *testing.T) {
	bin := buildFixtureExt(t)
	for _, phase := range []string{"connect", "register"} {
		t.Run(phase, func(t *testing.T) {
			t.Setenv("FIXTURE_HANDSHAKE_MODE", phase)
			h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
			t.Cleanup(func() { h.Shutdown("test done") })
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			ctx = context.WithValue(ctx, runtimeParentKey{}, t.Context())
			ctx = context.WithValue(ctx, reloadingKey{}, uint64(1))
			ctx = context.WithValue(ctx, reloadHandshakeTimeoutKey{}, 250*time.Millisecond)
			me := &managedExt{config: ExtConfig{Name: "fixture", Path: bin}, host: h}
			_, err := h.startExt(ctx, me, false)
			if !errors.Is(err, context.DeadlineExceeded) || !strings.Contains(err.Error(), "reload handshake timed out") {
				t.Fatalf("reload error = %v, want the handshake deadline", err)
			}
			if ctx.Err() != nil {
				t.Fatal("reload waited for the caller's deadline")
			}
			loadErr, ok := errors.AsType[*LoadError](err)
			if !ok || loadErr.Phase != phase {
				t.Fatalf("error = %#v, want phase %s", err, phase)
			}
			if loadErr.Hint != "" {
				t.Fatalf("reload deadline misclassified as an SDK or payload failure: %s", loadErr.Hint)
			}
			select {
			case <-me.exited():
			default:
				t.Fatal("timed-out process was not killed and reaped")
			}
		})
	}
}

func TestReloadDeadlineIncludesStalledRetainedAdmission(t *testing.T) {
	root := writeHeightReportingModule(t, "example.com/stalledadmission", "stalled-admission")
	cfg := packedFactoryConfig("stalled-admission", root, "example.com/stalledadmission", "admission")
	h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
	t.Cleanup(func() { h.Shutdown("test done") })
	h.SetHeightFunc(func() int { return 42 })
	h.SetConfigLoader(func() ([]ExtConfig, error) { return []ExtConfig{cfg}, nil })
	h.SetReloadHandshakeTimeoutFunc(func() (time.Duration, error) { return 250 * time.Millisecond, nil })
	if _, errs := h.LoadAll(t.Context(), []ExtConfig{cfg}); len(errs) != 0 {
		t.Fatal(errs)
	}
	old := h.exts[cfg.Name]
	share := old.packedProcess.share
	writer, peer := net.Pipe()
	original := share.admit
	share.admit = writer
	t.Cleanup(func() {
		h.Shutdown("test done")
		_ = original.Close()
		_ = peer.Close()
	})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx = context.WithValue(ctx, runtimeParentKey{}, t.Context())
	result := make(chan error, 1)
	var work sync.WaitGroup
	work.Go(func() { _, err := h.Reload(ctx); result <- err })
	var reloadErr error
	select {
	case reloadErr = <-result:
	case <-time.After(5 * time.Second):
		t.Error("reload did not bound the retained process's admission write")
		cancel()
		_ = peer.Close()
		reloadErr = <-result
	}
	work.Wait()
	if !errors.Is(reloadErr, context.DeadlineExceeded) || !strings.Contains(reloadErr.Error(), "reload handshake timed out") {
		t.Fatalf("admission error = %v, want the reload deadline", reloadErr)
	}
	if h.exts[cfg.Name] != old {
		t.Fatal("failed admission discarded the serving generation")
	}
	assertPackedHeight(t, t.Context(), h, 42)
	h.Shutdown("verify writer drain")
	select {
	case <-share.exitedChan():
	default:
		t.Fatal("shutdown left the retained process running")
	}
}

func TestInitialLoadHasNoReloadHandshakeDeadline(t *testing.T) {
	h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
	t.Cleanup(func() { h.Shutdown("test done") })
	ctx := context.WithValue(t.Context(), reloadHandshakeTimeoutKey{}, time.Nanosecond)
	if _, err := h.Load(ctx, ExtConfig{Name: "fixture", Path: buildFixtureExt(t)}); err != nil {
		t.Fatalf("initial load applied a reload-only deadline: %v", err)
	}
}

func TestSocketGenerationTeardownDoesNotUnlinkSuccessor(t *testing.T) {
	h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
	t.Cleanup(func() { h.Shutdown("test done") })
	oldPath, err := h.sockPathFor("reload-member")
	if err != nil {
		t.Fatal(err)
	}
	old := &managedExt{config: ExtConfig{Name: "reload-member"}, sockPath: oldPath}
	newPath, err := h.sockPathFor("reload-member")
	if err != nil {
		t.Fatal(err)
	}
	listener, address, err := ListenExtension(newPath, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	h.stopFailedPackedMember(old)
	conn, err := (&net.Dialer{}).DialContext(t.Context(), "unix", address)
	if err != nil {
		t.Fatalf("old generation teardown removed the successor's socket: %v", err)
	}
	_ = conn.Close()
}

func TestPackedGoReloadLoopNeverHangs(t *testing.T) {
	root := writeHeightReportingModule(t, "example.com/reloadloop", "reload-loop")
	path := filepath.Join(root, "extension.go")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	source = []byte(strings.Replace(string(source), `e := sdk.New("reload-loop")`, `e := sdk.New("reload-loop")
	if _, err := sdk.ParseFrontmatter("---\nkey: value\n---\nbody"); err != nil { panic(err) }
	if _, err := e.Events().On("reload-loop", func(sdk.Context, any) error { return nil }); err != nil { panic(err) }`, 1))
	if err := os.WriteFile(path, source, 0o600); err != nil {
		t.Fatal(err)
	}
	h := NewHostWithConfigRoot(t.TempDir(), t.TempDir())
	h.SetConfigLoader(func() ([]ExtConfig, error) {
		return []ExtConfig{packedFactoryConfig("reload-loop", root, "example.com/reloadloop", "loop")}, nil
	})
	t.Cleanup(func() { h.Shutdown("test done") })
	h.SetHeightFunc(func() int { return 42 })
	ctx, cancel := buildContext(t)
	defer cancel()
	if _, err := h.Reload(ctx); err != nil {
		t.Fatal(err)
	}
	assertPackedHeight(t, ctx, h, 42)
	h.mu.Lock()
	share := h.exts["reload-loop"].packedProcess.share
	h.mu.Unlock()
	for i := range 200 {
		step, cancel := context.WithTimeout(ctx, 10*time.Second)
		// A reload deadline bounds admission, not the session's retained process lifetime.
		step = context.WithValue(step, runtimeParentKey{}, ctx)
		_, err := h.Reload(step)
		if err == nil {
			assertPackedHeight(t, step, h, 42)
		}
		cancel()
		if err != nil {
			t.Fatalf("reload %d: %v", i+1, err)
		}
		h.mu.Lock()
		current := h.exts["reload-loop"].packedProcess.share
		h.mu.Unlock()
		if current != share {
			t.Fatalf("reload %d replaced the retained process", i+1)
		}
	}
	h.Shutdown("loop complete")
	select {
	case <-share.exitedChan():
	default:
		t.Fatal("shutdown left the retained process running")
	}
	if _, err := os.Stat(h.sockRuntimeDir); !os.IsNotExist(err) {
		t.Fatalf("socket directory remains after shutdown: %v", err)
	}
}
