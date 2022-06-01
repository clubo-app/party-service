package rpc

import (
	"context"
	"log"

	"github.com/clubo-app/packages/utils"
	"github.com/clubo-app/party-service/dto"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/cridenour/go-postgis"
)

func (s partyServer) CreateParty(c context.Context, req *pg.CreatePartyRequest) (*pg.Party, error) {
	start := req.StartDate.AsTime()
	end := req.EndDate.AsTime()

	d := dto.Party{
		Title:         req.Title,
		UserId:        req.RequesterId,
		Location:      postgis.Point{X: float64(req.Lat), Y: float64(req.Long)},
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
		log.Println(err)
		return nil, utils.HandleError(err)
	}

	return p.ToGRPCParty(), nil
}
