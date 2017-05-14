package lungfish

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type callbackMethod func(*Event)

type Connection struct {
	token     string
	userId    string
	userName  string
	channel   string
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

type SlackUserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RealName  string `json:"real_name"`
	Email     string `json:"email"`
	Image24   string `json:"image_24"`
}

type SlackUser struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	IsAdmin bool             `json:"is_admin"`
	IsOwner bool             `json:"is_owner"`
	Profile SlackUserProfile `json:"profile"`
}

type UsersListResponseData struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error"`
	Users []SlackUser `json:"members"`
}

type UsersInfoResponseData struct {
	Ok    bool      `json:"ok"`
	Error string    `json:"error"`
	User  SlackUser `json:"user"`
}

type RtmStartResponseData struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Url   string `json:"url"`
	Self  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
}

func NewConnection(token string) *Connection {
	return &Connection{
		token:     token,
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

func (conn *Connection) Start() (*websocket.Conn, error) {
	res, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {conn.token}})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var resData RtmStartResponseData
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Fatal(err)
	}

	conn.userId = resData.Self.Id
	conn.userName = resData.Self.Name

	return websocket.Dial(resData.Url, "", "https://slack.com")
}

func (conn *Connection) Run() error {
	ws, err := conn.Start()
	if err != nil {
		return err
	}

	conn.handleEvents(conn.receive(ws))
	return nil
}

func (conn *Connection) receive(ws *websocket.Conn) <-chan map[string]interface{} {
	ch := make(chan map[string]interface{})
	go func() {
		for {
			var data map[string]interface{}
			websocket.JSON.Receive(ws, &data)
			fmt.Printf("%+v\n", data)
			ch <- data
		}
	}()

	return ch
}

func (conn *Connection) handleEvents(ch <-chan map[string]interface{}) {
	for {
		data := <-ch
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

func (conn *Connection) PostMessage(text string) {
	res, err := http.PostForm("https://slack.com/api/chat.postMessage", url.Values{
		"token":   {conn.token},
		"channel": {conn.channel},
		"text":    {text},
		"as_user": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
}

func (conn *Connection) GetUsersList() UsersListResponseData {
	res, err := http.PostForm("https://slack.com/api/users.list", url.Values{
		"token":   {conn.token},
		"channel": {conn.channel},
		"as_user": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}

	var resData UsersListResponseData
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return resData
}

func (conn *Connection) GetUserInfo(userId string) UsersInfoResponseData {
	res, err := http.PostForm("https://slack.com/api/users.info", url.Values{
		"token":   {conn.token},
		"user":    {userId},
		"as_user": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}

	var resData UsersInfoResponseData

	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return resData
}

func (conn *Connection) RegisterChannel(channel string) {
	conn.channel = channel
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
