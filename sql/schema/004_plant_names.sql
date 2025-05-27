-- +goose Up
create table plant_names (
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
  plant_id uuid,
  --
  -- table data
  lang_code text,
  common_name text
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


