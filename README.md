# Welcome to simple-auth!
This is a golang application for simple authentication of users using postgres and redis.

Steps to start the application:
1. Clone the repo.
2. Enter `docker compose up` when in the project directory.

## API Endpoints

1. signup: Signup to the application using following cURL

> curl --location 'http://localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data '{
    "username": "maybe123",
    "password": "hellohemanth"
}'

2. login: Login to the application to get a jwt

> curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "maybe123",
    "password": "hellohemanth"
}'

4. revoke: use the fetched jwt to revoke

> curl --location --request POST 'http://localhost:8080/revoke' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2NjYyNTUsInVzZXJuYW1lIjoibWF5YmUxMjMifQ.LNAuJ33vA2l4W4V1IDZnK7xGy_WRirM1ptLyKoFgzi0'

5. refresh: you will not be able to refresh with revoked jwt, use a jwt without revoking to refresh

> curl --location 'http://localhost:8080/refresh' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2NjYyNTUsInVzZXJuYW1lIjoibWF5YmUxMjMifQ.LNAuJ33vA2l4W4V1IDZnK7xGy_WRirM1ptLyKoFgzi0'

6. Test api to test the authorization

> curl --location 'http://localhost:8080/test' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2NzE2MzcsInVzZXJuYW1lIjoibWF5YmUxMjMifQ.A9p3Ba8PyLhRETg7r4INKyt_iLKXJXdiAzy-g5wSlNo'

Note: Replace authorization header with valid jwts.