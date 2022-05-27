-- name: CreateParty :one
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
  $1, $2, $3, $4, ST_GeomFromWKB($5), $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateParty :one
UPDATE parties SET
    title = sqlc.narg('title'),
    street_address = sqlc.narg('street_address'),
    postal_code = sqlc.narg('postal_code'),
    state = sqlc.narg('state'),
    country = sqlc.narg('country'),
    start_date = sqlc.narg('start_date'),
    end_date = sqlc.narg('end_date')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteParty :one
DELETE FROM parties
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetParty :one
SELECT * FROM parties
WHERE id = $1 LIMIT 1;

-- name: GetManyParties :many
SELECT * FROM parties
WHERE id IN(sqlc.arg('ids')::text[])
LIMIT sqlc.arg('limit');

-- name: GetPartiesInRadius :many
SELECT *
FROM parties
WHERE ST_DWithin(
  location,
  ST_GeomFromWKB(sqlc.arg('bytes')::text),
  sqlc.arg('radius')::int
) LIMIT sqlc.arg('limit');

-- name: GetPartiesByUser :many
SELECT * FROM parties
WHERE user_id = $1
ORDER BY id desc
LIMIT $2
OFFSET $3;
