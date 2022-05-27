package repository

import (
	"time"

	pg "github.com/clubo-app/protobuf/party"
	"github.com/segmentio/ksuid"
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
		Lat:           float32(p.Location.FlatCoords()[0]),
		Long:          float32(p.Location.FlatCoords()[1]),
		StreetAddress: p.StreetAddress.String,
		PostalCode:    p.PostalCode.String,
		State:         p.State.String,
		Country:       p.Country.String,
		StartDate:     p.StartDate.Time.UTC().Format(time.RFC3339),
		EndDate:       p.EndDate.Time.UTC().Format(time.RFC3339),
		CreatedAt:     id.Time().UTC().Format(time.RFC3339),
	}
}
