package main

import (
	"log"

	"github.com/clubo-app/packages/stream"
	"github.com/clubo-app/party-service/config"
	"github.com/clubo-app/party-service/repository"
	"github.com/clubo-app/party-service/rpc"
	"github.com/clubo-app/party-service/service"
	"github.com/nats-io/nats.go"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	opts := []nats.Option{nats.Name("Party Service")}
	nc, err := stream.Connect(c.NATS_CLUSTER, opts)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()
	stream := stream.New(nc)

	d, err := repository.NewPartyRepository(c.DB_USER, c.DB_PW, c.DB_NAME, c.DB_HOST, c.DB_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	s := service.NewPartyService(d)

	p := rpc.NewPartyServer(s, stream)
	rpc.Start(p, c.PORT)
}
