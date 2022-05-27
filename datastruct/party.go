package datastruct

import (
	"time"

	pg "github.com/clubo-app/protobuf/party"
	"github.com/segmentio/ksuid"
)

type Party struct {
	Id            string    `json:"id"             db:"id"             validate:"required"`
	UserId        string    `json:"user_id"        db:"user_id"        validate:"required"`
	Title         string    `json:"title"          db:"title"          validate:"required"`
	IsPublic      bool      `json:"is_public"      db:"is_public"`
	Lat           float32   `json:"lat"            db:"lat"            validate:"required"`
	Long          float32   `json:"long"           db:"long"           validate:"required"`
	StreetAddress string    `json:"street_address" db:"street_address" validate:"required"`
	PostalCode    string    `json:"postal_code"    db:"postal_code"    validate:"required"`
	State         string    `json:"state"          db:"state"          validate:"required"`
	Country       string    `json:"country"        db:"country"        validate:"required"`
	StartDate     time.Time `json:"start_date"     db:"start_date"     validate:"required"`
	EndDate       time.Time `json:"end_date"       db:"end_date"       validate:"required"`
}

func (p Party) ToGRPCParty() *pg.Party {
	id, err := ksuid.Parse(p.Id)
	if err != nil {
		return &pg.Party{}
	}

	return &pg.Party{
		Id:            p.Id,
		UserId:        p.UserId,
		Title:         p.Title,
		IsPublic:      p.IsPublic,
		Lat:           p.Lat,
		Long:          p.Long,
		StreetAddress: p.StreetAddress,
		PostalCode:    p.PostalCode,
		State:         p.State,
		Country:       p.Country,
		StartDate:     p.StartDate.UTC().Format(time.RFC3339),
		EndDate:       p.EndDate.UTC().Format(time.RFC3339),
		CreatedAt:     id.Time().UTC().Format(time.RFC3339),
	}
}
