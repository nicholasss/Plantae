# verify that server is online
GET http://localhost:8080/api/v1/health
HTTP 200
Content-Type: text/html; charset=UTF-8
[Asserts]
xpath "string(/html/body)" contains "OK"

# reset database
POST http://localhost:8080/api/v1/super-admin/reset-users
Authorization: SuperAdminToken {{super_admin_token}}

HTTP 204

# create account

