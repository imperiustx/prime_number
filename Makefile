db:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it postgres12 psql -U root

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root prime_number

dropdb:
	docker exec -it postgres12 dropdb prime_number

.PHONY: db psql createdb dropdb
