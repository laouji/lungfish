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

func (client *Client) Start() (resData RtmStartResponseData, err error) {
	res, err := http.PostForm(baseUrl+endpointRtmStart, url.Values{
		"token": {client.token},
	})
	if err != nil {
		return resData, fmt.Errorf("request to %s endpoint failed: %s", endpointRtmStart, err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return resData, fmt.Errorf("failed to decode rtm.start response body: %s", err)
	}

	return resData, nil
}

func (client *Client) GetUserInfo(userId string) (resData UsersInfoResponseData, err error) {
	res, err := http.PostForm(baseUrl+endpointUsersInfo, url.Values{
		"token":   {client.token},
		"user":    {userId},
		"as_user": {"true"},
	})
	if err != nil {
		return resData, fmt.Errorf("users.info request failed: %s", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return resData, fmt.Errorf("users.info response was undecodable: %s", err)
	}
	return resData, nil
}

func (client *Client) PostMessage(channel string, text string) error {
	_, err := http.PostForm(baseUrl+endpointChatPostMessage, url.Values{
		"token":   {client.token},
		"channel": {channel},
		"text":    {text},
		"as_user": {"true"},
	})
	if err != nil {
		return fmt.Errorf("chat.postMessage request failed to post to channel %s: %s", channel, err)
	}
	return nil
}
