up:
	docker compose up --build

generate-mock:
	go generate ./...

test: generate-mock
	API_ENV=test go test ./...