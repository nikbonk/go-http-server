-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, expires_at, user_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetRefreshToken :one
select *
from refresh_tokens
where token like $1;

-- name: RevokeToken :one
update refresh_tokens
SET revoked_at = $1, updated_at = $2
where token like $3
returning *;
