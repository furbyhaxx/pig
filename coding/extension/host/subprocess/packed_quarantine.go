package subprocess

import (
	"fmt"
	"slices"
	"strings"
)

// armPackedProcessGenerations suppresses earlier crash reports while replacement cells stage.
func (h *Host) armPackedProcessGenerations() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for key := range h.packedProcesses {
		h.advancePackedCellGenerationLocked(key)
	}
}

func (h *Host) restoreServingPackedProcesses() {
	h.mu.Lock()
	processes := make(map[*packedProcessState]struct{})
	for _, me := range h.exts {
		if p := me.packedProcess; p != nil {
			h.packedCellActiveGeneration[p.key] = p.generation
			h.packedProcesses[p.key] = p
			processes[p] = struct{}{}
		}
	}
	h.mu.Unlock()
	for p := range processes {
		if p.stopping.Load() || h.shuttingDown.Load() || (p.parent != nil && p.parent.Err() != nil) {
			continue
		}
		select {
		case <-p.startWait():
			if !p.node {
				h.quarantinePackedCellGeneration(p.key, p.generation, "packed process exited during reload")
			}
		default:
		}
	}
}

// packedCellSupervisor returns the persistent circuit breaker for cellKey,
// creating it with the same DefaultSupervisorConfig() every isolated
// extension and every packed member already uses (CNC-002: no new
// crash-loop policy). Unlike a managedExt's own Supervisor, which a fresh
// Reload replaces along with the managedExt (so an isolated extension's
// crash history does not survive past its next reload), a cell key's
// Supervisor here is stored on the Host and outlives any one packed
// process, because "did this cell key keep crashing" is a question about
// repeated reloads, not about one running process.
func (h *Host) packedCellSupervisor(cellKey string) *Supervisor {
	h.packedCellSupervisorsMu.Lock()
	defer h.packedCellSupervisorsMu.Unlock()
	if h.packedCellSupervisors == nil {
		h.packedCellSupervisors = make(map[string]*Supervisor)
	}
	s, ok := h.packedCellSupervisors[cellKey]
	if !ok {
		s = NewSupervisor(DefaultSupervisorConfig())
		h.packedCellSupervisors[cellKey] = s
	}
	return s
}

// quarantinePackedCell disables a crashed packed cell and tears down every
// extension currently attached to that cell. This is the runtime-cell fission
// point: callers/planners can inspect QuarantinedCells and reload the affected
// extensions as isolated subprocesses instead of reusing the failed pack.
// Callers that never claim a packedCellGeneration (an explicit non-process
// crash report, e.g. a per-member socket close, or a test) call this
// directly, which always quarantines.
func (h *Host) quarantinePackedCell(cellKey, reason string) {
	h.quarantinePackedCellGeneration(cellKey, 0, reason)
}

// quarantinePackedCellGeneration ignores reports outside the currently admitted or restored generation.
func (h *Host) quarantinePackedCellGeneration(cellKey string, generation int, reason string) {
	if cellKey == "" {
		return
	}
	h.mu.Lock()
	if generation != 0 && h.packedCellActiveGeneration[cellKey] != generation {
		// A later Reload already claimed a fresh generation for this key
		// (releaseRetryableQuarantines/nextPackedCellGeneration) before this
		// report reached the lock: it is stale and must not quarantine the
		// replacement process just because it shares the old content-derived
		// key.
		h.mu.Unlock()
		return
	}
	delete(h.packedCellActiveGeneration, cellKey)
	if h.quarantinedCells == nil {
		h.quarantinedCells = make(map[string]string)
	}
	if h.quarantinedCells[cellKey] == "" {
		h.quarantinedCells[cellKey] = reason
	}
	var members []*managedExt
	logPaths := make(map[*managedExt]string)
	for name, me := range h.exts {
		if me.packedCellKey == cellKey {
			if h.crashHandler() != nil {
				logPaths[me] = me.retainStderrLog()
			}
			delete(h.exts, name)
			members = append(members, me)
		}
	}
	h.mu.Unlock()
	names := make([]string, len(members))
	for i, me := range members {
		names[i] = me.config.Name
	}
	slices.Sort(names)
	// The returned delay/error are for a caller that reschedules a restart
	// itself; packed cells have no such automatic restart (only an explicit
	// Reload repacks), so only the resulting IsDisabled() state matters here.
	_, _ = h.packedCellSupervisor(cellKey).RecordCrash()
	for _, me := range members {
		me.shuttingDown.Store(true)
		h.stopManaged(me, "packed cell quarantined")
		if report := h.crashHandler(); report != nil {
			report(me.config.Name, 0, true, withStderrLog(quarantineNotice(reason, names), logPaths[me]))
		}
	}
}

// quarantineNotice states plainly which extensions a packed-cell crash
// stopped and that /reload recreates the cell (fission splits a quarantined
// cell back into isolated extensions on the next reload; see PlanCells).
func quarantineNotice(reason string, stoppedNames []string) string {
	if len(stoppedNames) == 0 {
		return reason
	}
	return fmt.Sprintf("%s; %s stopped; /reload restarts them", reason, strings.Join(stoppedNames, ", "))
}

// releaseRetryableQuarantines clears quarantine for every packed cell key
// whose circuit breaker has not tripped, so an explicit Reload gives it
// another chance to repack (CNC-002; the doc comment on Supervisor.Reset
// already says "Used by /reload" - this is that call site, extended to the
// cell-key-scoped supervisor packed cells use). Pig's default is one shared
// process per packable language (Go, Rust, Python, and Node cells alike);
// quarantine from a live crash degrades that to isolated subprocesses, but
// only until the next explicit reload, not permanently, unless the same
// cell key crashes DefaultSupervisorConfig().MaxCrashes times within
// CrashWindow - the exact same threshold every isolated extension's own
// Supervisor already enforces, so a Node cell's crash-loop protection reads
// no differently from an isolated extension's. PlanCells reads
// QuarantinedCells() right after this runs, so a cleared key is eligible to
// repack in the same reload.
func (h *Host) releaseRetryableQuarantines() {
	h.mu.Lock()
	keys := make([]string, 0, len(h.quarantinedCells))
	for key := range h.quarantinedCells {
		keys = append(keys, key)
	}
	h.mu.Unlock()
	for _, key := range keys {
		if h.packedCellSupervisor(key).IsDisabled() {
			continue
		}
		h.mu.Lock()
		delete(h.quarantinedCells, key)
		h.mu.Unlock()
	}
}

// nextPackedCellGeneration reserves a unique ID and revokes the prior crash-report owner.
func (h *Host) nextPackedCellGeneration(key string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.advancePackedCellGenerationLocked(key)
}

func (h *Host) advancePackedCellGenerationLocked(key string) int {
	if h.packedCellGeneration == nil {
		h.packedCellGeneration = make(map[string]int)
	}
	if h.packedCellActiveGeneration == nil {
		h.packedCellActiveGeneration = make(map[string]int)
	}
	h.packedCellGeneration[key]++
	h.packedCellActiveGeneration[key] = h.packedCellGeneration[key]
	return h.packedCellGeneration[key]
}
