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

# setup
# ========================================================================
# testing users

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
# Refresh token
POST http://localhost:8080/api/v1/auth/refresh
Authorization: Bearer {{craig_refresh_token}}
HTTP 200
Content-Type: application/json; charset=utf-8
[Captures]
craig_token: jsonpath "$.token"

#
# Revoke token
POST http://localhost:8080/api/v1/auth/revoke
Authorization: Bearer {{craig_refresh_token}}
Content-Type: application/json; charset=utf-8
```json
{
  "id": "{{craig_id}}"
}
```
HTTP 204

#
# Login back in to user account for new refresh token
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
jsonpath "$.isAdmin" == false
jsonpath "$.tokenExpiresAt" exists
jsonpath "$.refreshTokenExpiresAt" exists
