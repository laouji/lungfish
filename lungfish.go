package lungfish

import (
	"fmt"
	"strings"

	"github.com/laouji/lungfish/api"
	"github.com/laouji/lungfish/rtm"
)

// unbuffered for the time being
var eventsChanBufferSize = 1

type callbackMethod func(*Event)

type Connection struct {
	token        string
	userId       string
	userName     string
	slackChannel string

	apiClient *api.Client
	rtmClient *rtm.Client
	reactions map[string]callbackMethod
}

type Event struct {
	data      map[string]interface{}
	EventType string
	rawText   string
	userId    string
	trigger   *Trigger
}

type Trigger struct {
	keyword string
	args    []string
}

func NewConnection(token string) *Connection {
	return &Connection{
		apiClient: api.NewClient(token),
		rtmClient: rtm.NewClient(eventsChanBufferSize),
		reactions: map[string]callbackMethod{},
	}
}

func createEvent(data map[string]interface{}) *Event {
	e := &Event{
		data:      data,
		EventType: data["type"].(string),
	}

	if userId, ok := data["user"]; ok {
		e.userId = userId.(string)
	}

	if e.EventType == "message" {
		e.rawText = data["text"].(string)

		args := strings.Split(strings.TrimSpace(e.rawText), " ")
		if len(args) > 1 {
			e.trigger = createTrigger(args[1], args[2:])
		}
	}

	return e
}

func createTrigger(keyword string, args []string) *Trigger {
	return &Trigger{
		keyword: keyword,
		args:    args,
	}
}

func (conn *Connection) Run() error {
	resData, err := conn.apiClient.Start()
	if err != nil {
		return fmt.Errorf("failed to start connection: %s", err)
	}

	conn.userId = resData.Self.Id
	conn.userName = resData.Self.Name

	eventsChan, err := conn.rtmClient.Start(resData.Url)
	if err != nil {
		return fmt.Errorf("failed to start rtm connection on %s: %s", resData.Url, err)
	}

	conn.handleEvents(eventsChan)
	return nil
}

func (conn *Connection) handleEvents(eventsChan <-chan map[string]interface{}) {
	for data := range eventsChan {
		if data == nil {
			continue
		}
		e := createEvent(data)

		switch data["type"].(string) {
		case "message":
			var isMention = strings.HasPrefix(data["text"].(string), "<@"+conn.userId+">")
			if !isMention {
				// ignore if bot's name not mentioned for now
				continue
			}

			if callback, ok := conn.reactions[e.Trigger().Keyword()]; ok {
				callback(e)
			}
		case "presence_change":
			presenceType := data["presence"].(string)
			if callback, ok := conn.reactions[presenceType]; ok {
				callback(e)
			}
		}
	}
}

func (conn *Connection) PostMessage(text string) error {
	return conn.apiClient.PostMessage(conn.slackChannel, text)
}

func (conn *Connection) GetUserInfo(userId string) (resData api.UsersInfoResponseData, err error) {
	return conn.apiClient.GetUserInfo(userId)
}

func (conn *Connection) RegisterChannel(slackChannel string) {
	conn.slackChannel = slackChannel
}

func (conn *Connection) RegisterReaction(triggerWord string, callback callbackMethod) {
	conn.reactions[triggerWord] = callback
}

func (conn *Connection) OwnUserId() string {
	return conn.userId
}

func (e *Event) Text() string {
	return e.rawText
}

func (e *Event) Trigger() *Trigger {
	return e.trigger
}

func (e *Event) UserId() string {
	return e.userId
}

func (t *Trigger) Keyword() string {
	return t.keyword
}
