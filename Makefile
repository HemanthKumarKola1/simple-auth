DB_URL=postgresql://root:secret@localhost:5432/user_auth?sslmode=disable
name=init_schema
entershell:
	docker exec -it postgres /bin/sh

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root user_auth

dropdb:
	docker exec -it postgres dropdb user_auth

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

.PHONY: entershell postgres createdb dropdb migrateup migratedown new_migration sqlc
