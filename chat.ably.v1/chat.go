package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	// SUB
	chatTopic.SubscribeAll(
		context.Background(),
		func(m *ably.Message) {
			log.Printf("received_msg_from_%s : %s\n", m.ClientID, m.Data)
		},
	)

	// Check for presence
	// first we enter the set of presence
	chatTopic.Presence.Enter(context.Background(), "")
	// we listen for any leave of joining
	chatTopic.Presence.SubscribeAll(
		context.Background(),
		func(pm *ably.PresenceMessage) {
			if pm.Action == ably.PresenceActionEnter {
				log.Printf("%s_has_joined\n", pm.ClientID)
			} else if pm.Action == ably.PresenceActionLeave {
				log.Printf("%s_has_left\n", pm.ClientID)
			}
		},
	)
	// if we terminted the program we need to exit the presence set
	exitPresenceChan := make(chan os.Signal, 1)
	signal.Notify(exitPresenceChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// block until we hear a closing signal
		<-exitPresenceChan
		// exit
		chatTopic.Presence.Leave(context.Background(), "")
		os.Exit(0)
	}()

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
