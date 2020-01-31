package lungfish

import "testing"

func TestCreateEvent_Success(t *testing.T) {
	data := map[string]interface{}{
		"type": "message",
		"user": "aaaaaa",
		"text": "@botname: command_name arg1 arg2",
	}
	e, err := parseEvent(data)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}

	if e.rawText != data["text"] {
		t.Fatalf("expected: %s, got: %s", data["text"], e.rawText)
	}

	if e.Type != data["type"] {
		t.Fatalf("expected: %s, got: %s", data["type"], e.Type)
	}

	if e.trigger.keyword != "command_name" {
		t.Fatalf("expected: %s, got: %s", "command_name", e.trigger.keyword)
	}

	if len(e.trigger.args) != 2 {
		t.Fatalf("expected slice with: %d elements, got: %d", 2, len(e.trigger.args))
	}
}

func TestCreateEvent_Unsupported(t *testing.T) {
	data := map[string]interface{}{
		"type": "something_unsupported",
	}
	_, err := parseEvent(data)
	if err != ErrUnsupportedEventType {
		t.Fatalf("expected %s error, got: %s", ErrUnsupportedEventType, err)
	}
}

func TestCreateEvent_NotEnoughArgs(t *testing.T) {
	data := map[string]interface{}{
		"type": "message",
		"text": "@botname",
	}
	e, err := parseEvent(data)
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}

	// trigger should at least not be a nil pointer
	if len(e.trigger.args) != 0 {
		t.Fatalf("expected slice with: %d elements, got: %d", 0, len(e.trigger.args))
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
