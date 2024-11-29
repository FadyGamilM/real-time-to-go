package main

import (
	"log"
	"os"

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
}
