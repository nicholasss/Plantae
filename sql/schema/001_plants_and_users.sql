-- +goose Up
create table plant_species (
	id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by uuid not null,
  updated_by uuid not null,
  deleted_by uuid,
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
  created_by uuid not null,
  updated_by uuid not null,
  deleted_by uuid,
  --
  -- table data
  join_date timestamp with time zone not null,
  is_admin boolean not null,
  email text not null,
  hashed_password text not null
);

-- +goose Down
drop table plant_species;

drop table users;
