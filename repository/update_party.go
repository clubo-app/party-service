package repository

import (
	"context"
	"database/sql"

	"github.com/paulmach/orb/encoding/wkb"
)

const updateParty = `-- name: UpdateParty :one
UPDATE parties SET
    title = $1,
    street_address = $2,
    postal_code = $3,
    state = $4,
    country = $5,
    start_date = $6,
    end_date = $7
WHERE id = $8
RETURNING id, user_id, title, is_public, location, street_address, postal_code, state, country, start_date, end_date
`

type UpdatePartyParams struct {
	Title         sql.NullString
	StreetAddress sql.NullString
	PostalCode    sql.NullString
	State         sql.NullString
	Country       sql.NullString
	StartDate     sql.NullTime
	EndDate       sql.NullTime
	ID            string
}

func (q *Queries) UpdateParty(ctx context.Context, arg UpdatePartyParams) (Party, error) {
	row := q.db.QueryRow(ctx, updateParty,
		arg.Title,
		arg.StreetAddress,
		arg.PostalCode,
		arg.State,
		arg.Country,
		arg.StartDate,
		arg.EndDate,
		arg.ID,
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
