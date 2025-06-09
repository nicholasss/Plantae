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
# Create user account
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json; charset=UTF-8
```json
{
  "client": "hurl-project-test",
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
  "client": "hurl-project-test",
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

#
# Test empty plant species table
GET http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
HTTP 204

#
# Add a plant to the plant_species table
# -- with all details
POST http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
```json
{
  "client": "hurl-project-test",
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
GET http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
HTTP 200
[Captures]
1_plant_species_id: jsonpath "$[0].id"
[Asserts]
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == "hurl-project-test"
jsonpath "$[0].updatedBy" == "hurl-project-test"
jsonpath "$[0].speciesName" == "{{1_species_name}}"
jsonpath "$[0].humanPoisonToxic" == false
jsonpath "$[0].petPoisonToxic" == false
jsonpath "$[0].humanEdible" == false
jsonpath "$[0].petEdible" == false

#
# Add a second plant to the plant_species
# -- with required details
POST http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
```json
{
  "client": "hurl-project-test",
  "speciesName": "{{2_species_name}}"
}
```
HTTP 204

#
# Requesting list of plant species in the table
# -- with two records in table
GET http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
HTTP 200
[Captures]
2_plant_species_id: jsonpath "$[0].id"
[Asserts]
jsonpath "$[0].speciesName" == "{{2_species_name}}"
jsonpath "$[0].createdAt" isIsoDate
jsonpath "$[0].updatedAt" isIsoDate
jsonpath "$[0].createdBy" == "hurl-project-test"
jsonpath "$[0].updatedBy" == "hurl-project-test"

#
# Updating plant with new information
POST http://localhost:8080/api/v1/admin/plants/{{2_plant_species_id}}
Authorization: Bearer {{lisa_token}}
```json
{
  "client": "hurl-project-test",
  "speciesName": "{{2_species_name}}",
  "humanPoisonToxic": {{2_human_poison_toxic}},
  "petPoisonToxic": {{2_pet_poison_toxic}},
  "humanEdible": {{2_human_edible}},
  "petEdible": {{2_pet_edible}}
}
```
HTTP 204
