package repository

import (
	"context"
	"database/sql"

	"github.com/paulmach/orb/encoding/wkb"
)

const createParty = `-- name: CreateParty :one
INSERT INTO parties (
  id, 
  user_id,
  title,
  is_public,
  location,
  street_address,
  postal_code,
  state,
  country,
  start_date,
  end_date
) VALUES (
  $1, $2, $3, $4, ST_GeomFromText($5), 
  $6, $7, $8, $9, $10, $11
)
RETURNING id, user_id, title, is_public, ST_AsBinary(location) as location, street_address, postal_code, state, country, start_date, end_date
`

type CreatePartyParams struct {
	ID            string
	UserID        string
	Title         string
	IsPublic      bool
	Point         string
	StreetAddress sql.NullString
	PostalCode    sql.NullString
	State         sql.NullString
	Country       sql.NullString
	StartDate     sql.NullTime
	EndDate       sql.NullTime
}

func (q *Queries) CreateParty(ctx context.Context, arg CreatePartyParams) (Party, error) {
	row := q.db.QueryRow(ctx, createParty,
		arg.ID,
		arg.UserID,
		arg.Title,
		arg.IsPublic,
		arg.Point,
		arg.StreetAddress,
		arg.PostalCode,
		arg.State,
		arg.Country,
		arg.StartDate,
		arg.EndDate,
	)
	var i Party
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.IsPublic,
		wkb.Scanner(&i.Location),
		&i.StreetAddress,
		&i.PostalCode,
		&i.State,
		&i.Country,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}
