package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	cs "github.com/billglover/cards/cards-service"
	"google.golang.org/grpc"
)

type csServer struct{}

func (s *csServer) Create(ctx context.Context, c *cs.Card) (*cs.Card, error) {
	log.Printf("creating card: %s\n", c)
	c.Title = "Demo Card"
	log.Printf("created card: %s\n", c)
	return c, nil
}

func (s *csServer) Delete(ctx context.Context, c *cs.Card) (*cs.Empty, error) {
	log.Printf("deleting card: %s\n", c)
	log.Printf("deleted card: %s\n", c)
	return &cs.Empty{}, nil
}

func newServer() *csServer {
	s := new(csServer)
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8040))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	cs.RegisterCardProviderServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
