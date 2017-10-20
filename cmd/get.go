package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	cs "github.com/billglover/cards/cards-service"
)

func init() {
	RootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a card",
	Long:  `Get a card by providing its id`,
	Run: func(cmd *cobra.Command, args []string) {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println("card identifier should be a number")
			return
		}
		card := &cs.Card{Id: uint64(i)}

		getCard(card)
	},
}

func getCard(c *cs.Card) error {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := cs.NewCardProviderClient(conn)

	card, err := client.Get(context.Background(), c)
	if err != nil {
		log.Printf("get card: {%+v} %s\n", card, "FAILED")
		log.Printf("unable to get card: %v", grpc.ErrorDesc(err))
		return nil
	}
	log.Printf("get card: {%+v} %s\n", card, "SUCCESS")
	fmt.Println(card.Display())

	return nil
}
