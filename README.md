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
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MjUwNTAsInVzZXJuYW1lIjoibWF5YmUxMjMifQ.r_-JnCua20FhO7L8jN9jXesuI0FPB0ySbZ7EHzY7s_w'

5. refresh: you will not be able to refresh with revoked jwt, use a jwt without revoking to refresh

> curl --location 'http://localhost:8080/refresh' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM2MjUwNTAsInVzZXJuYW1lIjoibWF5YmUxMjMifQ.r_-JnCua20FhO7L8jN9jXesuI0FPB0ySbZ7EHzY7s_w'
