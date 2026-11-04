package subprocess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MichaelKinsy/PiG/coding/extension"
)

type deferLoadedProvidersKey struct{}

type loadedProvider struct {
	declaration ProviderDecl
	config      extension.ProviderConfig
	oauth       *ProviderOAuthConfig
}

func (h *Host) loadProviderDeclarations(ctx context.Context, me *managedExt, conn *Conn, declarations []ProviderDecl, isRestart bool) error {
	providers := make([]loadedProvider, 0, len(declarations))
	for _, declaration := range declarations {
		provider := loadedProvider{declaration: declaration}
		if native := declaration.Native; native != nil {
			if native.ID == "" || native.Key == "" || native.Auth == nil {
				return errors.New("native provider requires id, callback key and auth declaration")
			}
		} else {
			if isRestart {
				continue
			}
			if err := json.Unmarshal(declaration.Config, &provider.config); err != nil {
				return fmt.Errorf("provider %s: decode config: %w", declaration.Name, err)
			}
			var err error
			provider.oauth, err = decodeProviderOAuth(declaration.Name, declaration.Config)
			if err != nil {
				return err
			}
		}
		providers = append(providers, provider)
	}
	if deferred, _ := ctx.Value(deferLoadedProvidersKey{}).(bool); deferred && !isRestart {
		me.pendingProviders = providers
		return nil
	}
	if isRestart {
		for _, provider := range providers {
			if err := h.registerNativeProvider(ctx, me, conn, provider.declaration.Native, false); err != nil {
				return err
			}
		}
		return nil
	}
	h.publishLoadedProviders(ctx, me, conn, providers)
	return nil
}

func (h *Host) publishLoadedProviders(ctx context.Context, me *managedExt, conn *Conn, providers []loadedProvider) {
	for _, provider := range providers {
		if native := provider.declaration.Native; native != nil {
			proxy := h.newNativeProviderProxy(me, conn, native)
			h.queueNativeProvider(proxy)
			h.publishNativeProvider(proxy)
			extension.CallInitiated(ctx)
			if h.uiBridge != nil {
				h.uiBridge.PublishModelCatalog()
			}
		} else {
			h.queuePreparedLoadedProvider(me, conn, provider.declaration, provider.config, provider.oauth)
		}
	}
}
