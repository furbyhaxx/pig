package coding

import (
	"strings"
	"testing"

	"github.com/MichaelKinsy/PiG/coding/extension"
)

// An unknown deliverAs used to take the steering path on sendMessage and be
// ignored on an idle sendUserMessage, so a misspelt value silently changed the
// delivery instead of failing (SPEC-0007 D3).
func TestSendMessageRejectsUnknownDeliverAs(t *testing.T) {
	h := newRecoveryHarness(t, harnessOptions{}, fauxReply("answer", "stop", 0))
	for _, value := range []extension.DeliverAs{"sideways", "FollowUp", "follow_up", "nextturn"} {
		err := h.session.SendMessage(extension.CustomMessageRef{CustomType: "notice", Content: "MISSING-DELIVERY", Display: true}, &extension.SendMessageOptions{DeliverAs: value})
		if err == nil {
			t.Fatalf("SendMessage with deliverAs %q = nil, want a rejection", value)
		}
		if !strings.Contains(err.Error(), string(value)) {
			t.Errorf("SendMessage error %q does not name the rejected value %q", err, value)
		}
		if projectionContains(h, "MISSING-DELIVERY") {
			t.Fatalf("the rejected deliverAs %q still reached the session", value)
		}
	}
}

func TestSendUserMessageRejectsUnknownDeliverAs(t *testing.T) {
	h := newRecoveryHarness(t, harnessOptions{}, fauxReply("answer", "stop", 0))
	for _, value := range []extension.DeliverAs{"sideways", "FollowUp", "follow_up"} {
		err := h.session.SendUserMessage(t.Context(), "MISSING-DELIVERY", &extension.SendUserMessageOptions{DeliverAs: value})
		if err == nil {
			t.Fatalf("SendUserMessage with deliverAs %q = nil, want a rejection", value)
		}
		if !strings.Contains(err.Error(), string(value)) {
			t.Errorf("SendUserMessage error %q does not name the rejected value %q", err, value)
		}
		if projectionContains(h, "MISSING-DELIVERY") {
			t.Fatalf("the rejected deliverAs %q still reached the session", value)
		}
	}
}

// nextTurn is a valid sendMessage mode and not a sendUserMessage mode: the
// type is shared, so only the call rejects it.
func TestSendUserMessageRejectsNextTurn(t *testing.T) {
	h := newRecoveryHarness(t, harnessOptions{}, fauxReply("answer", "stop", 0))
	err := h.session.SendUserMessage(t.Context(), "MISSING-DELIVERY", &extension.SendUserMessageOptions{DeliverAs: extension.DeliverAsNextTurn})
	if err == nil {
		t.Fatal("SendUserMessage with deliverAs nextTurn = nil, want a rejection")
	}
	if !strings.Contains(err.Error(), "nextTurn") {
		t.Errorf("SendUserMessage error %q does not name nextTurn", err)
	}
	if projectionContains(h, "MISSING-DELIVERY") {
		t.Error("the rejected nextTurn user message still reached the session")
	}
}

func TestKnownDeliverAsValuesAreAccepted(t *testing.T) {
	for _, value := range []extension.DeliverAs{"", extension.DeliverAsSteer, extension.DeliverAsFollowUp, extension.DeliverAsNextTurn} {
		h := newRecoveryHarness(t, harnessOptions{}, fauxReply("answer", "stop", 0))
		if err := h.session.SendMessage(extension.CustomMessageRef{CustomType: "notice", Content: "note", Display: true}, &extension.SendMessageOptions{DeliverAs: value}); err != nil {
			t.Fatalf("SendMessage with deliverAs %q = %v", value, err)
		}
	}
	for _, value := range []extension.DeliverAs{"", extension.DeliverAsSteer, extension.DeliverAsFollowUp} {
		h := newRecoveryHarness(t, harnessOptions{}, fauxReply("answer", "stop", 0))
		if err := h.session.SendUserMessage(t.Context(), "hello", &extension.SendUserMessageOptions{DeliverAs: value}); err != nil {
			t.Fatalf("SendUserMessage with deliverAs %q = %v", value, err)
		}
	}
}
