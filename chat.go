package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	userName := os.Args[1]
	client, _ := ably.NewRealtime(
		ably.WithKey(os.Getenv("ABLY_KEY")),
		ably.WithClientID(userName),
		ably.WithEchoMessages(false),
	)
	log.Println(client.Auth.ClientID() + " joined")

	// then we define the topic (the channel)
	chatTopic := client.Channels.Get("chat")

	// define a reader to read from the stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error_reading_the_msg.", err.Error())
			continue
		}

		msg = strings.ReplaceAll(msg, "\n", "")

		if err = chatTopic.Publish(context.Background(), "message", msg); err != nil {
			log.Println("error_publishing_msg.", err.Error())
			continue
		}
	}
}
