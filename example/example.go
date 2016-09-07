package main

import (
	"flag"
	"github.com/laouji/lungfish"
	"log"
)

var (
	token   = flag.String("t", "", "token: your slackbot's RTI API token")
	channel = flag.String("c", "#general", "channel: the channel you want the bot to reply in")
)

func main() {
	flag.Parse()
	if len(*token) == 0 {
		log.Fatal("token required")
	}
	conn := lungfish.NewConnection(*token)

	conn.RegisterChannel(*channel)
	conn.RegisterReaction("hello", func(e *lungfish.Event) {
		userInfo := conn.GetUserInfo(e.UserId())
		if !userInfo.Ok {
			log.Println(e.EventType + ": " + userInfo.Error)
			conn.PostMessage("error: " + userInfo.Error)
		} else {
			conn.PostMessage("o hai @" + userInfo.User.Name)
		}
	})

	conn.Run()
}
