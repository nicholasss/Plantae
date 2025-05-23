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
POST http://localhost:8080/api/v1/create-user
Content-Type: application/json; charset=UTF-8
```json
{
  "client": "hurl-project-test",
  "email": "{{craig_email}}",
  "password": "{{craig_password}}"
}
```
HTTP 201
Content-Type: application/json
[Captures]
craig_id: jsonpath "$.id"
[Asserts]
jsonpath "$.isAdmin" == false
jsonpath "$.email" == "{{craig_email}}"

#
# Login to user account
POST http://localhost:8080/api/v1/login-user
Content-Type: application/json; charset=UTF-8
```json
{
  "client": "hurl-project-test",
  "email": "{{craig_email}}",
  "password": "{{craig_password}}"
}
```
HTTP 200
Content-Type: application/json
[Captures]
craig_token: jsonpath "$.token"
craig_refresh_token: jsonpath "$.refreshToken"
[Asserts]
jsonpath "$.isAdmin" == false
