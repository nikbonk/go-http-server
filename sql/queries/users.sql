-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: ResetUsers :exec
delete from users;

-- name: UpdateUsers :one
update users
SET email = $1, hashed_password = $2, updated_at = $3
where id = $4
returning *;
