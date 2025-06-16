-- +goose Up
create table watering_needs (
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
  watering_type text not null,
  description text not null,
  -- type is either soil depth or interval watering
  -- values below fill in the 'X'
  -- 'X' number of days between waterings
  -- or
  -- 'X' inches of soil is dry between watering
  grow_season float not null,
  transition_season float not null,
  dormant_season float not null
);

alter table plant_species
  add column watering_needs_id uuid;

alter table plant_species
  add constraint fk_watering_needs
  foreign key (watering_needs_id)
  references watering_needs(id);

-- +goose Down
alter table plant_species
  drop constraint fk_watering_needs;

alter table plant_species
  drop column watering_needs_id;

drop table watering_needs;
