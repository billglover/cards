package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	cs "github.com/billglover/cards/cards-service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	RootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a card",
	Long:  `Delete a card by providing its id.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println("card identifier should be a number")
			return
		}
		card := &cs.Card{Id: uint64(i)}
		deleteCard(card)
	},
}

func deleteCard(c *cs.Card) error {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := cs.NewCardProviderClient(conn)

	_, err = client.Delete(context.Background(), c)
	if err != nil {
		log.Printf("deleting card: {%+v} %s\n", c, "FAILED")
		log.Printf("unable to delete card: %v", grpc.ErrorDesc(err))
		return nil
	}
	log.Printf("deleting card: {%+v} %s\n", c, "SUCCESS")

	return nil
}
