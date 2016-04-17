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
	eventType string
	rawText   string
	userId    string
	trigger   *Trigger
}

type Trigger struct {
	keyword string
	args    []string
}

type Member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UsersListResponseData struct {
	Ok      bool     `json:"ok"`
	Error   string   `json:"error"`
	Members []Member `json:"members"`
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
		eventType: data["type"].(string),
	}

	if e.eventType == "message" {
		e.rawText = data["text"].(string)
		e.userId = data["user"].(string)

		args := strings.Split(strings.TrimSpace(e.rawText), " ")
		log.Printf(args[1])
		e.trigger = createTrigger(args[1], args[2:])
	}

	return e
}

func createTrigger(keyword string, args []string) *Trigger {
	return &Trigger{
		keyword: keyword,
		args:    args,
	}
}

func (con *Connection) Loop() {
	ws, err := con.Start()
	if err != nil {
		log.Fatal(err)
	}

	for {
		var data map[string]interface{}
		websocket.JSON.Receive(ws, &data)
		fmt.Printf("%v\n", data)
		fmt.Println(con.BotUserId)

		if eventType, ok := data["type"].(string); ok && eventType == "message" {
			var isMention = strings.HasPrefix(data["text"].(string), "<@"+con.userId+">")
			if !isMention {
				// ignore if bot's name not mentioned for now
				continue
			}

			e := createEvent(data)
			if callback, ok := con.reactions[e.Trigger().Keyword()]; ok {
				callback(e)
			}
		}
	}
}

func (con *Connection) Start() (*websocket.Conn, error) {
	res, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {con.token}})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var resData RtmStartResponseData
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", resData)

	con.userId = resData.Self.Id
	con.userName = resData.Self.Name

	return websocket.Dial(resData.Url, "", "https://slack.com")
}

func (con *Connection) PostMessage(text string) {
	res, err := http.PostForm("https://slack.com/api/chat.postMessage", url.Values{
		"token":   {con.token},
		"channel": {con.channel},
		"text":    {text},
		"as_user": {"true"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
}

func (con *Connection) GetUsersList() UsersListResponseData {
	res, err := http.PostForm("https://slack.com/api/users.list", url.Values{
		"token":   {con.token},
		"channel": {con.channel},
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

func (con *Connection) RegisterChannel(channel string) {
	con.channel = channel
}

func (con *Connection) RegisterReaction(triggerWord string, callback callbackMethod) {
	con.reactions[triggerWord] = callback
}

func (con *Connection) BotUserId() string {
	return con.userId
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
