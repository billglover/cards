package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	cs "github.com/billglover/cards/cards-service"
	"google.golang.org/grpc"
)

type csServer struct {
	mysql *cs.DB
}

func (s *csServer) Create(ctx context.Context, c *cs.Card) (*cs.Card, error) {
	log.Printf("creating card: %s\n", c)

	tx, err := s.mysql.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	uid, err := tx.CreateCard(c)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c.Id = uint64(uid)

	log.Printf("created card: %s\n", c)
	return c, nil
}

func (s *csServer) Delete(ctx context.Context, c *cs.Card) (*cs.Empty, error) {
	log.Printf("deleting card: %s\n", c)

	tx, err := s.mysql.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	_, err = tx.DeleteCard(c)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

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

	// Demo a DB connection
	db, err := cs.Open("root@/CardsService?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := newServer()
	srv.mysql = db

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	cs.RegisterCardProviderServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
