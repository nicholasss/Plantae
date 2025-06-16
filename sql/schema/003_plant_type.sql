-- +goose Up
create table plant_types (
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
	name text not null unique,
  description text not null,
  -- environment description
	max_temperature_celsius float,
	min_temperature_celsius float,
	max_humidity_percent float,
	min_humidity_percent float,
  -- soil description
  soil_organic_mix text,
  soil_grit_mix text,
  soil_drainage_mix text
);

alter table plant_species
  add column plant_type_id uuid;

alter table plant_species
	add constraint fk_plant_types
	foreign key (plant_type_id)
	references plant_types(id);

-- +goose Down
alter table plant_species
	drop constraint fk_plant_types;

alter table plant_species
  drop column plant_type_id;

drop table plant_types;
