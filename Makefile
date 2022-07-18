up:
	docker compose up --build

test:
	API_ENV=test go test ./...

mock:
	mockgen -source=db/auth_repository.go -destination=db/mocks/db_mock.go -package=mocks