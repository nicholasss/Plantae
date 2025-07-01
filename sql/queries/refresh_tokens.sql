-- name: CreateRefreshToken :one
insert into refresh_tokens (
  refresh_token,
  created_at, updated_at,
  created_by, updated_by,
  expires_at, user_id
) values (
  $1,
  now(), now(),
  $2, $2,
  $3, $2
) returning
  refresh_token,
  created_at, updated_at,
  created_by, updated_by,
  expires_at, user_id;

-- name: GetUserFromRefreshToken :one
select * from refresh_tokens
where
  refresh_token = $1 and
  deleted_by is null and
  revoked_at is null
order by created_at desc
limit 1;

-- name: GetValidRefreshTokenFromUserID :one
select * from refresh_tokens
where
  user_id = $1 and
  deleted_by is null and
  revoked_at is null
order by created_at desc
limit 1;

-- name: RevokeRefreshTokenWithToken :one
update refresh_tokens
set
  updated_at = now(),
  updated_by = $2,
  revoked_at = now(),
  revoked_by = $2
where refresh_token = $1
returning user_id;
