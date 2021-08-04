build-api:
	
	go build -o pt-api cmd/payload-tracker-api/main.go

build-consumer:

	go build -o pt-consumer cmd/payload-tracker-consumer/main.go

build-all:

	go build -o pt-api cmd/payload-tracker-api/main.go
	go build -o pt-consumer cmd/payload-tracker-consumer/main.go

lint:
	gofmt -l .
	gofmt -s -w .

seed:
	go build cmd/db-seeder/main.go
	go run cmd/db-seeder/main.go
