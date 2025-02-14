-- name: InsertNamePlant :one
insert into plants (
	id, species_name
) values (
	$1, $2
) returning *;

-- name: InsertAllPlant :one
insert into plants (
	id, biome_id, room_id, water_schedule_id,
	species_name, pet_toxic, human_toxic, human_edible,
	avg_ideal_temp_c, avg_ideal_humid,
	ideal_light_hours, ideal_light_intensity
) values (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) returning *;

-- name: GetPlantByName :one
select * from plants
	where species_name like $1;

