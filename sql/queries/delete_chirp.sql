-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE chirps.id = $1
AND chirps.user_id = $2;