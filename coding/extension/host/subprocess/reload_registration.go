package subprocess

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/MichaelKinsy/PiG/coding/extension"
)

func reloadStagedCellKey(item stagedManagedExt, cellKeys map[string]string) string {
	if key := cellKeys[item.name]; key != "" {
		return key
	}
	if key := item.me.packedCellKey; key != "" {
		return key
	}
	return "extension:" + item.name
}

func (h *Host) abortStaged(staged []stagedManagedExt, reason string) {
	for _, item := range staged {
		item.me.shuttingDown.Store(true)
		if conn := item.me.connection(); conn != nil {
			conn.fail(errors.New(reason))
		}
	}
	for _, item := range staged {
		h.stopManaged(item.me, reason)
	}
	for _, item := range staged {
		item.me.reap()
	}
}

func (h *Host) discardFailedReloadCells(staged []stagedManagedExt, cellKeys map[string]string, failed map[string]bool) []stagedManagedExt {
	var discarded []stagedManagedExt
	staged = slices.DeleteFunc(staged, func(item stagedManagedExt) bool {
		if !failed[reloadStagedCellKey(item, cellKeys)] {
			return false
		}
		discarded = append(discarded, item)
		return true
	})
	h.abortStaged(discarded, "reload cell rolled back")
	return staged
}

func (h *Host) loadExtensionRegistration(ctx context.Context, me *managedExt, reg *RegisterPayload) error {
	if deferred, _ := ctx.Value(deferLoadedProvidersKey{}).(bool); deferred {
		me.pendingRegistration = reg
		return nil
	}
	return h.registerExtensionAPI(me, reg)
}

func (h *Host) prepareReloadRegistrations(staged []stagedManagedExt, cellKeys map[string]string) ([]stagedManagedExt, []cellFailure) {
	var order []string
	groups := make(map[string][]stagedManagedExt)
	for _, item := range staged {
		key := reloadStagedCellKey(item, cellKeys)
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], item)
	}
	servers := h.providerRuntime.McpServers().List()
	var failures []cellFailure
	failed := make(map[string]bool)
	for _, key := range order {
		if !slices.ContainsFunc(groups[key], func(item stagedManagedExt) bool {
			return item.me.pendingRegistration != nil && len(item.me.pendingRegistration.McpServers) != 0
		}) {
			continue
		}
		candidate := extension.CreateExtensionRuntime()
		for _, server := range servers {
			candidate.McpServers().Register(server)
		}
		var registrationErr error
		for _, item := range groups[key] {
			if reg := item.me.pendingRegistration; reg != nil {
				for _, server := range reg.McpServers {
					if err := candidate.RegisterMcpServer(extConfigOrigin(item.me.config), server.Name, server.Config); err != nil {
						registrationErr = fmt.Errorf("reload cell %s: %w", key, err)
						break
					}
					prepared, _ := candidate.McpServers().Get(server.Name)
					item.me.pendingMcpServers = append(item.me.pendingMcpServers, prepared)
				}
			}
			if registrationErr != nil {
				break
			}
		}
		if registrationErr == nil {
			servers = candidate.McpServers().List()
			continue
		}
		for _, item := range groups[key] {
			failed[item.name] = true
			failures = append(failures, cellFailure{cfg: item.me.config, err: registrationErr})
		}
		h.abortStaged(groups[key], "reload registration rolled back")
	}
	return slices.DeleteFunc(staged, func(item stagedManagedExt) bool { return failed[item.name] }), failures
}

func (h *Host) publishReloadRegistrations(me *managedExt) {
	if reg := me.pendingRegistration; reg != nil {
		for _, server := range me.pendingMcpServers {
			h.providerRuntime.McpServers().Register(server)
			h.trackMcpServer(me, server.Name)
		}
		for _, ref := range reg.UnregisterVirtualModels {
			h.providerRuntime.UnregisterVirtualModel(ref.Provider, ref.ID)
			h.untrackVirtualModel(me, ref)
		}
		for _, model := range reg.VirtualModels {
			h.providerRuntime.QueueVirtualModelRegistration(h.virtualModelDefinition(me, model), extConfigOrigin(me.config))
			h.trackVirtualModel(me, VirtualModelRef{Provider: model.Provider, ID: model.ID})
		}
		if h.uiBridge != nil {
			h.uiBridge.RegisterExtConn(me.config.Name, me.connection(), reg.WantsModelRegistry)
		}
		me.pendingRegistration, me.pendingMcpServers = nil, nil
	}
	if len(me.pendingProviders) != 0 {
		h.publishLoadedProviders(me.parentCtx, me, me.connection(), me.pendingProviders)
		me.pendingProviders = nil
	}
}
