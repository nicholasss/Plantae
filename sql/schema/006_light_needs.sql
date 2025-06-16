-- +goose Up
create table light_needs (
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
  description text not null
);

alter table plant_species
  add column light_needs_id uuid;

alter table plant_species
  add constraint fk_light_needs
  foreign key (light_needs_id)
  references light_needs(id);

-- +goose Down
alter table plant_species
  drop constraint fk_light_needs;

alter table plant_species
  drop column light_needs_id;

drop table light_needs;
