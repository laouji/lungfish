package lungfish

import "testing"

func TestCreateEvent(t *testing.T) {
	data := make(map[string]interface{})
	data["text"] = "@lungfish blahblah extra args"
	data["type"] = "message"
	data["user"] = "me"

	e := createEvent(data)

	if e.eventType != data["type"] {
		t.Fatalf("expected: %s, got: %s", data["type"], e.eventType)
	}

	if e.trigger.keyword != "blahblah" {
		t.Fatalf("expected: %s, got: %s", "blahblah", e.trigger.keyword)
	}

	if len(e.trigger.args) != 2 {
		t.Fatalf("expected slice with: %d elements, got: %d", 2, len(e.trigger.args))
	}
}

func TestRegisterReaction(t *testing.T) {
	c := &Connection{
		token:     "token",
		userId:    "abcde",
		userName:  "Florence",
		channel:   "#fatm",
		reactions: map[string]callbackMethod{},
	}

	trigger := "ping"
	callback := func(e *Event) {}
	c.RegisterReaction(trigger, callback)

	if _, ok := c.reactions[trigger]; !ok {
		t.Fatalf("expected key named %s to be set in c.reactions", trigger)
	}
}
