.SILENT:
.PHONY: run migrate-Users migrate-Users-drop user-data-service-up start
INDEX_PATH ?= temp/dist/index.html
BRANCH := main
USERS_DB_URL := 'postgres://postgres:admin@0.0.0.0:5433/users?sslmode=disable'
USERS_MIGRATIONS_PATH := ./internal/user_data/db/migrations

${INDEX_PATH} frontend:
	mkdir temp && cd temp && git clone https://github.com/vr009/our_little_chatik_frontend.git --recursive &&\
	cd our_little_chatik_frontend && git checkout ${BRANCH} && yarn install && yarn build && cp -r dist ..

run:
	docker-compose -f docker-compose-test.yml up --remove-orphans --build

start: frontend run
	echo "Starting project..."

migrate-users:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} up

migrate-users-drop:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} drop

.PHONY: integration-peer-test
integration-peer-test:
	docker-compose -f docker-compose-test.yml down &&\
	docker ps && \
	docker-compose -f docker-compose-test.yml up -d test-db-peer && sleep 1 &&\
	docker-compose -f docker-compose-test.yml up -d --build test-peer &&\
	go clean -testcache && JWT_SIGNED_KEY=test PEER_HOST=localhost PEER_PORT=8089 \
	go test ./internal/peer/cmd/... &&\
	docker-compose -f docker-compose-test.yml down

integration-user-data-test:
	docker-compose -f docker-compose-test.yml down &&\
	docker-compose -f docker-compose-test.yml up -d test-db-user-data && sleep 1 &&\
	docker-compose -f docker-compose-test.yml up -d --build test-user-data &&\
	go clean -testcache && TEST_HOST=http://localhost:8086 go test ./internal/user_data/cmd/... &&\
	TEST_HOST=http://localhost:8086 go test -bench=. -run=^Benchmark -benchmem -benchtime=100x ./internal/user_data/cmd/... &&\
	docker-compose -f docker-compose-test.yml down

.PHONY: integration-chat-test
integration-chat-test:
	docker-compose -f docker-compose-test.yml down &&\
	docker-compose -f docker-compose-test.yml up -d test-db-peer test-db-chat && sleep 1 &&\
	docker-compose -f docker-compose-test.yml up -d --build test-chat &&\
	go clean -testcache && TEST_HOST=http://localhost:8083 JWT_SIGNED_KEY=test \
	DATABASE_URL="user=test password=test host=localhost port=5433 dbname=chats" \
	REDIS_HOST="localhost" REDIS_PORT="6379" REDIS_PASSWORD="test" go test ./internal/chat/cmd/... &&\
	docker-compose -f docker-compose-test.yml down

## test: run all integration tests
.PHONY: integration
integration: integration-user-data-test integration-chat-test integration-peer-test

## test: run all unit tests
.PHONY: unit
unit:
	go test -v -race -short ./... -coverprofile coverage.out -covermode atomic && go tool cover -func coverage.out

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## dep: download dependencies
.PHONY: dep
dep:
	go mod download

## lint: run linter
.PHONY: lint
lint:
	golangci-lint run -v -c golangci-lint.yml
