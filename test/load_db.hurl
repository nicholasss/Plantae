#
# Verify that server is online
GET http://localhost:8080/api/v1/health
HTTP 200
Content-Type: text/html; charset=utf-8
[Asserts]
xpath "string(/html/body)" contains "OK"

#
# Reset user table
POST http://localhost:8080/api/v1/super-admin/reset-users
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset plant_species table
POST http://localhost:8080/api/v1/super-admin/reset-plant-species
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset plant_names table
POST http://localhost:8080/api/v1/super-admin/reset-plant-names
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset plant_types table
POST http://localhost:8080/api/v1/super-admin/reset-plant-types
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset water_needs table
POST http://localhost:8080/api/v1/super-admin/reset-water
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset light_needs table
POST http://localhost:8080/api/v1/super-admin/reset-light
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Create admin account
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json; charset=utf-8
```json
{
  "email": "{{lisa_email}}",
  "password": "{{lisa_password}}",
  "langCodePref": "{{lisa_lang_code}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
lisa_id: jsonpath "$.id"

#
# Login to user account
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json; charset=utf-8
```json
{
  "email": "{{lisa_email}}",
  "password": "{{lisa_password}}"
}
```
HTTP 200
Content-Type: application/json; charset=utf-8
[Captures]
lisa_token: jsonpath "$.token"
lisa_refresh_token: jsonpath "$.refreshToken"
[Asserts]
jsonpath "$.langCodePref" == "{{lisa_lang_code}}"

#
# Promote user account to admin account
POST http://localhost:8080/api/v1/super-admin/promote-user
Authorization: SuperAdminToken {{super_admin_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "id": "{{lisa_id}}"
}
```
HTTP 200
Content-Type: application/json; charset=utf-8

#
# Create user account
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json; charset=utf-8
```json
{
  "email": "{{craig_email}}",
  "password": "{{craig_password}}",
  "langCodePref": "{{craig_lang_code}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
craig_id: jsonpath "$.id"

#
# Login to user account
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json; charset=utf-8
```json
{
  "email": "{{craig_email}}",
  "password": "{{craig_password}}"
}
```
HTTP 200
Content-Type: application/json; charset=utf-8
[Captures]
craig_token: jsonpath "$.token"
craig_refresh_token: jsonpath "$.refreshToken"

# setup
# ========================================================================
# load plant data into database

#
# Add plant 1 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "speciesName": "{{1_species_name}}",
  "humanPoisonToxic": {{1_human_poison_toxic}},
  "petPoisonToxic": {{1_pet_poison_toxic}},
  "humanEdible": {{1_human_edible}},
  "petEdible": {{1_pet_edible}}
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
1_plant_species_id: jsonpath "$.id"
[Asserts]
jsonpath "$.speciesName" == "{{1_species_name}}"
jsonpath "$.humanPoisonToxic" == {{1_human_poison_toxic}}
jsonpath "$.petPoisonToxic" == {{1_pet_poison_toxic}}
jsonpath "$.humanEdible" == {{1_human_edible}}
jsonpath "$.petEdible" == {{1_pet_edible}}

#
# Add plant 2 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "speciesName": "{{2_species_name}}",
  "humanPoisonToxic": {{2_human_poison_toxic}},
  "petPoisonToxic": {{2_pet_poison_toxic}},
  "humanEdible": {{2_human_edible}},
  "petEdible": {{2_pet_edible}}
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
2_plant_species_id: jsonpath "$.id"
[Asserts]
jsonpath "$.speciesName" == "{{2_species_name}}"
jsonpath "$.humanPoisonToxic" == {{2_human_poison_toxic}}
jsonpath "$.petPoisonToxic" == {{2_pet_poison_toxic}}
jsonpath "$.humanEdible" == {{2_human_edible}}
jsonpath "$.petEdible" == {{2_pet_edible}}

# Adding species
# ==============
# Adding names

#
# Adding common names in english & spanish for 1_plant_species
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1a_species_common_langcode}}",
  "commonName": "{{1a_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1a_species_common_langcode}}"
jsonpath "$.commonName" == "{{1a_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1b_species_common_langcode}}",
  "commonName": "{{1b_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1b_species_common_langcode}}"
jsonpath "$.commonName" == "{{1b_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1c_species_common_langcode}}",
  "commonName": "{{1c_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1c_species_common_langcode}}"
jsonpath "$.commonName" == "{{1c_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1d_species_common_langcode}}",
  "commonName": "{{1d_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1d_species_common_langcode}}"
jsonpath "$.commonName" == "{{1d_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1e_species_common_langcode}}",
  "commonName": "{{1e_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1e_species_common_langcode}}"
jsonpath "$.commonName" == "{{1e_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1f_species_common_langcode}}",
  "commonName": "{{1f_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{1_plant_species_id}}"
jsonpath "$.langCode" == "{{1f_species_common_langcode}}"
jsonpath "$.commonName" == "{{1f_species_common_name}}"

#
# Adding common names in english & spanish for 2_plant_species
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2a_species_common_langcode}}",
  "commonName": "{{2a_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2a_species_common_langcode}}"
jsonpath "$.commonName" == "{{2a_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2b_species_common_langcode}}",
  "commonName": "{{2b_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2b_species_common_langcode}}"
jsonpath "$.commonName" == "{{2b_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2c_species_common_langcode}}",
  "commonName": "{{2c_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2c_species_common_langcode}}"
jsonpath "$.commonName" == "{{2c_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2d_species_common_langcode}}",
  "commonName": "{{2d_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2d_species_common_langcode}}"
jsonpath "$.commonName" == "{{2d_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2e_species_common_langcode}}",
  "commonName": "{{2e_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2e_species_common_langcode}}"
jsonpath "$.commonName" == "{{2e_species_common_name}}"

#
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantID": "{{2_plant_species_id}}",
  "langCode": "{{2f_species_common_langcode}}",
  "commonName": "{{2f_species_common_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.plantID" == "{{2_plant_species_id}}"
jsonpath "$.langCode" == "{{2f_species_common_langcode}}"
jsonpath "$.commonName" == "{{2f_species_common_name}}"

# Adding names
# ==============
# Adding types

#
# adding plant type 1 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{1_plant_type_name}}",
  "description": "{{1_plant_type_description}}",
  "maxTemperatureCelsius": {{1_plant_type_maxtc}},
  "minTemperatureCelsius": {{1_plant_type_mintc}},
  "maxHumidityPercent": {{1_plant_type_maxph}},
  "minHumidityPercent": {{1_plant_type_minph}},
  "soilOrganicMix": "{{1_plant_type_soilom}}",
  "soilGritMix": "{{1_plant_type_soilgm}}",
  "soilDrainageMix": "{{1_plant_type_soildm}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
1_plant_type_id: jsonpath "$.id"
[Asserts]
jsonpath "$.name" == "{{1_plant_type_name}}"
jsonpath "$.description" == "{{1_plant_type_description}}"
jsonpath "$.maxTemperatureCelsius" == {{1_plant_type_maxtc}}
jsonpath "$.minTemperatureCelsius" == {{1_plant_type_mintc}}
jsonpath "$.maxHumidityPercent" == {{1_plant_type_maxph}}
jsonpath "$.minHumidityPercent" == {{1_plant_type_minph}}
jsonpath "$.soilOrganicMix" == "{{1_plant_type_soilom}}"
jsonpath "$.soilGritMix" == "{{1_plant_type_soilgm}}"
jsonpath "$.soilDrainageMix" == "{{1_plant_type_soildm}}"

#
# adding plant type 2 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{2_plant_type_name}}",
  "description": "{{2_plant_type_description}}",
  "maxTemperatureCelsius": {{2_plant_type_maxtc}},
  "minTemperatureCelsius": {{2_plant_type_mintc}},
  "maxHumidityPercent": {{2_plant_type_maxph}},
  "minHumidityPercent": {{2_plant_type_minph}},
  "soilOrganicMix": "{{2_plant_type_soilom}}",
  "soilGritMix": "{{2_plant_type_soilgm}}",
  "soilDrainageMix": "{{2_plant_type_soildm}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
2_plant_type_id: jsonpath "$.id"
[Asserts]
jsonpath "$.name" == "{{2_plant_type_name}}"
jsonpath "$.description" == "{{2_plant_type_description}}"
jsonpath "$.maxTemperatureCelsius" == {{2_plant_type_maxtc}}
jsonpath "$.minTemperatureCelsius" == {{2_plant_type_mintc}}
jsonpath "$.maxHumidityPercent" == {{2_plant_type_maxph}}
jsonpath "$.minHumidityPercent" == {{2_plant_type_minph}}
jsonpath "$.soilOrganicMix" == "{{2_plant_type_soilom}}"
jsonpath "$.soilGritMix" == "{{2_plant_type_soilgm}}"
jsonpath "$.soilDrainageMix" == "{{2_plant_type_soildm}}"

#
# adding plant type 3 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{3_plant_type_name}}",
  "description": "{{3_plant_type_description}}",
  "maxTemperatureCelsius": {{3_plant_type_maxtc}},
  "minTemperatureCelsius": {{3_plant_type_mintc}},
  "maxHumidityPercent": {{3_plant_type_maxph}},
  "minHumidityPercent": {{3_plant_type_minph}},
  "soilOrganicMix": "{{3_plant_type_soilom}}",
  "soilGritMix": "{{3_plant_type_soilgm}}",
  "soilDrainageMix": "{{3_plant_type_soildm}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
3_plant_type_id: jsonpath "$.id"
[Asserts]
jsonpath "$.name" == "{{3_plant_type_name}}"
jsonpath "$.description" == "{{3_plant_type_description}}"
jsonpath "$.maxTemperatureCelsius" == {{3_plant_type_maxtc}}
jsonpath "$.minTemperatureCelsius" == {{3_plant_type_mintc}}
jsonpath "$.maxHumidityPercent" == {{3_plant_type_maxph}}
jsonpath "$.minHumidityPercent" == {{3_plant_type_minph}}
jsonpath "$.soilOrganicMix" == "{{3_plant_type_soilom}}"
jsonpath "$.soilGritMix" == "{{3_plant_type_soilgm}}"
jsonpath "$.soilDrainageMix" == "{{3_plant_type_soildm}}"

#
# adding plant type 4 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{4_plant_type_name}}",
  "description": "{{4_plant_type_description}}",
  "maxTemperatureCelsius": {{4_plant_type_maxtc}},
  "minTemperatureCelsius": {{4_plant_type_mintc}},
  "maxHumidityPercent": {{4_plant_type_maxph}},
  "minHumidityPercent": {{4_plant_type_minph}},
  "soilOrganicMix": "{{4_plant_type_soilom}}",
  "soilGritMix": "{{4_plant_type_soilgm}}",
  "soilDrainageMix": "{{4_plant_type_soildm}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
4_plant_type_id: jsonpath "$.id"
[Asserts]
jsonpath "$.name" == "{{4_plant_type_name}}"
jsonpath "$.description" == "{{4_plant_type_description}}"
jsonpath "$.maxTemperatureCelsius" == {{4_plant_type_maxtc}}
jsonpath "$.minTemperatureCelsius" == {{4_plant_type_mintc}}
jsonpath "$.maxHumidityPercent" == {{4_plant_type_maxph}}
jsonpath "$.minHumidityPercent" == {{4_plant_type_minph}}
jsonpath "$.soilOrganicMix" == "{{4_plant_type_soilom}}"
jsonpath "$.soilGritMix" == "{{4_plant_type_soilgm}}"
jsonpath "$.soilDrainageMix" == "{{4_plant_type_soildm}}"

#
# set 1_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{3_plant_type_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8

# 
# set 2_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{4_plant_type_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8

# Adding types
# ==============
# Adding light needs

#
# add light 1 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{1_light_name}}",
  "description": "{{1_light_description}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
1_light_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "{{1_light_name}}"
jsonpath "$.description" == {{1_light_description}}

#
# add light 2 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{2_light_name}}",
  "description": "{{2_light_description}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
2_light_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "{{2_light_name}}"
jsonpath "$.description" == {{2_light_description}}

#
# add light 3 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{3_light_name}}",
  "description": "{{3_light_description}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
3_light_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "{{3_light_name}}"
jsonpath "$.description" == {{3_light_description}}

#
# add light 4 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "name": "{{4_light_name}}",
  "description": "{{4_light_description}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
4_light_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "{{4_light_name}}"
jsonpath "$.description" == {{4_light_description}}

#
# set 1_plant_species as light need
POST http://localhost:8080/api/v1/admin/light/link/{{2_light_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8

# 
# set 2_plant_species as light need
POST http://localhost:8080/api/v1/admin/light/link/{{1_light_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8

# Adding light needs
# ==============
# Adding water needs

#
# adding water record for plant 1
POST http://localhost:8080/api/v1/admin/water
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantType": "{{1_plant_water_type}}",
  "description": "{{1_plant_water_description}}",
  "drySoilMM": {{1_plant_water_mm}}
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
1_water_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.plantType" == "{{1_plant_water_type}}"
jsonpath "$.description" == "{{1_plant_water_description}}"
jsonpath "$.drySoilMM" == {{1_plant_water_mm}}

#
# adding water record for plant 2
POST http://localhost:8080/api/v1/admin/water
Authorization: Bearer {{lisa_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantType": "{{2_plant_water_type}}",
  "description": "{{2_plant_water_description}}",
  "drySoilDays": {{2_plant_water_days}}
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
2_water_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.plantType" == "{{2_plant_water_type}}"
jsonpath "$.description" == "{{2_plant_water_description}}"
jsonpath "$.drySoilDays" == {{2_plant_water_days}}

#
# set water 1 to plant 1
POST http://localhost:8080/api/v1/admin/water/link/{{1_water_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.waterNeedID" == "{{1_water_id}}"
jsonpath "$.plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$.plantSpeciesName" == "{{1_species_name}}"

#
# set water 2 to plant 2
POST http://localhost:8080/api/v1/admin/water/link/{{2_water_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$.waterNeedID" == "{{2_water_id}}"
jsonpath "$.plantSpeciesID" == "{{2_plant_species_id}}"
jsonpath "$.plantSpeciesName" == "{{2_species_name}}"
