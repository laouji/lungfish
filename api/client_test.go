package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/laouji/lungfish/api"
)

func TestRtmStart_Success(t *testing.T) {
	expectedUserId := "user_id"
	expectedFormValues := map[string]string{
		"token": "dummyToken",
	}

	mockServer := getMockServer(t, expectedFormValues, []byte(fmt.Sprintf(`{
	"ok": true,
	"self": {
		"id": "%s",
		"name": "botname"
	}
}`, expectedUserId)))
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, expectedFormValues["token"])
	resData, err := client.Start()
	if err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}

	if userId := resData.Self.Id; userId != expectedUserId {
		t.Errorf("expected userId to be %s but got %s", expectedUserId, userId)
	}
}

func TestGetUserInfo_Success(t *testing.T) {
	expectedUserName := "userName"
	expectedFormValues := map[string]string{
		"token":   "dummyToken",
		"as_user": "true",
		"user":    "someUserId",
	}

	mockServer := getMockServer(t, expectedFormValues, []byte(fmt.Sprintf(`{
	"ok": true,
	"user": {
		"id": "%s",
		"name": "%s"
	}
}`, expectedFormValues["user"], expectedUserName)))
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, expectedFormValues["token"])
	resData, err := client.GetUserInfo(expectedFormValues["user"])
	if err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}

	if userId := resData.User.Id; userId != expectedFormValues["user"] {
		t.Errorf("expected userId to be %s but got %s", expectedFormValues["user"], userId)
	}

	if userName := resData.User.Name; userName != expectedUserName {
		t.Errorf("expected userName to be %s but got %s", expectedUserName, userName)
	}
}

func TestPostMessage_Success(t *testing.T) {
	expectedFormValues := map[string]string{
		"token":   "dummyToken",
		"as_user": "true",
		"channel": "#channel",
		"text":    "message contents",
	}

	mockServer := getMockServer(t, expectedFormValues, []byte(`{"ok": true}`))
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, expectedFormValues["token"])
	if err := client.PostMessage(expectedFormValues["channel"], expectedFormValues["text"]); err != nil {
		t.Errorf("expected error to be nil but got %s", err)
	}
}

func getMockServer(t *testing.T, expectedFormValues map[string]string, responseBody []byte) *httptest.Server {
	handlerFn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("expected error to be nil but got %s", err)
		}

		for key, expectedValue := range expectedFormValues {
			if actualValue := req.Form.Get(key); actualValue != expectedValue {
				t.Errorf("expected form value for %s to be %s but got %s", key, expectedValue, actualValue)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBody)
	})
	return httptest.NewServer(handlerFn)
}
