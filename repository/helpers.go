package repository

import (
	"fmt"
	"net/url"

	pg "github.com/clubo-app/protobuf/party"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	g "github.com/golang-migrate/migrate/v4/source/github"
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

const version = 1

func validateSchema(url url.URL) error {
	url.Scheme = "pgx"
	url2 := fmt.Sprintf("%v%v", url.String(), "?sslmode=disable")
	g := g.Github{}
	d, err := g.Open("github://clubo-app/party-service/repository/migrations")
	if err != nil {
		return err
	}
	defer d.Close()

	m, err := migrate.NewWithSourceInstance("github", d, url2)

	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	defer m.Close()
	return nil
}
