package rpc

import (
	"context"
	"time"

	"github.com/clubo-app/packages/utils"
	"github.com/clubo-app/party-service/dto"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/cridenour/go-postgis"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) UpdateParty(c context.Context, req *pg.UpdatePartyRequest) (*pg.Party, error) {
	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid start date")
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid end date")
	}

	id, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	p := dto.Party{
		ID:            id.String(),
		UserId:        req.RequesterId,
		Title:         req.Title,
		Location:      postgis.Point{X: float64(req.Lat), Y: float64(req.Long)},
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
