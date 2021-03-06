package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	mysql "github.com/billglover/cards/cards-mysql"
	//neo "github.com/billglover/cards/cards-neo"
	//mgo "github.com/billglover/cards/cards-mgo"
	cs "github.com/billglover/cards/cards-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// csServer represents a cards service server. It holds references to the
// databases used to store cards and decks.
// TODO: use an interface here
type csServer struct {
	db *mysql.DB
}

// Create creates an instance of a card in the database. It returns the
// card that has been created or an error.
func (s *csServer) Create(ctx context.Context, c *cs.Card) (*cs.Card, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	uid, err := tx.CreateCard(c)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	c.Id = uid

	log.Printf("created card: %+v\n", c)
	return c, nil
}

// Delete removes an instance of a card from the database. It returns an
// empty response or an error.
func (s *csServer) Delete(ctx context.Context, c *cs.Card) (*cs.Empty, error) {

	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	rows, err := tx.DeleteCard(c)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	if rows != 1 {
		log.Printf("unable to delete card: unexpected rows impacted %d (want 1)\n", rows)
		tx.Rollback()
		return nil, grpc.Errorf(codes.NotFound, "card not found")
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	log.Printf("deleted card: %+v\n", c)
	return &cs.Empty{}, nil
}

// Embed embeds one card inside another. It returns a copy of the pairing
// if successfully created or an error if embedding was unsuccessful.
func (s *csServer) Embed(ctx context.Context, p *cs.Pairing) (*cs.Pairing, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.EmbedCard(p.Parent, p.Child)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	if rows != 1 {
		log.Printf("unable to embed cards: unexpected rows impacted %d (want 1)\n", rows)
		tx.Rollback()
		return nil, grpc.Errorf(codes.NotFound, "embedding failed")
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	log.Printf("embedded cards: %+v\n", p)

	return p, nil
}

// Remove removes one card from inside another. It returns an empty response
// if pairing is successfully removed or an error if removal was unsuccessful.
func (s *csServer) Remove(ctx context.Context, p *cs.Pairing) (*cs.Empty, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.RemoveCard(p.Parent, p.Child)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	if rows != 1 {
		log.Printf("unable to remove cards: unexpected rows impacted %d (want 1)\n", rows)
		tx.Rollback()
		return nil, grpc.Errorf(codes.NotFound, "remove failed")
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	log.Printf("remove cards: %+v\n", p)

	return &cs.Empty{}, nil
}

// Get returns a card by querying based on its id. It returns a card or an error.
func (s *csServer) Get(ctx context.Context, c *cs.Card) (*cs.Card, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	card, err := tx.GetCard(c)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	log.Printf("returned card: %+v\n", card)

	return card, nil
}

// newServer is a helper method that returns a new instance of the cards
// service server.
func newServer() *csServer {
	s := new(csServer)
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8040))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := mysql.Open("root@/CardsService?charset=utf8")
	//db, err := mgo.Open("mongodb://127.0.0.1:27017")
	//db, err := neo.Open("bolt://neo4j:password@localhost:7687")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := newServer()
	srv.db = db

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	cs.RegisterCardProviderServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
