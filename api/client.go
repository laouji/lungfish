package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}

func (client *Client) postForm(
	endpoint string,
	payload url.Values,
	responseData interface{},
) error {
	res, err := http.PostForm(endpoint, payload)
	if err != nil {
		return fmt.Errorf("request to %s endpoint with payload %s failed: %s", endpoint, payload, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request to %s endpoint with payload %s resulted in status code %d", endpoint, payload, res.StatusCode)
	}

	// if the caller passed nil, we don't need to decode the response body
	if responseData == nil {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(responseData)
	if err != nil {
		return fmt.Errorf("failed to decode %s response body: %s", endpoint, err)
	}
	return nil
}

func (client *Client) Start() (responseData RtmStartResponseData, err error) {
	err = client.postForm(baseUrl+endpointRtmStart, url.Values{
		"token": {client.token},
	}, &responseData)
	if err != nil {
		return responseData, err
	}
	return responseData, nil
}

func (client *Client) GetUserInfo(userId string) (responseData UsersInfoResponseData, err error) {
	err = client.postForm(baseUrl+endpointUsersInfo, url.Values{
		"token":   {client.token},
		"user":    {userId},
		"as_user": {"true"},
	}, &responseData)
	if err != nil {
		return responseData, err
	}
	return responseData, nil
}

func (client *Client) PostMessage(channel string, text string) error {
	return client.postForm(baseUrl+endpointChatPostMessage, url.Values{
		"token":   {client.token},
		"channel": {channel},
		"text":    {text},
		"as_user": {"true"},
	}, nil)
}
