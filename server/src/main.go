package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
)

func main() {
	Do()
}

func Do() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	cutoffTime := time.Now().Add(-7 * 24 * time.Hour)
	iter := firestoreClient.Collection("tokens").Where("timestamp", ">", cutoffTime).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error getting Messaging client: %v\n", err)
		}

		registrationToken := doc.Ref.ID

		// See documentation on defining a message payload.
		message := &messaging.Message{
			Data: map[string]string{
				"hello": "there",
			},
			Token: registrationToken,
		}

		// Send a message to the device corresponding to the provided
		// registration token.
		response, err := messagingClient.Send(ctx, message)
		if err != nil {
			log.Fatalln(err)
		}
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response)
	}
}
