DB_URL=postgresql://root:secret@localhost:5432/users?sslmode=disable
name=init_schema

entershell:
	docker exec -it postgres /bin/sh

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

createdb:
	docker exec -it simple-auth-db-1 createdb --username=root --owner=root users

dropdb:
	docker exec -it simple-auth-db-1 dropdb users

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

.PHONY: entershell postgres createdb dropdb migrateup migratedown new_migration sqlc
