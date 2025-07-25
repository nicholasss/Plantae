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
# Create user account
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
[Asserts]
jsonpath "$.id" exists
jsonpath "$.langCodePref" == "{{lisa_lang_code}}"
jsonpath "$.joinDate" isIsoDate
jsonpath "$.isAdmin" == false
jsonpath "$.email" == "{{lisa_email}}"

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
jsonpath "$.joinDate" isIsoDate
jsonpath "$.isAdmin" == false
jsonpath "$.tokenExpiresAt" exists
jsonpath "$.refreshTokenExpiresAt" exists

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
# Add first plant to the plant_species table
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
jsonpath "$.id" exists
jsonpath "$.speciesName" == "{{1_species_name}}"
jsonpath "$.humanPoisonToxic" == {{1_human_poison_toxic}}
jsonpath "$.petPoisonToxic" == {{1_pet_poison_toxic}}
jsonpath "$.humanEdible" == {{1_human_edible}}
jsonpath "$.petEdible" == {{1_pet_edible}}

#
# Add second plant to the plant_species table
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
jsonpath "$.id" exists
jsonpath "$.speciesName" == "{{2_species_name}}"
jsonpath "$.humanPoisonToxic" == {{2_human_poison_toxic}}
jsonpath "$.petPoisonToxic" == {{2_pet_poison_toxic}}
jsonpath "$.humanEdible" == {{2_human_edible}}
jsonpath "$.petEdible" == {{2_pet_edible}}

# setup
# ========================================================================
# testing users plants

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
[Asserts]
jsonpath "$.id" exists
jsonpath "$.langCodePref" == "{{craig_lang_code}}"
jsonpath "$.joinDate" isIsoDate
jsonpath "$.isAdmin" == false
jsonpath "$.email" == "{{craig_email}}"

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
[Asserts]
jsonpath "$.langCodePref" == "{{craig_lang_code}}"
jsonpath "$.joinDate" isIsoDate
jsonpath "$.isAdmin" == false
jsonpath "$.tokenExpiresAt" exists
jsonpath "$.refreshTokenExpiresAt" exists

#
# List empty table
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{lisa_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$" == null

#
# Create first users plant
POST http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantSpeciesID": "{{1_plant_species_id}}",
  "adoptionDate": "{{1_plant_adoption}}",
  "plantName": "{{1_plant_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
1_my_plant_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$.plantSpeciesName" == "{{1_species_name}}"
jsonpath "$.adoptionDate" == "{{1_plant_adoption}}"
jsonpath "$.plantName" == "{{1_plant_name}}"

#
# Get list with just one plant
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$[0].id" == "{{1_my_plant_id}}"
jsonpath "$[0].plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$[0].plantSpeciesName" == "{{1_species_name}}"
jsonpath "$[0].adoptionDate" == "{{1_plant_adoption}}"
jsonpath "$[0].plantName" == "{{1_plant_name}}"

#
# Update first users plant
PUT http://localhost:8080/api/v1/my/plants/{{1_my_plant_id}}
Authorization: Bearer {{craig_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "adoptionDate": "{{1_plant_new_adoption}}",
  "plantName": "{{1_plant_new_name}}"
}
```
HTTP 204

#
# Get list with one plant, but updated info
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$[0].id" == "{{1_my_plant_id}}"
jsonpath "$[0].plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$[0].plantSpeciesName" == "{{1_species_name}}"
jsonpath "$[0].adoptionDate" == "{{1_plant_new_adoption}}"
jsonpath "$[0].plantName" == "{{1_plant_new_name}}"

#
# Create second users plant
POST http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "plantSpeciesID": "{{2_plant_species_id}}",
  "adoptionDate": "{{2_plant_adoption}}",
  "plantName": "{{2_plant_name}}"
}
```
HTTP 201
Content-Type: application/json; charset=utf-8
[Captures]
2_my_plant_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.plantSpeciesID" == "{{2_plant_species_id}}"
jsonpath "$.plantSpeciesName" == "{{2_species_name}}"
jsonpath "$.adoptionDate" == "{{2_plant_adoption}}"
jsonpath "$.plantName" == "{{2_plant_name}}"

#
# Get list with two plants
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$[0].id" == "{{2_my_plant_id}}"
jsonpath "$[0].plantSpeciesID" == "{{2_plant_species_id}}"
jsonpath "$[0].plantSpeciesName" == "{{2_species_name}}"
jsonpath "$[0].adoptionDate" == "{{2_plant_adoption}}"
jsonpath "$[0].plantName" == "{{2_plant_name}}"
jsonpath "$[1].id" == "{{1_my_plant_id}}"
jsonpath "$[1].plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$[1].plantSpeciesName" == "{{1_species_name}}"
jsonpath "$[1].adoptionDate" == "{{1_plant_new_adoption}}"
jsonpath "$[1].plantName" == "{{1_plant_new_name}}"

#
# Update second users plant
PUT http://localhost:8080/api/v1/my/plants/{{2_my_plant_id}}
Authorization: Bearer {{craig_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "adoptionDate": "{{2_plant_new_adoption}}",
  "plantName": "{{2_plant_new_name}}"
}
```
HTTP 204

#
# Get list with two plants, with second plant updated
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$[0].id" == "{{2_my_plant_id}}"
jsonpath "$[0].plantSpeciesID" == "{{2_plant_species_id}}"
jsonpath "$[0].plantSpeciesName" == "{{2_species_name}}"
jsonpath "$[0].adoptionDate" == "{{2_plant_new_adoption}}"
jsonpath "$[0].plantName" == "{{2_plant_new_name}}"
jsonpath "$[1].id" == "{{1_my_plant_id}}"
jsonpath "$[1].plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$[1].plantSpeciesName" == "{{1_species_name}}"
jsonpath "$[1].adoptionDate" == "{{1_plant_new_adoption}}"
jsonpath "$[1].plantName" == "{{1_plant_new_name}}"

#
# Delete plant 2
DELETE http://localhost:8080/api/v1/my/plants/{{2_my_plant_id}}
Authorization: Bearer {{craig_token}}
HTTP 204

#
# Get list with one plant, with second plant deleted
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$[0].id" == "{{1_my_plant_id}}"
jsonpath "$[0].plantSpeciesID" == "{{1_plant_species_id}}"
jsonpath "$[0].plantSpeciesName" == "{{1_species_name}}"
jsonpath "$[0].adoptionDate" == "{{1_plant_new_adoption}}"
jsonpath "$[0].plantName" == "{{1_plant_new_name}}"

#
# Delete plant 2 again and fail
DELETE http://localhost:8080/api/v1/my/plants/{{2_my_plant_id}}
Authorization: Bearer {{craig_token}}
HTTP 400

#
# Delete plant 1
DELETE http://localhost:8080/api/v1/my/plants/{{1_my_plant_id}}
Authorization: Bearer {{craig_token}}
HTTP 204

#
# Get list with no plants
GET http://localhost:8080/api/v1/my/plants
Authorization: Bearer {{craig_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Asserts]
jsonpath "$" == null
