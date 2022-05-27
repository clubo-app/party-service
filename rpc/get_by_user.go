package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	pg "github.com/clubo-app/protobuf/party"
)

func (s partyServer) GetByUser(c context.Context, req *pg.GetByUserRequest) (*pg.PagedParties, error) {
	ps, err := s.ps.GetByUser(c, req.UserId, int32(req.Limit), int32(req.Offset))
	if err != nil {
		return nil, utils.HandleError(err)
	}

	var pp []*pg.Party
	for _, p := range ps {
		pp = append(pp, p.ToGRPCParty())
	}

	return &pg.PagedParties{Parties: pp}, nil
}
