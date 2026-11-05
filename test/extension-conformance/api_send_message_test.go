package extensionconformance

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/MichaelKinsy/PiG/coding"
)

// TestConformance_ExtensionAPIStrictDeliverAs drives pi.sendMessage with a deliverAs value outside the contract in every SDK. The host, not the SDK,
// must reject it: an SDK that substituted a default would still produce the custom message, so the expected host error is a value no fallback produces.
// SPEC-0007 D3 keeps the one shared validator on the Session boundary.
func TestConformance_ExtensionAPIStrictDeliverAs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping conformance suite in short mode (builds subprocess fixtures)")
	}
	t.Parallel()

	for _, tc := range sdkHarnessCases() {
		t.Run(tc.name, func(t *testing.T) {
			h := tc.make(t)
			t.Cleanup(func() {
				if h.cleanup != nil {
					h.cleanup()
				}
				if h.host != nil {
					h.host.Shutdown("test done")
				}
			})
			if h.bridge == nil {
				t.Skip("no host bridge for " + tc.name)
			}
			services, err := coding.CreateAgentSessionServices(coding.CreateAgentSessionServicesOptions{CWD: t.TempDir(), AgentDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			session, err := coding.NewSession(services, coding.SessionOptions{NoSession: true, SkipBuiltinTools: true})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = session.Close() })
			coding.BindExtensionHostSessionActions(h.bridge, func() *coding.Session { return session })

			cmd, ok := findCommand(h.runner, "send_message_bogus_deliver_as")
			if !ok {
				t.Fatal("send_message_bogus_deliver_as command not registered")
			}
			err = cmd.Handler(context.Background(), "")
			if err != nil && !strings.Contains(err.Error(), `unsupported deliverAs "bogus"`) {
				t.Fatalf("bogus deliverAs returned %v; want the host's unsupported deliverAs error", err)
			}
			if err == nil && !strings.Contains(tc.name, "node") {
				t.Fatal("bogus deliverAs returned success")
			}
			// Pi's Node loader declares pi.sendMessage void (loader.js:305-308), so its rejection is logged by the host and cannot be returned. A later host call on the same ordered connection proves the first was processed; a fallback default would then have appended the message.
			*h.notify = nil
			runConformanceCommand(t, h, "active_tools_get")
			pollUntilConformance(t, "the ordered host call after the rejected send never answered", func() bool {
				for _, notification := range *h.notify {
					if strings.HasPrefix(notification, "active_tools:") {
						return true
					}
				}
				return false
			})
			for _, entry := range session.Entries() {
				if strings.Contains(string(entry.Raw()), "bogus-deliver-as") {
					t.Fatalf("bogus deliverAs appended a custom message: %s", entry.Raw())
				}
			}
		})
	}
}

// (packages/coding-agent/src/core/extensions/types.ts:1703-1706, agent-session.ts:2293-2328 sendCustomMessage): every SDK sends a custom message through the host,
// which hands it to a real coding.Session through the production binder every mode uses (coding.BindExtensionHostSessionActions). An idle Session whose
// message does not trigger a turn appends it to the log as a custom_message entry (agent-session.ts sendCustomMessage: `else { this.agent.state.messages.push(appMessage); this.sessionManager.appendCustomMessageEntry(...) }`).
// The fixture's customType and content ("notice", "no-turn", "default") are values no SDK fallback produces. Both fixtures avoid a model turn: one passes
// triggerTurn:false, the other no options, and an idle Session only starts a turn when triggerTurn is true.
func TestConformance_ExtensionAPISendMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping conformance suite in short mode (builds subprocess fixtures)")
	}
	t.Parallel()

	for _, tc := range sdkHarnessCases() {
		t.Run(tc.name, func(t *testing.T) {
			h := tc.make(t)
			t.Cleanup(func() {
				if h.cleanup != nil {
					h.cleanup()
				}
				if h.host != nil {
					h.host.Shutdown("test done")
				}
			})
			if h.bridge == nil {
				t.Skip("no host bridge for " + tc.name)
			}
			services, err := coding.CreateAgentSessionServices(coding.CreateAgentSessionServicesOptions{CWD: t.TempDir(), AgentDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			session, err := coding.NewSession(services, coding.SessionOptions{NoSession: true, SkipBuiltinTools: true})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = session.Close() })
			coding.BindExtensionHostSessionActions(h.bridge, func() *coding.Session { return session })

			customMessages := func() []string {
				var found []string
				for _, entry := range session.Entries() {
					if entry.Base().Type != "custom_message" {
						continue
					}
					var custom struct {
						CustomType string `json:"customType"`
						Content    any    `json:"content"`
						Display    bool   `json:"display"`
					}
					if json.Unmarshal(entry.Raw(), &custom) == nil && custom.Display {
						found = append(found, custom.CustomType+"="+toString(custom.Content))
					}
				}
				return found
			}
			if got := customMessages(); len(got) != 0 {
				t.Fatalf("the Session starts with custom messages %v", got)
			}
			runConformanceCommand(t, h, "send_message_no_turn")
			runConformanceCommand(t, h, "send_message_default")
			want := []string{"notice=no-turn", "notice=default"}
			pollUntilConformance(t, "the Session never received the messages", func() bool {
				return len(customMessages()) >= len(want)
			})
			if got := customMessages(); !slices.Equal(got, want) {
				t.Fatalf("custom messages = %v, want %v", got, want)
			}
		})
	}
}
