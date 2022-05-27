package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	pg "github.com/clubo-app/protobuf/party"
)

func (s partyServer) GetManyParties(ctx context.Context, req *pg.GetManyPartiesRequest) (*pg.GetManyPartiesResponse, error) {
	ps, err := s.ps.GetMany(ctx, req.Ids)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	parties := make([]*pg.Party, len(ps))

	for i, p := range ps {
		parties[i] = p.ToGRPCParty()
	}

	return &pg.GetManyPartiesResponse{Parties: parties}, nil
}
