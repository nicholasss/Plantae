#
# Verify that server is online
GET http://localhost:8080/api/v1/health
HTTP 200
Content-Type: text/html; charset=UTF-8
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
# Create user account
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json; charset=UTF-8
```json
{
  "email": "{{lisa_email}}",
  "password": "{{lisa_password}}"
}
```
HTTP 201
Content-Type: application/json
[Captures]
lisa_id: jsonpath "$.id"
[Asserts]
jsonpath "$.isAdmin" == false
jsonpath "$.email" == "{{lisa_email}}"

#
# Login to user account
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json; charset=UTF-8
```json
{
  "email": "{{lisa_email}}",
  "password": "{{lisa_password}}"
}
```
HTTP 200
Content-Type: application/json
[Captures]
lisa_token: jsonpath "$.token"
lisa_refresh_token: jsonpath "$.refreshToken"
[Asserts]
jsonpath "$.isAdmin" == false
jsonpath "$.tokenExpiresAt" exists
jsonpath "$.refreshTokenExpiresAt" exists

#
# Promote user account to admin account
POST http://localhost:8080/api/v1/super-admin/promote-user
Authorization: SuperAdminToken {{super_admin_token}}
Content-Type: application/json; charset=UTF-8
```json
{
  "id": "{{lisa_id}}"
}
```
HTTP 204

# setup
# ========================================================================
# testing plant species

#
# Test empty plant species table
GET http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# Add a plant to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{1_species_name}}",
  "humanPoisonToxic": {{1_human_poison_toxic}},
  "petPoisonToxic": {{1_pet_poison_toxic}},
  "humanEdible": {{1_human_edible}},
  "petEdible": {{1_pet_edible}}
}
```
HTTP 204

#
# Requesting list of plant species in the table
# -- with only one 
GET http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
HTTP 200
[Captures]
1_plant_species_id: jsonpath "$[0].id"
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == {{lisa_id}}
jsonpath "$[0].updatedBy" == {{lisa_id}}
jsonpath "$[0].speciesName" == "{{1_species_name}}"
jsonpath "$[0].humanPoisonToxic" == {{1_human_poison_toxic}}
jsonpath "$[0].petPoisonToxic" == {{1_pet_poison_toxic}}
jsonpath "$[0].humanEdible" == {{1_human_edible}}
jsonpath "$[0].petEdible" == {{1_pet_edible}}

#
# Add a second plant to the plant_species
# -- with required details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{2_species_name}}"
}
```
HTTP 204

#
# Requesting list of plant species in the table
# -- with two records in table
GET http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
HTTP 200
[Captures]
2_plant_species_id: jsonpath "$[0].id"
[Asserts]
jsonpath "$" count == 2
jsonpath "$[0].speciesName" == "{{2_species_name}}"
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == "{{lisa_id}}"
jsonpath "$[0].updatedBy" == "{{lisa_id}}"

#
# Updating plant with new information
PUT http://localhost:8080/api/v1/admin/plant-species/{{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "humanPoisonToxic": {{2_human_poison_toxic}},
  "petPoisonToxic": {{2_pet_poison_toxic}},
  "humanEdible": {{2_human_edible}},
  "petEdible": {{2_pet_edible}}
}
```
HTTP 204

#
# Requesting list of plant species in the table
# -- with one record updated after it was created
GET http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
jsonpath "$" count == 2
jsonpath "$[0].speciesName" == "{{2_species_name}}"
jsonpath "$[0].humanPoisonToxic" == {{2_human_poison_toxic}}
jsonpath "$[0].petPoisonToxic" == {{2_pet_poison_toxic}}
jsonpath "$[0].humanEdible" == {{2_human_edible}}
jsonpath "$[0].petEdible" == {{2_pet_edible}}
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == "{{lisa_id}}"
jsonpath "$[0].updatedBy" == "{{lisa_id}}"

#
# Requesting deletion of a plants species record
DELETE http://localhost:8080/api/v1/admin/plant-species/{{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# Requesting list of plants after deletion of 1_plant_species
GET http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].speciesName" == "{{2_species_name}}"
jsonpath "$[0].humanPoisonToxic" == {{2_human_poison_toxic}}
jsonpath "$[0].petPoisonToxic" == {{2_pet_poison_toxic}}
jsonpath "$[0].humanEdible" == {{2_human_edible}}
jsonpath "$[0].petEdible" == {{2_pet_edible}}
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == "{{lisa_id}}"
jsonpath "$[0].updatedBy" == "{{lisa_id}}"

# testing plant species
# ========================================================================
# testing plant names

#
# Adding common names in english & spanish for 2_plant_species
POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1a_species_common_langcode}}",
  "commonName": "{{1a_species_common_name}}"
}
```
HTTP 204

POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1b_species_common_langcode}}",
  "commonName": "{{1b_species_common_name}}"
}
```
HTTP 204

POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1c_species_common_langcode}}",
  "commonName": "{{1c_species_common_name}}"
}
```
HTTP 204

POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1d_species_common_langcode}}",
  "commonName": "{{1d_species_common_name}}"
}
```
HTTP 204

POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1e_species_common_langcode}}",
  "commonName": "{{1e_species_common_name}}"
}
```
HTTP 204

POST http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
```json
{
  "plantID": "{{1_plant_species_id}}",
  "langCode": "{{1f_species_common_langcode}}",
  "commonName": "{{1f_species_common_name}}"
}
```
HTTP 204

#
# Get all common names of plants
GET http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 6

#
# Get all common names of plants, in english
GET http://localhost:8080/api/v1/admin/plant-names?lang=en
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 3
jsonpath "$[0].plantId" == "{{1_plant_species_id}}"
jsonpath "$[0].commonName" == "{{1c_species_common_name}}"
jsonpath "$[0].langCode" == "en"
jsonpath "$[1].plantId" == "{{1_plant_species_id}}"
jsonpath "$[1].commonName" == "{{1b_species_common_name}}"
jsonpath "$[1].langCode" == "en"
jsonpath "$[2].plantId" == "{{1_plant_species_id}}"
jsonpath "$[2].commonName" == "{{1a_species_common_name}}"
jsonpath "$[2].langCode" == "en"
[Captures]
1c_species_common_name_id: jsonpath "$[0].id"
1b_species_common_name_id: jsonpath "$[1].id"
1a_species_common_name_id: jsonpath "$[2].id"

#
# Get all common names of plants, in spanish
GET http://localhost:8080/api/v1/admin/plant-names?lang=es
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 3
jsonpath "$[0].plantId" == "{{1_plant_species_id}}"
jsonpath "$[0].commonName" == "{{1f_species_common_name}}"
jsonpath "$[0].langCode" == "es"
jsonpath "$[1].plantId" == "{{1_plant_species_id}}"
jsonpath "$[1].commonName" == "{{1e_species_common_name}}"
jsonpath "$[1].langCode" == "es"
jsonpath "$[2].plantId" == "{{1_plant_species_id}}"
jsonpath "$[2].commonName" == "{{1d_species_common_name}}"
jsonpath "$[2].langCode" == "es"
[Captures]
1f_species_common_name_id: jsonpath "$[0].id"
1e_species_common_name_id: jsonpath "$[1].id"
1d_species_common_name_id: jsonpath "$[2].id"

#
# Mark the english names as deleted
DELETE http://localhost:8080/api/v1/admin/plant-names/{{1a_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

DELETE http://localhost:8080/api/v1/admin/plant-names/{{1b_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

DELETE http://localhost:8080/api/v1/admin/plant-names/{{1c_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# Get all common names of plants, after deleting the english ones
GET http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 3
jsonpath "$[0].langCode" == "es"
jsonpath "$[1].langCode" == "es"
jsonpath "$[2].langCode" == "es"

#
# Mark the spanish names as deleted
DELETE http://localhost:8080/api/v1/admin/plant-names/{{1d_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

DELETE http://localhost:8080/api/v1/admin/plant-names/{{1e_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

DELETE http://localhost:8080/api/v1/admin/plant-names/{{1f_species_common_name_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# Get all common names of plants, after deleting the spanish ones
GET http://localhost:8080/api/v1/admin/plant-names
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" == null

# testing plant names
# ========================================================================
# testing plant types

#
# adding plant type 1 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{1_plant_type_name}}",
  "description": "{{1_plant_type_description}}"
}
```
HTTP 204

#
# adding plant type 2 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{2_plant_type_name}}",
  "description": "{{2_plant_type_description}}"
}
```
HTTP 204

#
# adding plant type 3 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{3_plant_type_name}}",
  "description": "{{3_plant_type_description}}"
}
```
HTTP 204

#
# adding plant type 4 with minimal information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{4_plant_type_name}}",
  "description": "{{4_plant_type_description}}"
}
```
HTTP 204

#
# listing plant types and capturing ids with minimal data
GET http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 4
[Captures]
1_plant_type_id: jsonpath "$[3].id"
2_plant_type_id: jsonpath "$[2].id"
3_plant_type_id: jsonpath "$[1].id"
4_plant_type_id: jsonpath "$[0].id"

#
# adding more information to plant type 1
PUT http://localhost:8080/api/v1/admin/plant-types/{{1_plant_type_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "maxTemperatureCelsius": {{1_plant_type_maxtc}},
  "minTemperatureCelsius": {{1_plant_type_mintc}},
  "maxHumidityPercent": {{1_plant_type_maxph}},
  "minHumidityPercent": {{1_plant_type_minph}},
  "soilOrganicMix": "{{1_plant_type_soilom}}",
  "soilGritMix": "{{1_plant_type_soilgm}}",
  "soilDrainageMix": "{{1_plant_type_soildm}}"
}
```
HTTP 204

#
# adding more information to plant type 2
PUT http://localhost:8080/api/v1/admin/plant-types/{{2_plant_type_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "maxTemperatureCelsius": {{2_plant_type_maxtc}},
  "minTemperatureCelsius": {{2_plant_type_mintc}},
  "maxHumidityPercent": {{2_plant_type_maxph}},
  "minHumidityPercent": {{2_plant_type_minph}},
  "soilOrganicMix": "{{2_plant_type_soilom}}",
  "soilGritMix": "{{2_plant_type_soilgm}}",
  "soilDrainageMix": "{{2_plant_type_soildm}}"
}
```
HTTP 204

#
# adding more information to plant type 3
PUT http://localhost:8080/api/v1/admin/plant-types/{{3_plant_type_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "maxTemperatureCelsius": {{3_plant_type_maxtc}},
  "minTemperatureCelsius": {{3_plant_type_mintc}},
  "maxHumidityPercent": {{3_plant_type_maxph}},
  "minHumidityPercent": {{3_plant_type_minph}},
  "soilOrganicMix": "{{3_plant_type_soilom}}",
  "soilGritMix": "{{3_plant_type_soilgm}}",
  "soilDrainageMix": "{{3_plant_type_soildm}}"
}
```
HTTP 204

#
# adding more information to plant type 4
PUT http://localhost:8080/api/v1/admin/plant-types/{{4_plant_type_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "maxTemperatureCelsius": {{4_plant_type_maxtc}},
  "minTemperatureCelsius": {{4_plant_type_mintc}},
  "maxHumidityPercent": {{4_plant_type_maxph}},
  "minHumidityPercent": {{4_plant_type_minph}},
  "soilOrganicMix": "{{4_plant_type_soilom}}",
  "soilGritMix": "{{4_plant_type_soilgm}}",
  "soilDrainageMix": "{{4_plant_type_soildm}}"
}
```
HTTP 204

#
# listing plant types and asserting responses
GET http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 4
jsonpath "$[3].id" == "{{1_plant_type_id}}"
jsonpath "$[3].name" == "{{1_plant_type_name}}"
jsonpath "$[3].description" == "{{1_plant_type_description}}"
jsonpath "$[3].maxTemperatureCelsius" == {{1_plant_type_maxtc}}
jsonpath "$[3].minTemperatureCelsius" == {{1_plant_type_mintc}}
jsonpath "$[3].maxHumidityPercent" == {{1_plant_type_maxph}}
jsonpath "$[3].minHumidityPercent" == {{1_plant_type_minph}}
jsonpath "$[3].soilOrganicMix" == "{{1_plant_type_soilom}}"
jsonpath "$[3].soilGritMix" == "{{1_plant_type_soilgm}}"
jsonpath "$[3].soilDrainageMix" == "{{1_plant_type_soildm}}"
jsonpath "$[2].id" == "{{2_plant_type_id}}"
jsonpath "$[2].name" == "{{2_plant_type_name}}"
jsonpath "$[2].description" == "{{2_plant_type_description}}"
jsonpath "$[2].maxTemperatureCelsius" == {{2_plant_type_maxtc}}
jsonpath "$[2].minTemperatureCelsius" == {{2_plant_type_mintc}}
jsonpath "$[2].maxHumidityPercent" == {{2_plant_type_maxph}}
jsonpath "$[2].minHumidityPercent" == {{2_plant_type_minph}}
jsonpath "$[2].soilOrganicMix" == "{{2_plant_type_soilom}}"
jsonpath "$[2].soilGritMix" == "{{2_plant_type_soilgm}}"
jsonpath "$[2].soilDrainageMix" == "{{2_plant_type_soildm}}"
jsonpath "$[1].id" == "{{3_plant_type_id}}"
jsonpath "$[1].name" == "{{3_plant_type_name}}"
jsonpath "$[1].description" == "{{3_plant_type_description}}"
jsonpath "$[1].maxTemperatureCelsius" == {{3_plant_type_maxtc}}
jsonpath "$[1].minTemperatureCelsius" == {{3_plant_type_mintc}}
jsonpath "$[1].maxHumidityPercent" == {{3_plant_type_maxph}}
jsonpath "$[1].minHumidityPercent" == {{3_plant_type_minph}}
jsonpath "$[1].soilOrganicMix" == "{{3_plant_type_soilom}}"
jsonpath "$[1].soilGritMix" == "{{3_plant_type_soilgm}}"
jsonpath "$[1].soilDrainageMix" == "{{3_plant_type_soildm}}"
jsonpath "$[0].id" == "{{4_plant_type_id}}"
jsonpath "$[0].name" == "{{4_plant_type_name}}"
jsonpath "$[0].description" == "{{4_plant_type_description}}"
jsonpath "$[0].maxTemperatureCelsius" == {{4_plant_type_maxtc}}
jsonpath "$[0].minTemperatureCelsius" == {{4_plant_type_mintc}}
jsonpath "$[0].maxHumidityPercent" == {{4_plant_type_maxph}}
jsonpath "$[0].minHumidityPercent" == {{4_plant_type_minph}}
jsonpath "$[0].soilOrganicMix" == "{{4_plant_type_soilom}}"
jsonpath "$[0].soilGritMix" == "{{4_plant_type_soilgm}}"
jsonpath "$[0].soilDrainageMix" == "{{4_plant_type_soildm}}"

#
# deleting plant type 1
DELETE http://localhost:8080/api/v1/admin/plant-types/{{1_plant_type_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# deleting plant type 2
DELETE http://localhost:8080/api/v1/admin/plant-types/{{2_plant_type_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# listing plant types and asserting responses after deleting type 1 & 2
GET http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$" count == 2
jsonpath "$[1].id" == "{{3_plant_type_id}}"
jsonpath "$[1].name" == "{{3_plant_type_name}}"
jsonpath "$[1].description" == "{{3_plant_type_description}}"
jsonpath "$[1].maxTemperatureCelsius" == {{3_plant_type_maxtc}}
jsonpath "$[1].minTemperatureCelsius" == {{3_plant_type_mintc}}
jsonpath "$[1].maxHumidityPercent" == {{3_plant_type_maxph}}
jsonpath "$[1].minHumidityPercent" == {{3_plant_type_minph}}
jsonpath "$[1].soilOrganicMix" == "{{3_plant_type_soilom}}"
jsonpath "$[1].soilGritMix" == "{{3_plant_type_soilgm}}"
jsonpath "$[1].soilDrainageMix" == "{{3_plant_type_soildm}}"
jsonpath "$[0].id" == "{{4_plant_type_id}}"
jsonpath "$[0].name" == "{{4_plant_type_name}}"
jsonpath "$[0].description" == "{{4_plant_type_description}}"
jsonpath "$[0].maxTemperatureCelsius" == {{4_plant_type_maxtc}}
jsonpath "$[0].minTemperatureCelsius" == {{4_plant_type_mintc}}
jsonpath "$[0].maxHumidityPercent" == {{4_plant_type_maxph}}
jsonpath "$[0].minHumidityPercent" == {{4_plant_type_minph}}
jsonpath "$[0].soilOrganicMix" == "{{4_plant_type_soilom}}"
jsonpath "$[0].soilGritMix" == "{{4_plant_type_soilgm}}"
jsonpath "$[0].soilDrainageMix" == "{{4_plant_type_soildm}}"

# testing plant types
# ========================================================================
# testing setting plant to types

#
# Reset plant_species table
POST http://localhost:8080/api/v1/super-admin/reset-plant-species
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset plant_types table
POST http://localhost:8080/api/v1/super-admin/reset-plant-types
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Add plant 1 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{1_species_name}}",
  "humanPoisonToxic": {{1_human_poison_toxic}},
  "petPoisonToxic": {{1_pet_poison_toxic}},
  "humanEdible": {{1_human_edible}},
  "petEdible": {{1_pet_edible}}
}
```
HTTP 204

#
# Add plant 2 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{2_species_name}}",
  "humanPoisonToxic": {{2_human_poison_toxic}},
  "petPoisonToxic": {{2_pet_poison_toxic}},
  "humanEdible": {{2_human_edible}},
  "petEdible": {{2_pet_edible}}
}
```
HTTP 204

#
# adding plant type 1 with full information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
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
HTTP 204

#
# adding plant type 2 with full information
POST http://localhost:8080/api/v1/admin/plant-types
Authorization: Bearer {{lisa_token}}
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
HTTP 204


#
# set 1_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{3_plant_type_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# set 2_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{4_plant_type_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# unset 1_plant_species as type
DELETE http://localhost:8080/api/v1/admin/plant-types/link/{{3_plant_type_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# unset 2_plant_species as type
DELETE http://localhost:8080/api/v1/admin/plant-types/link/{{4_plant_type_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# again, set 1_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{3_plant_type_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# again, set 2_plant_species as type
POST http://localhost:8080/api/v1/admin/plant-types/link/{{4_plant_type_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# testing setting plant to types
# ========================================================================
# testing setting plant to light need

#
# Reset plant_species table
POST http://localhost:8080/api/v1/super-admin/reset-plant-species
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Reset light_needs table
POST http://localhost:8080/api/v1/super-admin/reset-light
Authorization: SuperAdminToken {{super_admin_token}}
HTTP 204

#
# Add plant 1 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{1_species_name}}",
  "humanPoisonToxic": {{1_human_poison_toxic}},
  "petPoisonToxic": {{1_pet_poison_toxic}},
  "humanEdible": {{1_human_edible}},
  "petEdible": {{1_pet_edible}}
}
```
HTTP 204

#
# Add plant 2 to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plant-species
Authorization: Bearer {{lisa_token}}
```json
{
  "speciesName": "{{2_species_name}}",
  "humanPoisonToxic": {{2_human_poison_toxic}},
  "petPoisonToxic": {{2_pet_poison_toxic}},
  "humanEdible": {{2_human_edible}},
  "petEdible": {{2_pet_edible}}
}
```
HTTP 204

#
# add light 1 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{1_light_name}}",
  "description": "{{1_light_description}}"
}
```
HTTP 204

#
# add light 2 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{2_light_name}}",
  "description": "{{2_light_description}}"
}
```
HTTP 204

#
# add light 2 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{3_light_name}}",
  "description": "{{3_light_description}}"
}
```
HTTP 204

#
# add light 4 to light_needs table
POST http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{4_light_name}}",
  "description": "{{4_light_description}}"
}
```
HTTP 204

#
# get list of light needs
GET http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$" count == 4
jsonpath "$[0].name" == "{{4_light_name}}"
jsonpath "$[0].description" == "{{4_light_description}}"
jsonpath "$[1].name" == "{{3_light_name}}"
jsonpath "$[1].description" == "{{3_light_description}}"
jsonpath "$[2].name" == "{{2_light_name}}"
jsonpath "$[2].description" == "{{2_light_description}}"
jsonpath "$[3].name" == "{{1_light_name}}"
jsonpath "$[3].description" == "{{1_light_description}}"
[Captures]
1_light_id: jsonpath "$[3].id"
2_light_id: jsonpath "$[2].id"
3_light_id: jsonpath "$[1].id"
4_light_id: jsonpath "$[0].id"

# 
# update light 1
PUT http://localhost:8080/api/v1/admin/light/{{1_light_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{1_light_name}}",
  "description": "{{1_light_description_alt}}"
}
```
HTTP 204

# 
# update light 2
PUT http://localhost:8080/api/v1/admin/light/{{2_light_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{2_light_name}}",
  "description": "{{2_light_description_alt}}"
}
```
HTTP 204

# 
# update light 3
PUT http://localhost:8080/api/v1/admin/light/{{3_light_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{3_light_name}}",
  "description": "{{3_light_description_alt}}"
}
```
HTTP 204

# 
# update light 4
PUT http://localhost:8080/api/v1/admin/light/{{4_light_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "name": "{{4_light_name}}",
  "description": "{{4_light_description_alt}}"
}
```
HTTP 204

#
# get list of light needs
GET http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$" count == 4
jsonpath "$[0].id" == "{{4_light_id}}"
jsonpath "$[0].name" == "{{4_light_name}}"
jsonpath "$[0].description" == "{{4_light_description_alt}}"
jsonpath "$[1].id" == "{{3_light_id}}"
jsonpath "$[1].name" == "{{3_light_name}}"
jsonpath "$[1].description" == "{{3_light_description_alt}}"
jsonpath "$[2].id" == "{{2_light_id}}"
jsonpath "$[2].name" == "{{2_light_name}}"
jsonpath "$[2].description" == "{{2_light_description_alt}}"
jsonpath "$[3].id" == "{{1_light_id}}"
jsonpath "$[3].name" == "{{1_light_name}}"
jsonpath "$[3].description" == "{{1_light_description_alt}}"

# 
# delete light 1
DELETE http://localhost:8080/api/v1/admin/light/{{1_light_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# delete light 2
DELETE http://localhost:8080/api/v1/admin/light/{{2_light_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# delete light 3
DELETE http://localhost:8080/api/v1/admin/light/{{3_light_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

# 
# delete light 4
DELETE http://localhost:8080/api/v1/admin/light/{{4_light_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# get empty list
GET http://localhost:8080/api/v1/admin/light
Authorization: Bearer {{lisa_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$" == null

# testing setting plant to light need
# ========================================================================
# testing water need

#
# adding water record for plant 1
POST http://localhost:8080/api/v1/admin/water
Authorization: Bearer {{lisa_token}}
```json
{
  "plantType": "{{1_plant_water_type}}",
  "description": "{{1_plant_water_description}}",
  "drySoilMM": {{1_plant_water_mm}}
}
```
HTTP 200
[Captures]
1_water_id: jsonpath "$.id"

#
# adding water record for plant 2
POST http://localhost:8080/api/v1/admin/water
Authorization: Bearer {{lisa_token}}
```json
{
  "plantType": "{{2_plant_water_type}}",
  "description": "{{2_plant_water_description}}",
  "drySoilDays": {{2_plant_water_days}}
}
```
HTTP 200
[Captures]
2_water_id: jsonpath "$.id"

#
# set water 1 to plant 1
POST http://localhost:8080/api/v1/admin/water/link/{{1_water_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# set water 2 to plant 2
POST http://localhost:8080/api/v1/admin/water/link/{{2_water_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# unset water 1 to plant 1
DELETE http://localhost:8080/api/v1/admin/water/link/{{1_water_id}}?plant-species-id={{1_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# unset water 2 to plant 2
DELETE http://localhost:8080/api/v1/admin/water/link/{{2_water_id}}?plant-species-id={{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204


#
# delete water 1
DELETE http://localhost:8080/api/v1/admin/water/{{1_water_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# delete water 2
DELETE http://localhost:8080/api/v1/admin/water/{{2_water_id}}
Authorization: Bearer {{lisa_token}}
HTTP 204
