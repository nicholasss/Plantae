-- +goose Up
create table water_needs (
  id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by uuid not null,
  updated_by uuid not null,
  deleted_by uuid,
  --
  -- type is either soil depth or interval watering
  -- values below fill in the 'X'
  -- 'X' number of days between waterings
  -- or
  -- 'X' millimeters of soil is dry between watering
  -- table data
  plant_type text not null,
  description text not null,
  dry_soil_mm integer,
  dry_soil_days integer
);

alter table plant_species
  add column water_needs_id uuid;

alter table plant_species
  add constraint fk_water_needs
  foreign key (water_needs_id)
  references water_needs(id);

-- +goose Down
alter table plant_species
  drop constraint fk_water_needs;

alter table plant_species
  drop column water_needs_id;

drop table water_needs;
