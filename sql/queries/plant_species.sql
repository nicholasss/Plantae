-- name: CreatePlantSpecies :one
insert into plant_species (
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
) values (
	gen_random_uuid(), now(), now(), $1, $2, $3, $4, $5, $6, $7
) returning *;

-- name: GetPlantSpeciesByName :one
select * from plant_species
	where species_name like $1
  limit 1;

-- name: GetPlantSpeciesByID :one
select * from plant_species
  where id = $1
  limit 1;

-- name: GetAllPlantSpeciesOrderedByUpdated :many
select * from plant_species
  order by updated_at desc;

-- name: GetAllPlantSpeciesOrderedByCreated :many
select * from plant_species
  order by created_at desc;

-- name: UpdatePlantSpeciesByID :exec
update plant_species
  set human_poison_toxic = $2,
	pet_poison_toxic = $3,
	human_edible = $4,
  pet_edible = $5
where id = $1;

-- name: MarkAsDeletedPlantSpeciesByID :exec
update plant_species
  set deleted_at = now(),
  deleted_by = $2
where id = $1;
