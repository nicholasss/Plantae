-- name: CreatePlant :one
insert into plants (
	id, created_at, updated_at,
	created_by, updated_by, species_name,
	human_poison_toxic, pet_poison_toxic,
	human_edible, pet_edible
) values (
	gen_random_uuid(), now(), now(), $1, $2, $3, $4, $5, $6, $7
) returning *;

-- name: GetPlantByName :one
select * from plants
	where species_name like $1
  and deleted_at is null
  limit 1;

-- name: GetPlantByID :one
select * from plants
  where id = $1
  and deleted_at is null
  limit 1;

-- name: GetAllPlantsOrderedByUpdated :many
select * from plants
  where deleted_at is null
  order by updated_at desc;

-- name: GetAllPlantsOrderedByCreated :many
select * from plants
  where deleted_at is null
  order by created_at desc;

-- name: UpdatePlantsPropertiesByID :exec
update plants
  set human_poison_toxic = $2,
	pet_poison_toxic = $3,
	human_edible = $4,
  pet_edible = $5
where id = $1
  and deleted_at is null;

-- name: MarkPlantAsDeletedByID :exec
update plants
  set deleted_at = now(),
  deleted_by = $2
where id = $1;
