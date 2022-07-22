up:
	docker compose up --build

test:
	API_ENV=test go test ./...

mock:
	mockgen -source=services/auth.go -destination=db/mocks/auth_mock.go -package=mocks