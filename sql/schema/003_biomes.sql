-- +goose Up
create table biomes (
  id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by text not null,
  updated_by text not null,
  deleted_by text,
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
	add constraint fk_biome
	foreign key (biome_id)
	references biomes(id);

-- +goose Down
alter table plant_species
	drop constraint fk_biome;

drop table biomes;
