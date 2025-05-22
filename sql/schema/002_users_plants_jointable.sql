-- +goose Up
create table users_plants (
  id uuid primary key,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  --
  created_by text not null,
  updated_by text not null,
  deleted_by text,
  --
  -- foreign keys
  plant_id uuid,
  user_id uuid,
  --
  -- table data
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
  foreign key (plant_id)
  references plants(id);

-- +goose Down
drop table users_plants;
