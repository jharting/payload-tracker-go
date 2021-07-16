build-api:
	
	go build -o payload-tracker-api cmd/payload-tracker-api/main.go

build-consumer:

	go build -o payload-tracker-consumer cmd/payload-tracker-consumer/main.go

build-all:

	go build -o payload-tracker-api cmd/payload-tracker-api/main.go
	go build -o payload-tracker-consumer cmd/payload-tracker-consumer/main.go
