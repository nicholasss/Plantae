-- name: CreateUser :one
insert into users (
  id, created_at, updated_at, deleted_at,
  created_by, updated_by, deleted_by,
  is_admin, email, hashed_password
) values (
  gen_random_uuid(), now(), now(), now(),
  $1, $2, $3, $4, $5, $6
) returning *;
