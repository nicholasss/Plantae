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
  iup.adoption_date,
  iup.name as plant_name,
  ps.id as species_id,
  ps.species_name as plant_species_name
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
with users_plant as (
  select
    id, plant_id, adoption_date, name
  from users_plants
  where
    deleted_at is null and
    user_id = $1
  order by updated_at desc
)
select
  up.id as users_plant_id,
  up.adoption_date,
  up.name as plant_name,
  ps.id as plant_species_id,
  ps.species_name
from
  users_plant as up
join
  plant_species as ps on up.plant_id = ps.id;

-- name: GetAllUsersPlantsOrderedByCreated :many
with users_plant as (
  select 
    id, plant_id, adoption_date, name
  from users_plants
  where
    deleted_at is null and
    user_id = $1
  order by created_at desc
)
select
  up.id as users_plant_id,
  up.adoption_date,
  up.name as plant_name,
  ps.id as plant_species_id,
  ps.species_name
from
  users_plant as up
join
  plant_species as ps on up.plant_id = ps.id;
