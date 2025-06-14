-- name: CreatePlantName :one
insert into plant_names (
	id,
  created_at, updated_at,
	created_by, updated_by,
  plant_id,
  lang_code,
  common_name
) values (
  gen_random_uuid(),
  now(), now(),
  $1, $1,
  $2,
  $3,
  $4
) returning *;

-- name: ResetPlantNamesTable :exec
delete from plant_names;

-- name: MarkPlantNameAsDeletedByID :exec
update plant_names
  set
  deleted_at = now(),
  deleted_by = $2
where id = $1;

-- name: GetAllPlantNamesOrderedByCreated :many
select 
	id,
  plant_id,
  lang_code,
  common_name
from plant_names
  where deleted_at is null
  order by created_at desc;

-- name: GetAllPlantNamesForLanguageOrderedByCreated :many
select 
	id,
  plant_id,
  lang_code,
  common_name
from plant_names
  where lang_code ilike $1
  and deleted_at is null
  order by created_at desc;
