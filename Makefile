up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/mailer_mock.go -package=mocks github.com/decagonhq/meddle-api/services Mailer
	 mockgen -destination=mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
	 mockgen -destination=mocks/auth_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db AuthRepository
	 mockgen -destination=mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService
	 mockgen -destination=mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

test: generate-mock
	 MEDDLE_ENV=test go test ./...
# API_ENV=test go test ./...
