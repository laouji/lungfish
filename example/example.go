package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/laouji/lungfish"
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
		userInfo, err := conn.GetUserInfo(e.UserId)
		if err != nil {
			log.Fatalf("error fetching user info for user id %s", e.UserId)
		}

		if !userInfo.Ok {
			log.Fatalf("failed to get user info %s", userInfo.Error)
		} else {
			message := fmt.Sprintf("o hai <@%s> (%s)", userInfo.User.Id, userInfo.User.Profile.RealName)
			if err = conn.PostMessage(message); err != nil {
				log.Fatalf("failed to post message %s", err)
			}
		}
	})

	err := conn.Run()
	if err != nil {
		log.Fatal(err)
	}
}
