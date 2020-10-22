db:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it postgres12 psql -U root

migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root prime_number

dropdb:
	docker exec -it postgres12 dropdb prime_number

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/prime_number?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/prime_number?sslmode=disable" -verbose down


.PHONY: db psql migrateinit createdb dropdb migrateup migratedown
