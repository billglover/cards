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
	RootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one card from within another",
	Long:  `Remove one card from within another.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		pID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println("parent identifier should be a number")
			return
		}

		cID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println("child identifier should be a number")
			return
		}

		pairing := &cs.Pairing{
			Parent: &cs.Card{Id: uint64(pID)},
			Child:  &cs.Card{Id: uint64(cID)},
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
