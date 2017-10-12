package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	cs "github.com/billglover/cards/cards-service"
	"google.golang.org/grpc"
)

func createCard(client cs.CardProviderClient, c *cs.Card) {
	log.Printf("creating card: %s\n", c)
	card, err := client.Create(context.Background(), c)
	if err != nil {
		log.Printf("%v.GetFeatures(_) = _, %v: \n", client, err)
	}
	log.Printf("created card: %s\n", card)
}

func deleteCard(client cs.CardProviderClient, c *cs.Card) {
	log.Printf("deleting card: %s\n", c)
	_, err := client.Delete(context.Background(), c)
	if err != nil {
		log.Printf("%v.GetFeatures(_) = _, %v: \n", client, err)
	}
	log.Printf("deleted card: %s\n", c)
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

	createCard(client, &cs.Card{Id: 1})
	deleteCard(client, &cs.Card{Id: 1})
}
