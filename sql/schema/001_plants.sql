-- +goose Up
create table plants (
	id uuid primary key,
	biome_id uuid not null,
	room_id uuid not null,
	water_schedule_id uuid not null,
	species_name text not null unique,
	pet_toxic bool not null,
	human_toxic bool not null,
	human_edible bool not null,
	avg_ideal_temp_c float not null,
	avg_ideal_humid float not null,
	ideal_light_hours float not null,
	ideal_light_intensity text not null
);

-- +goose Down
drop table plants;
