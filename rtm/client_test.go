package rtm_test

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/laouji/lungfish/rtm"
)

func TestClient_WebsocketHandshake(t *testing.T) {
	reqCount := 0
	handlerFn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqCount++
		key := req.Header.Get("Sec-WebSocket-Key")

		w.Header().Set("Upgrade", "websocket")
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Sec-WebSocket-Accept", websocketAcceptFromKey(key))
		w.WriteHeader(http.StatusSwitchingProtocols)
		w.Write([]byte{}) //http ResponseWriter closes connection
	})
	mockServer := httptest.NewServer(handlerFn)
	defer mockServer.Close()

	endpoint := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	client := rtm.NewClient(1)
	_, err := client.Start(endpoint)
	if err != nil {
		t.Errorf("expected nil error, got: %s", err)
	}

	if reqCount != 1 {
		t.Errorf("expected 1handshake request, got: %d", reqCount)
	}
}

func websocketAcceptFromKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	// magic string from https://tools.ietf.org/html/rfc6455
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
