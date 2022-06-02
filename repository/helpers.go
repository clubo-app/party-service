package repository

import (
	pg "github.com/clubo-app/protobuf/party"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p Party) ToGRPCParty() *pg.Party {
	id, err := ksuid.Parse(p.ID)
	if err != nil {
		return nil
	}

	return &pg.Party{
		Id:            p.ID,
		UserId:        p.UserID,
		Title:         p.Title,
		IsPublic:      p.IsPublic,
		Lat:           float32(p.Location.X()),
		Long:          float32(p.Location.Y()),
		StreetAddress: p.StreetAddress.String,
		PostalCode:    p.PostalCode.String,
		State:         p.State.String,
		Country:       p.Country.String,
		StartDate:     timestamppb.New(p.StartDate.Time),
		EndDate:       timestamppb.New(p.EndDate.Time),
		CreatedAt:     timestamppb.New(id.Time()),
	}
}
