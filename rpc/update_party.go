package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	"github.com/clubo-app/party-service/dto"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/paulmach/orb"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) UpdateParty(c context.Context, req *pg.UpdatePartyRequest) (*pg.Party, error) {
	start := req.StartDate.AsTime()
	end := req.EndDate.AsTime()

	id, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	p := dto.Party{
		ID:            id.String(),
		UserId:        req.RequesterId,
		Title:         req.Title,
		Location:      orb.Point{float64(req.Long), float64(req.Lat)},
		StreetAddress: req.StreetAddress,
		PostalCode:    req.PostalCode,
		State:         req.State,
		Country:       req.Country,
		StartDate:     start,
		EndDate:       end,
	}

	d, err := s.ps.Update(c, p)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return d.ToGRPCParty(), nil
}
