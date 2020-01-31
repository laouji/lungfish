package lungfish

import "testing"

func TestCreateEvent(t *testing.T) {
	data := map[string]interface{}{
		"type": "message",
		"user": "aaaaaa",
		"text": "@botname: command_name arg1 arg2",
	}
	e := createEvent(data)
	if e.rawText != data["text"] {
		t.Fatalf("expected: %s, got: %s", data["text"], e.rawText)
	}

	if e.EventType != data["type"] {
		t.Fatalf("expected: %s, got: %s", data["type"], e.EventType)
	}

	if e.trigger.keyword != "command_name" {
		t.Fatalf("expected: %s, got: %s", "command_name", e.trigger.keyword)
	}

	if len(e.trigger.args) != 2 {
		t.Fatalf("expected slice with: %d elements, got: %d", 2, len(e.trigger.args))
	}
}

func TestRegisterChannel(t *testing.T) {
	conn := NewConnection("dummytoken")
	conn.RegisterChannel("#general")
	if conn.slackChannel != "#general" {
		t.Fatalf("expected: #general, got: %s", conn.slackChannel)
	}
}

func TestRegisterReaction(t *testing.T) {
	trigger := "hello"
	conn := NewConnection("dummytoken")
	conn.RegisterReaction(trigger, func(e *Event) { return })

	if _, ok := conn.reactions["hello"]; !ok {
		t.Fatalf("expected key named %s to be set in c.reactions", trigger)
	}
}
