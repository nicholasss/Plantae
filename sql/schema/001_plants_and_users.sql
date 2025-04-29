-- +goose Up
create table plant_species (
	id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by timestamp with time zone not null,
  updated_by timestamp with time zone null,
  deleted_by timestamp with time zone,
  --
  --
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
  created_by timestamp with time zone not null,
  updated_by timestamp with time zone null,
  deleted_by timestamp with time zone,
  --
  --
  join_date timestamp with time zone not null,
  is_admin boolean not null,
  email text not null,
  hashed_password text
);

-- +goose Down
drop table plant_species;

drop table users;
