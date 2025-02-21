include .env
.PHONY: clean critic security lint test build run

APP_NAME = apiserver
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/database/migrations
DB_CONN=postgres://$(DB_USER:"%"=%):$(DB_PASSWORD:"%"=%)@$(DB_HOST:"%"=%):$(DB_PORT:"%"=%)/$(DB_NAME:"%"=%)?sslmode=$(DB_SSL_MODE:"%"=%)

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

migrate_create:
	@read -p "migration name (do not use space): " NAME \
  	&& migrate create -ext sql -dir $(MIGRATIONS_FOLDER) $${NAME}

migrate_up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_CONN)" up

migrate_down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_CONN)" down

# migrate_force:
# 	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

migrate_force:
	@read -p "please enter the migration version (the migration filename prefix): " VERSION \
  	&& migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_CONN)" force $${VERSION}

migrate_version:
	@migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_CONN)" version 

migrate_drop:
	@migrate -path $(MIGRATIONS_FOLDER) -database "$(DB_CONN)" drop

generate:
	go generate ./...

# development
dev: generate
	go run github.com/cosmtrek/air

install:
	cd .. && go install github.com/golang-migrate/migrate/v4 && \
	go install -u github.com/swaggo/swag/cmd/swag && go install -u github.com/cosmtrek/air && \
	go install github.com/vektra/mockery/v2/.../ && \
	cd ${APP_NAME} && swag init