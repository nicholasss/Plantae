-- name: CreateUser :one
insert into users (
  id, created_at, updated_at,
  created_by, updated_by,
  is_admin, email, hashed_password
) values (
  gen_random_uuid(), now(), now(),
  $1, $2, $3, $4, $5
) returning *;

-- name: GetUserByEmailSafe :one
select 
  id, created_at, updated_at, deleted_at,
  created_by, updated_by, deleted_by,
  is_admin, email
from users
  where email like $1
  limit 1;

-- name: GetUsersByIDAll :one
select * from users
  where id like $1
  limit 1;

