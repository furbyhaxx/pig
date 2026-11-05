package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/MichaelKinsy/PiG/ai"
)

// Upstream holds deliverAs "nextTurn" messages apart from the steering and
// follow-up queues and injects them beside the user message before the first
// model call (agent-session.ts:1226). Pig routed them into FollowUp, which
// continues the running turn instead, so a message meant as context for the
// next prompt reached the model a call too late: and on an idle agent sat in a
// queue with no goroutine to drain it.
func TestNextTurnMessagesArriveWithTheUserMessage(t *testing.T) {
	agent := newTestAgentForNextTurn(t)

	agent.QueueNextTurn(AgentMessage{Custom: map[string]any{
		"role": RoleCustom, "customType": "chain_state", "content": "step 2 of 3",
	}})
	if len(agent.pendingNextTurn) != 1 {
		t.Fatalf("QueueNextTurn deferred %d messages, want 1", len(agent.pendingNextTurn))
	}

	if _, err := agent.Send(context.Background(), "continue"); err != nil {
		t.Fatalf("Send: %v", err)
	}

	// The deferred message must sit immediately after the user message, so the
	// first model call of the turn already carries it.
	userAt := -1
	for i, m := range agent.messages {
		if m.User != nil {
			userAt = i
		}
	}
	if userAt < 0 {
		t.Fatal("no user message in history")
	}
	if userAt+1 >= len(agent.messages) {
		t.Fatal("the deferred message was never injected")
	}
	next := agent.messages[userAt+1]
	if next.Custom == nil || next.Custom["customType"] != "chain_state" {
		t.Errorf("message after the user message = %+v; want the deferred chain_state", next)
	}
	if len(agent.pendingNextTurn) != 0 {
		t.Error("the deferred message was injected but not cleared")
	}
}

// A deferred message must not be delivered as a follow-up continuation.
func TestNextTurnDoesNotEnterTheFollowUpQueue(t *testing.T) {
	agent := newTestAgentForNextTurn(t)
	agent.QueueNextTurn(AgentMessage{Custom: map[string]any{"role": RoleCustom, "customType": "x"}})

	if agent.HasQueuedMessages() {
		t.Error("a nextTurn message entered the steering or follow-up queue")
	}
}

func TestNextTurnLeafStampAndDiscard(t *testing.T) {
	a := newTestAgentForNextTurn(t)
	leaf := new("a")
	a.BindNextTurnSession("session", func() *string { return leaf })
	queue := func(name string) { a.QueueNextTurn(AgentMessage{Custom: map[string]any{"customType": name}}) }
	queue("a-1")
	queue("a-2")
	leaf = new("b")
	queue("b")
	discarded := a.DiscardPendingNextTurn("session", leaf)
	if len(discarded) != 2 || discarded[0].Custom["customType"] != "a-1" || discarded[1].Custom["customType"] != "a-2" {
		t.Fatalf("discarded = %+v, want a-1, a-2", discarded)
	}
	kept := a.TakePendingNextTurn()
	if len(kept) != 1 || kept[0].Custom["customType"] != "b" {
		t.Fatalf("kept = %+v, want b", kept)
	}
	leaf = nil
	queue("root")
	if got := a.DiscardPendingNextTurn("session", nil); len(got) != 0 {
		t.Fatalf("root discarded = %+v", got)
	}
	if got := a.RebindNextTurnSession("other-session", func() *string { return nil }); len(got) != 1 {
		t.Fatalf("replacement discarded = %+v, want root", got)
	}
	if got := a.DiscardPendingNextTurn("other-session", nil); len(got) != 0 {
		t.Fatalf("duplicate discard = %+v", got)
	}
}

func BenchmarkNextTurnLeafDiscard(b *testing.B) {
	for _, size := range []int{1, 1000} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			a := &Agent{}
			leaf := new("a")
			a.BindNextTurnSession("session", func() *string { return leaf })
			message := AgentMessage{Custom: map[string]any{"customType": "benchmark"}}
			b.ReportAllocs()
			for b.Loop() {
				for range size {
					a.QueueNextTurn(message)
				}
				a.DiscardPendingNextTurn("session", new("b"))
			}
		})
	}
}

func newTestAgentForNextTurn(t *testing.T) *Agent {
	t.Helper()
	provider := &recordingProvider{seqs: [][]ai.AssistantMessageEvent{textSeq("")}}
	return mustNewAgent(AgentOptions{Model: &ai.Model{ID: "model", Provider: provider}})
}
