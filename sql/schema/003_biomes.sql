-- +goose Up
create table environment (
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
	koppen_class text not null unique,
	avg_summer_temp_c float,
	avg_winter_temp_c float,
	avg_summer_humid float,
	avg_winter_humid float,
	annual_rain_mm int,
	annual_sun_hours int
);

alter table plant_species
	add constraint fk_environment
	foreign key (environment_id)
	references environment(id);

-- +goose Down
alter table plant_species
	drop constraint fk_environment;

drop table environment;
