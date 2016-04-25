package lungfish

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	conn := NewConnection("dummytoken")

	if conn.token != "dummytoken" {
		t.Error("For 'conn.token' expected: dummytoken, got: ", conn.token)
	}
}

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

func TestRegisterChannel(t *testing.T) {
	conn := NewConnection("dummytoken")
	conn.RegisterChannel("#general")
	if conn.channel != "#general" {
		t.Error("For 'conn.channel' expected: #general, got: ", conn.channel)
	}
}

func TestRegisterReaction(t *testing.T) {
	conn := NewConnection("dummytoken")
	conn.RegisterReaction("hello", func(e *Event) { return })

	_, ok := conn.reactions["hello"]
	if !ok {
		t.Error("No callback found for 'conn.reactions[\"hello\"]")
	}
}
