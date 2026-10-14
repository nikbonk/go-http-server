-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetChirps :many
select *
from chirps
order by created_at asc;

-- name: GetChirpsByAuthor :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpById :one
select *
from chirps
where id = $1;

-- name: DeleteChirpByID :exec
delete from chirps
where id = $1;
