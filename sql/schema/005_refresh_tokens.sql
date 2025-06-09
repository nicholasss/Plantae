-- +goose Up
create table refresh_tokens (
  refresh_token text primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by uuid not null,
  updated_by uuid not null,
  deleted_by uuid,
  --
  -- table data
  revoked_at timestamp,
  revoked_by text,
  expires_at timestamp not null,
  --
  -- table foreign key
  user_id uuid not null,
  constraint fk_user
  foreign key (user_id)
  references users(id)
  on delete cascade
);

-- +goose Down
drop table refresh_tokens;
