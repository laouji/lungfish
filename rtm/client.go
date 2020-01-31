package rtm

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

const origin = "https://slack.com"

type Client struct {
	eventsChan chan map[string]interface{}
	conn       *websocket.Conn
}

func NewClient(bufferSize int) *Client {
	return &Client{
		eventsChan: make(chan map[string]interface{}, bufferSize),
	}
}

func (client *Client) Start(endpoint string) (eventsChan <-chan map[string]interface{}, err error) {
	client.conn, err = websocket.Dial(endpoint, "", origin)
	if err != nil {
		return client.eventsChan, fmt.Errorf("failed to dial %s: %s", endpoint, err)
	}

	go client.receive()
	return client.eventsChan, nil
}

func (client *Client) receive() {
	for {
		var data map[string]interface{}
		if err := websocket.JSON.Receive(client.conn, &data); err != nil {
			// TODO passing errors back to main goroutine would be a good idea
			log.Printf("failed to receive on websocket: %s", err)
			close(client.eventsChan)
			break
		}
		client.eventsChan <- data
	}
}
