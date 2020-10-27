db:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it postgres12 psql -U root

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root prime_number

dropdb:
	docker exec -it postgres12 dropdb prime_number

run:
	go run ./cmd/server-api

lint:
	golangci-lint run

migrate:
	go run ./cmd/server-admin/ migrate

seed:
	go run ./cmd/server-admin/ seed

expvarmon:
	expvarmon -ports="6060" -endpoint="/debug/vars" -vars="requests,goroutines,errors,mem:memstats.Alloc"

.PHONY: db psql createdb dropdb run lint migrate seed expvarmon
