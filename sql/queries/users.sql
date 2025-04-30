-- name: CreateUser :one
insert into users (
  id, created_at, updated_at,
  created_by, updated_by,
  join_date,
  is_admin, email, hashed_password
) values (
  gen_random_uuid(), now(), now(),
  $1, $2, now(),
  $3, $4, $5
) returning *;

-- name: UpdateUserPasswordByID :exec
update users
set
  hashed_password = $2
where
  id = $1;

-- name: GetUserByEmailWithoutPassword :one
select 
  id, created_at, updated_at,
  created_by, updated_by,
  is_admin, email
from users
  where email like $1
  and deleted_at is null
  limit 1;

-- name: GetUserByIDWithoutPassword :one
select 
  id, created_at, updated_at,
  created_by, updated_by,
  is_admin, email
from users
  where id like $1
  and deleted_at is null
  limit 1;

-- name: GetUserByEmailWithPassword :one
select * from users
  where email like $1
  and deleted_at is null
  limit 1;

-- name: GetUserByIDWithPassword :one
select * from users
  where id like $1
  and deleted_at is null
  limit 1;

-- name: GetAllUsersWithoutPasswordByUpdated :many
select
  id, created_at, updated_at,
  created_by, updated_by,
  is_admin, email
from users
  where deleted_at is null
  order by updated_at desc;

-- name: GetAllUsersWithoutPasswordByJoinDate :many
select
  id, created_at, updated_at,
  created_by, updated_by,
  is_admin, email
from users
  where deleted_at is null
  order by join_date asc;
