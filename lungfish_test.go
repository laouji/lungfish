package lungfish

import (
	"testing"
)

func TestCreateEvent(t *testing.T) {
	data := map[string]interface{}{
		"type": "message",
		"user": "aaaaaa",
		"text": "@botname: command_name arg1 arg2",
	}
	e := createEvent(data)
	if e.rawText != data["text"] {
		t.Error("For 'e.rawText' expected: ", data["text"], ", got: ", e.rawText)
	}

	if e.trigger.keyword != "command_name" {
		t.Error("For 'e.trigger.keyword' expected: ", "command_name", ", got: ", e.trigger.keyword)
	}
}
