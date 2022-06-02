package repository

import (
	"context"

	"github.com/paulmach/orb/encoding/wkb"
)

const getPartiesInRadius = `-- name: GetPartiesInRadius :many
SELECT id, user_id, title, is_public, ST_AsBinary(location), street_address, postal_code, state, country, start_date, end_date
FROM parties
WHERE ST_DWithin(
  location,
  ST_GeomFromWKB($1::text),
  $2::int
) LIMIT $3
`

type GetPartiesInRadiusParams struct {
	Bytes  string
	Radius int32
	Limit  int32
}

func (q *Queries) GetPartiesInRadius(ctx context.Context, arg GetPartiesInRadiusParams) ([]Party, error) {
	rows, err := q.db.Query(ctx, getPartiesInRadius, arg.Bytes, arg.Radius, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Party
	for rows.Next() {
		var i Party
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
