.SILENT:
.PHONY: run migrate-Users migrate-Users-drop user-data-service-up

run:
	docker-compose up --remove-orphans --build

migrate-users:
	migrate -path ./internal/user_data/db/migrations -database 'postgres://postgres:admin@0.0.0.0:5433/users?sslmode=disable' up

migrate-users-drop:
	migrate -path ./internal/user_data/db/migrations -database 'postgres://postgres:adminy@0.0.0.0:5433/users?sslmode=disable' drop

user-data-service-up:
	docker-compose build && docker-compose up