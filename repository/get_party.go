package repository

import (
	"context"

	"github.com/paulmach/orb/encoding/wkb"
)

const getParty = `-- name: GetParty :one
SELECT id, user_id, title, is_public, ST_AsBinary(location) as location, street_address, postal_code, state, country, start_date, end_date FROM parties
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetParty(ctx context.Context, id string) (Party, error) {
	row := q.db.QueryRow(ctx, getParty, id)
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
