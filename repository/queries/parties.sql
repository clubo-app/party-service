-- name: DeleteParty :exec
DELETE FROM parties
WHERE id = $1 AND user_id = $2
RETURNING *;
