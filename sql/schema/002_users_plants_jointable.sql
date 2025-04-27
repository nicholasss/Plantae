-- +goose Up
create table users_plants (
  id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by timestamp with time zone not null,
  updated_by timestamp with time zone not null,
  deleted_by timestamp with time zone,
  --
  plant_species_id uuid,
  user_id uuid,
  --
  --
  adoption_date timestamp with time zone,
  name text
);

alter table users_plants
  add constraint fk_users
  foreign key (user_id)
  references users(id)
  on delete cascade;

alter table users_plants
  add constraint fk_plant_species
  foreign key (plant_species_id)
  references plant_species(id);

-- +goose Down
drop table users_plants;
