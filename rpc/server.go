package rpc

import (
	"log"
	"net"
	"strings"

	"github.com/clubo-app/packages/stream"
	"github.com/clubo-app/party-service/service"
	pg "github.com/clubo-app/protobuf/party"
	"google.golang.org/grpc"
)

type partyServer struct {
	ps     service.PartyService
	stream stream.Stream
	pg.UnimplementedPartyServiceServer
}

func NewPartyServer(ps service.PartyService, stream stream.Stream) pg.PartyServiceServer {
	return &partyServer{
		ps:     ps,
		stream: stream,
	}
}

func Start(s pg.PartyServiceServer, port string) {
	var sb strings.Builder
	sb.WriteString("0.0.0.0:")
	sb.WriteString(port)
	conn, err := net.Listen("tcp", sb.String())
	if err != nil {
		log.Fatalln(err)
	}

	grpc := grpc.NewServer()

	pg.RegisterPartyServiceServer(grpc, s)

	log.Println("Starting gRPC Server at: ", sb.String())
	if err := grpc.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
