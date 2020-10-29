db:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it postgres12 psql -U root

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root prime_number

dropdb:
	docker exec -it postgres12 dropdb prime_number

run:
	./bin/server
	# go run ./cmd/server-api

lint:
	golangci-lint run

migrate:
	./bin/admin migrate
	# go run ./cmd/server-admin/ migrate

seed:
	./bin/admin seed
	# go run ./cmd/server-admin/ seed

expvarmon:
	expvarmon -ports="6060" -endpoint="/debug/vars" -vars="requests,goroutines,errors,mem:memstats.Alloc"

private:
	./bin/admin keygen private.pem
	# go run ./cmd/server-admin keygen private.pem

build:
	docker image build -t server .

.PHONY: db psql createdb dropdb run lint migrate seed expvarmon private build
