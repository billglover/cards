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
	RootCmd.AddCommand(embedCmd)
}

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embed one card within another",
	Long:  `Embed one card within another.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pID := args[0]
		cID := args[1]

		pairing := &cs.Pairing{
			Parent: &cs.Card{Id: pID},
			Child:  &cs.Card{Id: cID},
		}
		embedCard(pairing)
	},
}

func embedCard(p *cs.Pairing) error {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8040), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := cs.NewCardProviderClient(conn)

	_, err = client.Embed(context.Background(), p)
	if err != nil {
		log.Printf("embedding card: {%+v} %s\n", p, "FAILED")
		log.Printf("unable to embed card: %v", grpc.ErrorDesc(err))
		return nil
	}
	log.Printf("embedding card: {%+v} %s\n", p, "SUCCESS")

	return nil
}
