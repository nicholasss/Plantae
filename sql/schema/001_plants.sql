-- +goose Up
create table plants (
	id uuid primary key,
	biome_id uuid,
	room_id uuid,
	water_schedule_id uuid,
	species_name text not null unique,
	pet_toxic bool,
	human_toxic bool,
	human_edible bool,
	avg_ideal_temp_c float not null,
	avg_ideal_humid float,
	ideal_light_hours float,
	ideal_light_intensity text
);

-- +goose Down
drop table plants;
