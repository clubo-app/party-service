package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	"github.com/clubo-app/party-service/repository"
	"github.com/clubo-app/protobuf/party"
	pg "github.com/clubo-app/protobuf/party"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) GeoSearch(ctx context.Context, req *party.GeoSearchRequest) (*party.PagedParties, error) {
	if req.Radius > 50000 {
		return nil, status.Error(codes.InvalidArgument, "Radius too high, max is 50000")
	}

	ps, err := s.ps.GeoSearch(ctx, repository.GeoSearchParams{
		Lat:      req.Lat,
		Long:     req.Long,
		Radius:   req.Radius,
		IsPublic: req.IsPublic,
		Offset:   req.Offset,
		Limit:    req.Limit,
	})
	if err != nil {
		return nil, utils.HandleError(err)
	}

	var pp []*pg.Party
	for _, p := range ps {
		pp = append(pp, p.ToGRPCParty())
	}

	return &pg.PagedParties{Parties: pp}, nil
}
