build:
	go build -o bin/main cmd/app/main.go

run:
	ENV=development nodemon --exec go run cmd/app/main.go --signal SIGTERM

docker:
	docker-compose up -d