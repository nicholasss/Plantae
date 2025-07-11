-- name: CreateWaterDryMM :one
insert into water_needs (
	id,
  created_at, updated_at,
	created_by, updated_by,
  plant_type, description,
  dry_soil_mm
) values (
  gen_random_uuid(),
  now(), now(),
  $1, $1,
  $2, $3,
  $4
) returning *;

-- name: CreateWaterDryDays :one
insert into water_needs (
	id,
  created_at, updated_at,
	created_by, updated_by,
  plant_type, description,
  dry_soil_days
) values (
  gen_random_uuid(),
  now(), now(),
  $1, $1,
  $2, $3,
  $4
) returning *;

-- name: ResetWaterNeedsTable :exec
delete from water_needs;

-- name: MarkWaterNeedAsDeletedByID :exec
update water_needs
  set
  deleted_at = now(),
  deleted_by = $2,
  updated_at = now(),
  updated_by = $2
where id = $1;

-- name: GetAllWaterNeedsOrderedByCreated :many
select
  id,
  plant_type,
  description,
  dry_soil_mm,
  dry_soil_days
from water_needs
  where deleted_at is null
  order by created_at desc;

-- name: UpdateWaterDryMMNeedsByID :exec
update water_needs
  set updated_at = now(),
  updated_by = $2,
  description = $3,
  dry_soil_mm = $4
where id = $1
  and deleted_at is null;

-- name: UpdateWaterDryDaysNeedsByID :exec
update water_needs
  set updated_at = now(),
  updated_by = $2,
  description = $3,
  dry_soil_days = $4
where id = $1
  and deleted_at is null;
