# lungfish
==========

A simplistic event-based framework for creating slackbots in golang

## DESCRIPTION

lungfish makes it easy to create a chatbot that can respond to mentions using user-defined callbacks
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

	conn.Loop()
}
```

## REACTIONS

A reaction is a keyword and callback function pair that lets you define the way your bot responds to commands.

```go
conn.RegisterReaction("hello", func(e *lungfish.Event) {
    conn.PostMessage("o hai")
})
```

## INSTALLATION

```
go get github.com/laouji/lungfish
```
