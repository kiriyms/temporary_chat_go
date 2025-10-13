run: build
	@./bin/app

build:
	@go build -o bin/app ./cmd

push:
	@git push origin main

docker.build:
	@docker build -t tempochat:local .

docker.run:
	@docker run --rm -p 8080:1323 --env-file .env tempochat:local