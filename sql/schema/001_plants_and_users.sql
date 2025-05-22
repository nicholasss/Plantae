-- +goose Up
create table plants (
	id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by text not null,
  updated_by text not null,
  deleted_by text,
  --
  -- foreign keys
  biome_id uuid,
  --
  -- table data
	species_name text not null unique,
	human_poison_toxic boolean,
	pet_poison_toxic boolean,
	human_edible boolean,
  pet_edible boolean
);

create table users (
  id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by text not null,
  updated_by text not null,
  deleted_by text,
  --
  -- table data
  join_date timestamp with time zone not null,
  is_admin boolean not null,
  email text not null,
  hashed_password text not null
);

-- +goose Down
drop table plants;

drop table users;
