# lungfish

A simplistic event-based framework for creating slackbots in golang

## DESCRIPTION

lungfish makes it easy to create a chatbot that can respond to mentions using user-defined callbacks.

See `example/example.go` for a working implementation of the bot

## ESTABLISHING A CONNECTION

```go
package main

import (
	"github.com/laouji/lungfish"
)

func main() {
	token := "<YOUR SLACK API TOKEN>"
	channel := "#channel-name"

	conn := lungfish.NewConnection(token)

    // register the channel you want your bot to join
	conn.RegisterChannel(channel)

    // bot logic goes here

	err := conn.Run()
  if err != nil {
    // handle connection error
  }
}
```

## REACTIONS

A reaction is a keyword and callback function pair that lets you define the way your bot responds to commands.

```go
conn.RegisterReaction("hello", func(e *lungfish.Event) {
    conn.PostMessage("o hai")
})
```

<img width="254" alt="51e4de98" src="https://cloud.githubusercontent.com/assets/2435916/14772595/23a2ff8a-0adb-11e6-8428-3c2467ff9669.png">

## COMMANDS

Here are some built-in API calls that come in handy in callbacks:

#### func (conn *Connection) PostMessage(text string)

Post a message to the registered channel:

```go
conn.PostMessage("hello world")
```

#### func (conn *Connection) GetUserInfo(userId string) UsersInfoResponseData

Get information about the user from a user ID via https://api.slack.com/methods/users.info

```go
userInfo := conn.GetUserInfo(e.UserId())

if !userInfo.Ok {
    log.Println("error: " + userInfo.Error)
} else {
    log.Println("user name is @" + userInfo.User.Name)
}
```

#### func (conn *Connection) GetUsersList() UsersListResponseData

Get information about all users in the registered channel via https://api.slack.com/methods/users.list

```go
conn.RegisterChannel("#channel-name")
usersList := conn.GetUsersList()
```

## INSTALLATION

```
go get github.com/laouji/lungfish
```
