package extensionconformance

import (
	"context"
	"encoding/json"
	"slices"
	"testing"
	"time"
)

func TestNativeEventBusToolCallerCancellationEndsListener(t *testing.T) {
	for _, isolation := range []string{"isolated", "shared-ok"} {
		t.Run(isolation, func(t *testing.T) {
			rig := newBusRig(t,
				busSpec("go", isolation, "emitter"),
				busSpec("go", isolation, "listener", busListenerAt{"toolwait", "record+awaitcancel"}),
			)
			tool, ok := rig.exts["emitter"].Tools["emitwait"]
			if !ok {
				t.Fatal("emitter has no emitwait tool")
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() {
				_, err := tool.Definition.Execute(ctx, "tc-tool-cancel", json.RawMessage(`{}`), nil)
				done <- err
			}()
			waitFor(t, func() bool { return slices.Contains(rig.lines(), `listener|toolwait:waiting|{"from":"emitter"}`) })
			cancel()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("cancelled tool returned success")
				}
			case <-time.After(10 * time.Second):
				t.Fatal("tool did not return after cancellation")
			}
			waitFor(t, func() bool { return slices.Contains(rig.lines(), `listener|toolwait:cancelled|"context canceled"`) })
		})
	}
}

func TestNativeEventBusCallerCancellationEndsListener(t *testing.T) {
	for _, language := range busNativeLanguages {
		for _, isolation := range []string{"isolated", "shared-ok"} {
			t.Run(language+"/"+isolation, func(t *testing.T) {
				rig := newBusRig(t,
					busSpec("go", isolation, "emitter"),
					busSpec("node", isolation, "prefix"),
					busSpec(language, isolation, "listener"),
					busSpec("go", isolation, "after"),
				)
				rig.must("prefix", "sub", "cancel record")
				rig.must("listener", "sub", "cancel awaitcancel")
				rig.must("after", "sub", "cancel record")
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				done := make(chan error, 1)
				go func() { done <- rig.exts["emitter"].Commands["emit"].Handler(ctx, "cancel json 73") }()
				waitFor(t, func() bool { return slices.Contains(rig.lines(), "listener|cancel:waiting|73") })
				cancel()
				select {
				case err := <-done:
					if err == nil {
						t.Fatal("cancelled caller returned success")
					}
				case <-time.After(10 * time.Second):
					t.Fatal("caller did not return after cancellation")
				}
				waitFor(t, func() bool { return slices.Contains(rig.lines(), `listener|cancel:cancelled|"context canceled"`) })
				waitFor(t, func() bool {
					retained := rig.host.EventBusRetention()
					return retained.Snapshots == 0 && retained.Emissions == 0
				})
				rig.expect("prefix|cancel|73", "listener|cancel:waiting|73", `listener|cancel:cancelled|"context canceled"`)
			})
		}
	}
}
