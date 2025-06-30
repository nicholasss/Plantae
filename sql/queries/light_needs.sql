-- name: CreateLightNeed :one
insert into light_needs (
	id,
  created_at, updated_at,
	created_by, updated_by,
  name, description
) values (
  gen_random_uuid(),
  now(), now(),
  $1, $1,
  $2, $3
) returning id, name, description;

-- name: ResetLightNeedsTable :exec
delete from light_needs;

-- name: MarkLightNeedAsDeletedByID :exec
update light_needs
  set
  deleted_at = now(),
  deleted_by = $2
where id = $1;

-- name: GetAllLightNeedsOrderedByCreated :many
select 
	id,
  name,
  description
from light_needs
  where deleted_at is null
  order by created_at desc;

-- name: UpdateLightNeedsByID :exec
update light_needs
  set updated_at = now(),
  updated_by = $2,
	name = $3,
  description = $4
where id = $1
  and deleted_at is null;
