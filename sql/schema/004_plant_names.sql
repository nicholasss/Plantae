-- +goose Up
create table plant_names (
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
  plant_id uuid not null,
  --
  -- table data
  lang_code text not null,
  common_name text not null
);

alter table plant_names
  add constraint fk_plants
  foreign key (plant_id)
  references plant_species(id)
  on delete cascade;

-- +goose Down
alter table plant_names
	drop constraint fk_plants;

drop table plant_names;


