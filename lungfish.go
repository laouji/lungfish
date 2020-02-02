package lungfish

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/laouji/lungfish/api"
	"github.com/laouji/lungfish/rtm"
)

var (
	// unbuffered for the time being
	eventsChanBufferSize = 1

	ErrUnsupportedEventType = errors.New("unsupported event type")
)

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
	Type   string
	UserId string // user who initiated the event

	rawData   map[string]interface{}
	rawText   string
	isMention bool // if the bot was mentioned or not
	trigger   *Trigger
}

// Text returns the raw text string sent to the bot
func (e *Event) Text() string {
	return e.rawText
}

func (e *Event) Trigger() *Trigger {
	return e.trigger
}

// Trigger is a struct containing the keyword that triggered a reaction and any following arguments
type Trigger struct {
	keyword string
	args    []string
}

func (t *Trigger) Keyword() string {
	return t.keyword
}

func NewConnection(token string) *Connection {
	return &Connection{
		apiClient: api.NewClient(api.BaseUrl, token),
		rtmClient: rtm.NewClient(eventsChanBufferSize),
		reactions: map[string]callbackMethod{},
	}
}

// RegisterChannel sets the slack channel that bot wishes to connect to
func (conn *Connection) RegisterChannel(slackChannel string) {
	conn.slackChannel = slackChannel
}

// RegisterReaction sets a callback method that can be triggered by commands sent to the bot
func (conn *Connection) RegisterReaction(triggerWord string, callback callbackMethod) {
	conn.reactions[triggerWord] = callback
}

// Run is a blocking method that triggers the bot to open a connectionw with the RTM endpoint
// and listen in on events happening in the channel.
// This method should be called after all the bot's reaction callback methods have been set up.
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

// PostMessage posts a message to the channel as the bot user
func (conn *Connection) PostMessage(text string) error {
	return conn.apiClient.PostMessage(conn.slackChannel, text)
}

// GetUserInfo fetches user profile information of a slack user given a user id
func (conn *Connection) GetUserInfo(userId string) (resData api.UsersInfoResponseData, err error) {
	return conn.apiClient.GetUserInfo(userId)
}

// OwnUserId is a helper function that returns the bot's own user id
func (conn *Connection) OwnUserId() string {
	return conn.userId
}

func (conn *Connection) parseEvent(rawData map[string]interface{}) (event *Event, err error) {
	event = &Event{rawData: rawData}
	if eventType, ok := rawData["type"].(string); ok {
		event.Type = eventType
	}

	// TODO: whitelist supported event types
	if event.Type != "message" {
		return nil, ErrUnsupportedEventType
	}

	if userId, ok := rawData["user"].(string); ok {
		event.UserId = userId
	}

	if text, ok := rawData["text"].(string); ok {
		event.rawText = text

		// bot expects space delimited commands
		args := strings.Split(strings.TrimSpace(text), " ")
		event.isMention, event.trigger = conn.parseArgs(args)
	}
	return event, nil
}

func (conn *Connection) parseArgs(args []string) (isMention bool, trigger *Trigger) {
	if len(args) < 1 {
		return false, &Trigger{}
	}

	// first argument is expected to be the bot's name
	isMention = strings.HasPrefix(args[0], "<@"+conn.userId+">")

	// TODO: trigger a help message if bot is called incorrectly
	if len(args) < 2 {
		return isMention, &Trigger{}
	}

	return isMention, &Trigger{
		keyword: args[1],
		args:    args[2:],
	}
}

func (conn *Connection) handleEvents(eventsChan <-chan map[string]interface{}) {
	for rawEvent := range eventsChan {
		if rawEvent == nil {
			continue
		}
		event, err := conn.parseEvent(rawEvent)
		if err != nil {
			log.Printf("skipping unparseable event %v: %s", rawEvent, err)
			continue
		}

		switch event.Type {
		case "message":
			if !event.isMention {
				// ignore if bot's name not mentioned for now
				continue
			}

			if callback, ok := conn.reactions[event.Trigger().Keyword()]; ok {
				callback(event)
			}
			// TODO: at the moment non message events are unexpected
		case "presence_change":
			presenceType := rawEvent["presence"].(string)
			if callback, ok := conn.reactions[presenceType]; ok {
				callback(event)
			}
		}
	}
}
