up:
	docker compose up --build

test:
	API_ENV=test go test ./...