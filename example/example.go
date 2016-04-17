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
		usersList := conn.GetUsersList()
		var memberName string
		for i := 0; i < len(usersList.Members); i++ {
			if usersList.Members[i].Id == e.UserId() {
				memberName = "@" + usersList.Members[i].Name
				break
			}
		}
		conn.PostMessage("o hai " + memberName)
	})

	conn.Loop()
}
