
-- name: GetAllViewPlantsOrderedByUpdated :many
select
  ps.id as plant_species_id,
  ps.species_name as plant_species_name,
  ps.human_poison_toxic,
  ps.pet_poison_toxic,
  ps.human_edible,
  ps.pet_edible,
  max(pn.lang_code) as lang_code,
  string_agg(pn.common_name, ', ') as common_names,
  pt.name as plant_type_name,
  pt.description as plant_type_description,
  ln.name as light_need_name,
  ln.description as light_need_description,
  wn.plant_type as water_need_type,
  wn.description as water_need_description,
  wn.dry_soil_mm as water_need_dry_soil_mm,
  wn.dry_soil_days as water_need_dry_soil_days
from
  plant_species as ps
join
  plant_names as pn on ps.id = pn.plant_id
left join
  plant_types as pt on ps.plant_type_id = pt.id
left join
  light_needs as ln on ps.light_needs_id = ln.id
left join
  water_needs as wn on ps.water_needs_id = wn.id
where
  pn.lang_code = $1 and
  ps.deleted_at is null and
  pn.deleted_at is null and
  ln.deleted_at is null and
  wn.deleted_at is null
group by
  ps.id,
  ps.species_name,
  ps.updated_at,
  ps.deleted_at,
  pn.deleted_at,
  pt.name,
  pt.description,
  ln.name,
  ln.description,
  wn.plant_type,
  wn.description,
  wn.dry_soil_mm,
  wn.dry_soil_days
order by
  ps.updated_at desc;

