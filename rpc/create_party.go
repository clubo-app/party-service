package rpc

import (
	"context"
	"time"

	"github.com/clubo-app/packages/utils"
	"github.com/clubo-app/party-service/dto"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/twpayne/go-geom"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) CreateParty(c context.Context, req *pg.CreatePartyRequest) (*pg.Party, error) {
	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid start date")
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid end date")
	}

	d := dto.Party{
		Title:         req.Title,
		UserId:        req.RequesterId,
		Location:      *geom.NewPointFlat(geom.XY, []float64{float64(req.Lat), float64(req.Long)}),
		IsPublic:      req.IsPublic,
		StreetAddress: req.StreetAddress,
		PostalCode:    req.PostalCode,
		State:         req.State,
		Country:       req.Country,
		StartDate:     start,
		EndDate:       end,
	}

	p, err := s.ps.Create(c, d)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return p.ToGRPCParty(), nil
}
