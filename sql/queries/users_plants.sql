-- name: GetAllUsersPlantsOrderedByUpdated :many
select * from users_plants
  where deleted_at is null
  order by updated_at desc;

-- name: GetAllUsersPlantsOrderedByCreated :many
select * from users_plants
  where deleted_at is null
  order by created_at desc;

