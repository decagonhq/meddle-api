up:
	docker compose up --build

test:
	API_ENV=test go test ./...

mock:
	mockgen -source=services/auth_service.go -destination=db/mocks/auth_service_mock.go -package=mocks
