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
# Add a plant to the plant species table
POST http://localhost:8080/api/v1/admin/plants
Authorization: Bearer {{lisa_token}}
```json
{
  "client": "hurl-project-test",
  "species": "{{plant_species}}",
  "humanPoisonToxic": {human_poison_toxic},
  "petPoisonToxic": {{pet_poison_toxic}},
  "humanEdible": {{human_edible}},
  "petEdible": {{pet_edible}}
}
```
HTTP 204
