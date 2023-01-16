.SILENT:
.PHONY: run migrate-Users migrate-Users-drop user-data-service-up start
INDEX_PATH ?= temp/dist/index.html

${INDEX_PATH} frontend:
	mkdir temp && cd temp && git clone https://github.com/vr009/our_little_chatik_frontend.git --recursive &&\
	cd our_little_chatik_frontend && npm install vite && npm run build && cp -r dist ..

start: frontend
	docker-compose build && docker-compose up -d && docker ps

run:
	docker-compose up --remove-orphans --build

migrate-users:
	migrate -path ./internal/user_data/db/migrations -database 'postgres://postgres:admin@0.0.0.0:5433/users?sslmode=disable' up

migrate-users-drop:
	migrate -path ./internal/user_data/db/migrations -database 'postgres://postgres:adminy@0.0.0.0:5433/users?sslmode=disable' drop

user-data-service-up:
	docker-compose build && docker-compose up