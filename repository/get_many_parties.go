package repository

import (
	"context"

	"github.com/paulmach/orb/encoding/wkb"
)

type GetManyPartiesParams struct {
	Ids   []string
	Limit int32
}

const getManyParties = `-- name: GetManyParties :many
SELECT id, user_id, title, is_public, ST_AsBinary(location), street_address, postal_code, state, country, start_date, end_date FROM parties
WHERE id IN($1::text[])
LIMIT $2
`

func (q *Queries) GetManyParties(ctx context.Context, arg GetManyPartiesParams) ([]Party, error) {
	rows, err := q.db.Query(ctx, getManyParties, arg.Ids, arg.Limit)
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
