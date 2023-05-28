.SILENT:
.PHONY: run migrate-Users migrate-Users-drop user-data-service-up start
INDEX_PATH ?= temp/dist/index.html
BRANCH := main
USERS_DB_URL := 'postgres://postgres:admin@0.0.0.0:5433/users?sslmode=disable'
USERS_MIGRATIONS_PATH := ./internal/user_data/db/migrations

${INDEX_PATH} frontend:
	mkdir temp && cd temp && git clone https://github.com/vr009/our_little_chatik_frontend.git --recursive &&\
	cd our_little_chatik_frontend && git checkout ${BRANCH} && yarn install && yarn build && cp -r dist ..

start: frontend
	docker-compose build && docker-compose up -d && docker ps

run:
	docker-compose up --remove-orphans --build

migrate-users:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} up

migrate-users-drop:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} drop

integration-user-data-test:
	docker-compose -f docker-compose-test.yml down &&\
	docker-compose -f docker-compose-test.yml up -d test-db-user-data && sleep 1 &&\
	docker-compose -f docker-compose-test.yml up -d test-user-data &&\
	go clean -testcache && TEST_HOST=http://localhost:8086 go test ./internal/user_data/cmd/... &&\
	TEST_HOST=http://localhost:8086 go test -bench=. -run=^Benchmark -benchmem -benchtime=100x ./internal/user_data/cmd/... &&\
	docker-compose -f docker-compose-test.yml down

integration-test: integration-user-data-test
