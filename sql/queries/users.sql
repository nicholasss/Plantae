-- name: CreateUser :one
insert into users (
  id, created_at, updated_at,
  created_by, updated_by, join_date,
  is_admin, email, hashed_password
) values (
  $1, now(), now(),
  $1, $1, now(),
  false, $2, $3
) returning id, join_date, is_admin, email;

-- name: ResetUsersTable :exec
delete from users;

-- name: UpdateUserPasswordByID :exec
update users
set
  hashed_password = $2
where
  id = $1;

-- name: PromoteUserToAdminByID :exec
update users
set
  is_admin = true
where
  id = $1;

-- name: DemoteUserFromAdminByID :exec
update users
set
  is_admin = false
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
  where id = $1
  and deleted_at is null
  limit 1;

-- name: GetUserByEmailWithPassword :one
select * from users
  where email like $1
  and deleted_at is null
  limit 1;

-- name: GetUserByIDWithPassword :one
select * from users
  where id = $1
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
