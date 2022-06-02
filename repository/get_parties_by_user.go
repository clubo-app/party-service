package repository

import (
	"context"

	"github.com/paulmach/orb/encoding/wkb"
)

const getPartiesByUser = `-- name: GetPartiesByUser :many
SELECT id, user_id, title, is_public, ST_AsBinary(location), street_address, postal_code, state, country, start_date, end_date FROM parties
WHERE user_id = $1
ORDER BY id desc
LIMIT $2
OFFSET $3
`

type GetPartiesByUserParams struct {
	UserID string
	Limit  int32
	Offset int32
}

func (q *Queries) GetPartiesByUser(ctx context.Context, arg GetPartiesByUserParams) ([]Party, error) {
	rows, err := q.db.Query(ctx, getPartiesByUser, arg.UserID, arg.Limit, arg.Offset)
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
