package subprocess

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"testing/synctest"
)

type gatedAdmitWriter struct {
	net.Conn
	started   chan struct{}
	closed    chan struct{}
	resume    <-chan struct{}
	startOnce sync.Once
	closeOnce sync.Once
}

func (w *gatedAdmitWriter) Write(p []byte) (int, error) {
	w.startOnce.Do(func() { close(w.started) })
	n, err := w.Conn.Write(p)
	<-w.resume
	return n, err
}

func (w *gatedAdmitWriter) Close() error {
	w.closeOnce.Do(func() { close(w.closed) })
	return w.Conn.Close()
}

func TestRetentionProbeDoesNotBlockAnotherProcessClaim(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		writer, peer := net.Pipe()
		defer func() { _ = peer.Close() }()
		blocked := newProcessShare("blocked", nil, make(chan struct{}), writer, nil)
		healthy := newProcessShare("healthy", nil, make(chan struct{}), discardAdmitWriter{io.Discard}, nil)
		retention := NewRuntimeRetention()
		retention.register(blocked)
		retention.register(healthy)
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		var work sync.WaitGroup
		first := make(chan error, 1)
		work.Go(func() {
			got, err := retention.claim(ctx, "blocked")
			if got != nil {
				got.release()
			}
			first <- err
		})
		synctest.Wait()
		type claimResult struct {
			share *processShare
			err   error
		}
		result := make(chan claimResult, 1)
		work.Go(func() { got, err := retention.claim(t.Context(), "healthy"); result <- claimResult{got, err} })
		synctest.Wait()
		var claimed *processShare
		select {
		case got := <-result:
			claimed = got.share
			if got.err != nil || claimed != healthy {
				t.Errorf("healthy claim = %v, %v", claimed, got.err)
			}
		default:
			t.Error("one stalled process blocked an unrelated retention claim")
		}
		cancel()
		synctest.Wait()
		work.Wait()
		if err := <-first; !errors.Is(err, context.Canceled) {
			t.Errorf("blocked claim error = %v", err)
		}
		if claimed == nil {
			select {
			case got := <-result:
				claimed = got.share
			default:
			}
		}
		if claimed != nil {
			claimed.release()
		}
		blocked.release()
		healthy.release()
	})
}

func TestProcessShareCancellationReturnsBeforeTheOwnerDrainsItsWriter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		writer, peer := net.Pipe()
		defer func() { _ = peer.Close() }()
		resume := make(chan struct{})
		gate := &gatedAdmitWriter{Conn: writer, started: make(chan struct{}), closed: make(chan struct{}), resume: resume}
		exited := make(chan struct{})
		share := newProcessShare("blocked", nil, exited, gate, func() { close(exited) })
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		sent := make(chan error, 1)
		var work sync.WaitGroup
		work.Go(func() { sent <- share.send(ctx, factoryAdmission{Op: "noop"}) })
		synctest.Wait()
		select {
		case <-gate.started:
		default:
			t.Error("admission did not begin its write")
		}
		cancel()
		synctest.Wait()
		select {
		case err := <-sent:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("send error = %v, want cancellation", err)
			}
		default:
			t.Error("cancellation waited for the blocked write")
		}
		released := make(chan struct{})
		work.Go(func() { share.release(); close(released) })
		synctest.Wait()
		select {
		case <-gate.closed:
		default:
			t.Error("the owner did not close its writer")
		}
		select {
		case <-released:
			t.Error("the owner returned before its writer finished")
		default:
		}
		close(resume)
		synctest.Wait()
		work.Wait()
		select {
		case <-released:
		default:
			t.Error("the owner did not finish after its writer drained")
		}
	})
}

type discardAdmitWriter struct{ io.Writer }

func (discardAdmitWriter) Close() error { return nil }

// A process that was killed but not yet reaped still takes control lines until the kernel closes its stdin, so neither the exit channel nor the noop probe refuses a claim on it. A member connection the host saw close unexpectedly marks the share dead, and claim must refuse it while the probe still succeeds.
func TestRetentionClaimRefusesShareWithClosedMemberConnection(t *testing.T) {
	t.Parallel()
	r := NewRuntimeRetention()
	share := newProcessShare("node:group", nil, make(chan struct{}), discardAdmitWriter{io.Discard}, nil)
	r.register(share)
	t.Cleanup(func() { share.release() })
	share.dead.Store(true)
	if got, err := r.claim(t.Context(), "node:group"); err != nil || got != nil {
		t.Fatalf("claim handed out a dead share: %v, %v", got, err)
	}
	share.dead.Store(false)
	got, err := r.claim(t.Context(), "node:group")
	if err != nil || got != share {
		t.Fatalf("claim refused a live share: %v, %v", got, err)
	}
	got.release()
}

// Every connection the host did not close withdraws its process from retention, packed or isolated: the isolated Node factory process is claimed by name and entry on the next reload exactly as a packed cell is claimed by group, and it too still takes control lines until it is reaped.
func TestUnexpectedConnectionCloseWithdrawsItsProcessFromRetention(t *testing.T) {
	t.Parallel()
	for _, packed := range []bool{false, true} {
		name := "isolated"
		if packed {
			name = "packed"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			hostEnd, peerEnd := net.Pipe()
			conn := NewConn("member", hostEnd)
			conn.Start(t.Context())
			h := NewHost(t.TempDir())
			key := "node-extension:member\x00entry"
			me := &managedExt{config: ExtConfig{Name: "member"}, host: h}
			share := newProcessShare(key, me, make(chan struct{}), discardAdmitWriter{io.Discard}, nil)
			if packed {
				me.packedCellKey = "cell"
				me.packedProcess = &packedProcessState{key: "cell", share: share}
			} else {
				me.share = share
			}
			withConn(me, conn)
			h.exts["member"] = me
			h.retention().register(share)
			_ = peerEnd.Close()
			h.handleIncoming(me, conn)
			if got, err := h.retention().claim(t.Context(), key); err != nil || got != nil {
				t.Fatalf("reload claimed a process with a closed connection: %v, %v", got, err)
			}
		})
	}
}

// A connection the host closes itself, as a reload closes a retiring generation, leaves the process claimable.
func TestHostClosedConnectionKeepsItsProcessRetained(t *testing.T) {
	t.Parallel()
	hostEnd, peerEnd := net.Pipe()
	defer func() { _ = peerEnd.Close() }()
	conn := NewConn("member", hostEnd)
	conn.Start(t.Context())
	h := NewHost(t.TempDir())
	key := "node-extension:member\x00entry"
	me := &managedExt{config: ExtConfig{Name: "member"}, host: h}
	share := newProcessShare(key, me, make(chan struct{}), discardAdmitWriter{io.Discard}, nil)
	me.share = share
	withConn(me, conn)
	h.retention().register(share)
	me.shuttingDown.Store(true)
	_ = conn.Close("retire")
	h.handleIncoming(me, conn)
	t.Cleanup(func() { share.release() })
	got, err := h.retention().claim(t.Context(), key)
	if err != nil || got != share {
		t.Fatalf("host retirement withdrew the process from retention: %v, %v", got, err)
	}
	got.release()
}
