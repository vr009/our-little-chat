.SILENT:
.PHONY: run migrate-Users migrate-Users-drop user-data-service-up start
INDEX_PATH ?= temp/dist/index.html
BRANCH := dev
USERS_DB_URL := 'postgres://postgres:admin@0.0.0.0:5433/users?sslmode=disable'
USERS_MIGRATIONS_PATH := ./internal/user_data/db/migrations
MOCKS_DESTINATION=internal/mocks

${INDEX_PATH} frontend:
	mkdir temp && cd temp && git clone https://github.com/vr009/our_little_chatik_frontend.git --recursive &&\
	cd our_little_chatik_frontend && git checkout ${BRANCH} && yarn install && yarn build

run:
	docker-compose -f docker-compose-test.yml up --remove-orphans --build

start: frontend run
	echo "Starting project..."

migrate-users:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} up

migrate-users-drop:
	migrate -path ${USERS_MIGRATIONS_PATH} -database ${USERS_DB_URL} drop

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

.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/users/internal/interfaces.go internal/chat/internal/interfaces.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$${file#*/}; done

swagger:
	swag init -g internal/users/cmd/main.go --output docs/

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	users.proto
