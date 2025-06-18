-- name: CreatePlantSpecies :one
insert into plant_species (
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
) values (
	gen_random_uuid(), now(), now(), $1, $2, $3, $4, $5, $6, $7
) returning *;

-- name: ResetPlantSpeciesTable :exec
delete from plant_species;

-- name: GetPlantSpeciesByName :one
select 
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
from plant_species
	where species_name like $1
  and deleted_at is null
  limit 1;

-- name: GetPlantSpeciesByID :one
select 
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
from plant_species
  where id = $1
  and deleted_at is null
  limit 1;

-- name: GetAllPlantSpeciesOrderedByUpdated :many
select 
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
from plant_species
  where deleted_at is null
  order by updated_at desc;

-- name: GetAllPlantSpeciesOrderedByCreated :many
select 
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
from plant_species
  where deleted_at is null
  order by created_at desc;

-- name: UpdatePlantSpeciesPropertiesByID :exec
update plant_species
  set updated_at = now(),
  updated_by = $2,
  human_poison_toxic = $3,
	pet_poison_toxic = $4,
	human_edible = $5,
  pet_edible = $6
where id = $1
  and deleted_at is null;

-- name: MarkPlantSpeciesAsDeletedByID :exec
update plant_species
  set deleted_at = now(),
  deleted_by = $2
where id = $1;

-- name: SetPlantSpeciesAsType :exec
update plant_species
  set plant_type_id = $2
where
  id = $1 and
  deleted_by is null;
