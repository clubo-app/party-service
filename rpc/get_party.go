package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) GetParty(c context.Context, req *pg.GetPartyRequest) (*pg.Party, error) {
	_, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	p, err := s.ps.Get(c, req.PartyId)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return p.ToGRPCParty(), nil
}
