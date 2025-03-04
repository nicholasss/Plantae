-- +goose Up
create table biomes (
	id uuid primary key,
	koppen_class text not null unique,
	avg_summer_temp_c float,
	avg_winter_temp_c float,
	avg_summer_humid float,
	avg_winter_humid float,
	annual_rain_mm int,
	annual_sun_hours int
);

alter table plants
	add constraint fk_biome
	foreign key (biome_id)
	references biomes(id);

-- +goose Down
alter table plants
	drop constraint fk_biome;

drop table biomes;
