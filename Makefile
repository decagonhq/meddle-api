up:
	docker compose up --build

test:
	generate-mock
	API_ENV=test go test ./...

generate-mock:
	go generate ./...