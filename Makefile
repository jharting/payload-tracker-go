all: build-all

build-all: pt-api pt-consumer pt-migration

pt-api:
	go build -o pt-api cmd/payload-tracker-api/main.go

pt-consumer:
	go build -o pt-consumer cmd/payload-tracker-consumer/main.go

pt-migration:
	go build -o pt-migration internal/migration/main.go

lint:
	gofmt -l .
	gofmt -s -w .

test:
	go test -p 1 -v ./...

pt-seeder:
	go build -o pt-seeder tools/db-seeder/main.go

run-seed: pt-seeder
	./pt-seeder

run-migration: pt-migration
	./pt-migration

run-vaccum:
	CLOWDER_DATABASE_PASSWORD=crc CLOWDER_DATABASE_HOSTNAME=127.0.0.1 CLOWDER_DATABASE_USERNAME=crc CLOWDER_DATABASE_NAME=crc ./tools/cron-job.sh

clean:
	go clean
	rm -f pt-api
	rm -f pt-consumer
	rm -f pt-migration
