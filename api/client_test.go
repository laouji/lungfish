package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/laouji/lungfish/api"
)

func TestPostMessage_Success(t *testing.T) {
	expectedToken := "dummy_token"
	expectedChannel := "#channel"
	expectedMessage := "message"

	reqCount := 0
	handlerFn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqCount++
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("expected error to be nil but got %s", err)
		}

		if token := req.Form.Get("token"); token != expectedToken {
			t.Errorf("expected token to be %s but got %s", expectedToken, token)
		}

		if channel := req.Form.Get("channel"); channel != expectedChannel {
			t.Errorf("expected channel to be %s but got %s", expectedChannel, channel)
		}

		if message := req.Form.Get("text"); message != expectedMessage {
			t.Errorf("expected text to be %s but got %s", expectedMessage, message)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	})
	mockServer := httptest.NewServer(handlerFn)
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, expectedToken)
	if err := client.PostMessage(expectedChannel, expectedMessage); err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}

	if reqCount != 1 {
		t.Errorf("expected request count to be 1 but got %d", reqCount)
	}
}
