package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	cs "github.com/billglover/cards/cards-service"
	"google.golang.org/grpc"
)

// createCard creates a card in the cards service
func createCard(client cs.CardProviderClient, c *cs.Card) {
	card, err := client.Create(context.Background(), c)
	if err != nil {
		log.Printf("creating card: {%+v} %s\n", card, "FAILED")
		log.Printf("unable to create card: %v", grpc.ErrorDesc(err))
	}
	log.Printf("creating card: {%+v} %s\n", card, "SUCCESS")
}

// deleteCard removes a card from the cards service
func deleteCard(client cs.CardProviderClient, c *cs.Card) {
	_, err := client.Delete(context.Background(), c)
	if err != nil {
		log.Printf("deleting card: {%+v} %s\n", c, "FAILED")
		log.Printf("unable to delete card: %v", grpc.ErrorDesc(err))
	}
	log.Printf("deleting card: {%+v} %s\n", c, "SUCCESS")
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := cs.NewCardProviderClient(conn)

	createCard(client, &cs.Card{Title: "demo card"})
	deleteCard(client, &cs.Card{Id: 1})
	createCard(client, &cs.Card{Title: "demo card"})
}
