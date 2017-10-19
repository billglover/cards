package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	cs "github.com/billglover/cards/cards-service"
)

func init() {
	RootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new card",
	Long:  `Create a new card by providing its title`,
	Run: func(cmd *cobra.Command, args []string) {
		createCard(args[0])
	},
}

func createCard(t string) error {

	c := &cs.Card{Title: t}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := cs.NewCardProviderClient(conn)

	card, err := client.Create(context.Background(), c)
	if err != nil {
		log.Printf("creating card: {%+v} %s\n", card, "FAILED")
		log.Printf("unable to create card: %v", grpc.ErrorDesc(err))
		return nil
	}
	log.Printf("creating card: {%+v} %s\n", card, "SUCCESS")

	return nil
}
