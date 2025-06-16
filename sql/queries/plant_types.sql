-- name: CreatePlantType :one
insert into plant_types (
	id,
  created_at, updated_at,
	created_by, updated_by,
  name, description,
  max_temperature_celsius, min_temperature_celsius,
  max_humidity_percent, min_humidity_percent,
  soil_organic_mix, soil_grit_mix, soil_drainage_mix
) values (
  gen_random_uuid(),
  now(), now(),
  $1, $1,
  $2, $3,
  $4, $5,
  $6, $7,
  $8, $9, $10
) returning *;

-- name: ResetPlantTypesTable :exec
delete from plant_types;

-- name: MarkPlantTypeAsDeletedByID :exec
update plant_types
  set
  deleted_at = now(),
  deleted_by = $2
where id = $1;

-- name: UpdatePlantTypesPropertiesByID :exec
update plant_types
  set updated_at = now(),
  max_temperature_celsius = $2,
  min_temperature_celsius = $3,
	max_humidity_percent = $4,
	min_humidity_percent = $5,
  soil_organic_mix = $6,
  soil_grit_mix = $7,
  soil_drainage_mix = $8
where id = $1
  and deleted_at is null;

-- name: GetAllPlantTypesOrderedByCreated :many
select 
	id,
  created_at, updated_at,
	created_by, updated_by,
  name, description,
  max_temperature_celsius, min_temperature_celsius,
  max_humidity_percent, min_humidity_percent,
  soil_organic_mix, soil_grit_mix, soil_drainage_mix
from plant_types
  where deleted_at is null
  order by created_at desc;
