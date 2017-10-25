package cmd

import (
	"context"
	"fmt"
	"log"

	cs "github.com/billglover/cards/cards-service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	RootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one card from within another",
	Long:  `Remove one card from within another.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pID := args[0]
		cID := args[1]

		pairing := &cs.Pairing{
			Parent: &cs.Card{Id: pID},
			Child:  &cs.Card{Id: cID},
		}
		removeCard(pairing)
	},
}

func removeCard(p *cs.Pairing) error {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := cs.NewCardProviderClient(conn)

	_, err = client.Remove(context.Background(), p)
	if err != nil {
		log.Printf("removing card: {%+v} %s\n", p, "FAILED")
		log.Printf("removing to embed card: %v", grpc.ErrorDesc(err))
		return nil
	}
	log.Printf("removing card: {%+v} %s\n", p, "SUCCESS")

	return nil
}
