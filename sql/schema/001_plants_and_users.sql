-- +goose Up
create table plant_species (
	id uuid primary key,
  created_at timestampz not null,
  updated_at timestampz not null,
  deleted_at timestampz,
  --
  created_by timestampz not null,
  updated_by timestampz not null,
  deleted_by timestampz,
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
  created_at timestampz not null,
  updated_at timestampz not null,
  deleted_at timestampz,
  --
  created_by timestampz not null,
  updated_by timestampz not null,
  deleted_by timestampz,
  --
  --
  join_date timestampz not null,
  is_admin boolean not null,
  email text not null,
  hashed_password text
);

-- +goose Down
drop table plant_species;

drop table users;
