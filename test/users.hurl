GET http://localhost:8080/health
HTTP 200
Content-Type: text/html; charset=UTF-8
[Asserts]
xpath "string(/html/body)" contains "OK"
