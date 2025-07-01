-- name: CreateUsersPlants :one
with inserted_users_plant as (
  insert into users_plants (
    id,
    created_at, updated_at,
    created_by, updated_by,
    plant_id, user_id,
    adoption_date, name
  ) values (
    gen_random_uuid(),
    now(), now(),
    $1, $1,
    $2, $3,
    $4, $5
  ) returning 
    id, plant_id, user_id,
    adoption_date, name
  )
select
  iup.id as users_plant_id,
  iup.user_id,
  iup.adoption_date,
  iup.name,
  ps.id as species_id,
  ps.species_name
from
  inserted_users_plant as iup
join
  plant_species as ps on iup.plant_id = ps.id;

-- name: UpdateUsersPlantByID :exec
update users_plants
set updated_at = now(),
  updated_by = $2,
  adoption_date = $3,
  name = $4
where id = $1
  and deleted_at is null;

-- name: GetAllUsersPlantsOrderedByUpdated :many
select
  id as users_plant_id, plant_id, adoption_date, name
from users_plants
where
  deleted_at is null and
  user_id = $1
order by updated_at desc;

-- name: GetAllUsersPlantsOrderedByCreated :many
select 
  id as users_plant_id, plant_id, adoption_date, name
from users_plants
where
  deleted_at is null and
  user_id = $1
order by created_at desc;

