up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/medication_history_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationHistoryService
	 mockgen -destination=mocks/medication_history_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationHistoryRepository
	 mockgen -destination=mocks/mailer_mock.go -package=mocks github.com/decagonhq/meddle-api/services Mailer
	 mockgen -destination=mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
	 mockgen -destination=mocks/auth_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db AuthRepository
	 mockgen -destination=mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService
	 mockgen -destination=mocks/push_notification.go -package=mocks github.com/decagonhq/meddle-api/services PushNotifier
	 mockgen -destination=mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository


test: generate-mock
	 MEDDLE_ENV=test go test ./...
