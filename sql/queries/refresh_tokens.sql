-- name: CreateRefreshToken :one
insert into refresh_tokens (
  refresh_token,
  created_at, updated_at,
  created_by, updated_by,
  expires_at, user_id
) values (
  $1,
  now(), now(),
  $2, $3,
  $4, $5
) returning
  refresh_token,
  created_at, updated_at,
  created_by, updated_by,
  expires_at, user_id;

-- name: GetUserFromRefreshToken :one
select * from refresh_tokens
where refresh_token = $1;

-- name: RevokeRefreshTokenWithToken :exec
update refresh_tokens
set
  updated_at = now(),
  updated_by = $2,
  revoked_at = now(),
  revoked_by = $3
where refresh_token = $1;
