include .env

validate_env:
	@if [ -z "${POSTGRES_PASSWORD}" ]; then echo "Error: POSTGRES_PASSWORD is not defined in .env"; exit 1; fi

run_db_container: validate_env
	docker run --name postgres12 -p 5433:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

start_db:
	docker start postgres12

stop_db:
	docker stop postgres12

create_db: validate_env
	docker exec -it postgres12 createdb --username=${DB_USER} --owner=${DB_USER} fit_forge

drop_db:
	docker exec -it postgres12 dropdb fit_forge

migrate_up_all: validate_env
	migrate -path migrations -database "${POSTGRES_URL}" -verbose up

migrate_down_all: validate_env
	migrate -path migrations -database "postgresql://${DB_USER}:${POSTGRES_PASSWORD}@localhost:5433/fit_forge?sslmode=disable" -verbose down

migrate_up_1: validate_env
	migrate -path migrations -database "postgresql://${DB_USER}:${POSTGRES_PASSWORD}@localhost:5433/fit_forge?sslmode=disable" -verbose up 1

migrate_down_1: validate_env
	migrate -path migrations -database "postgresql://${DB_USER}:${POSTGRES_PASSWORD}@localhost:5433/fit_forge?sslmode=disable" -verbose down 1

test:
	sh -c 'env $$(cat .env | xargs) go test ./...'

run:
	sh -c 'env $$(cat .env.prod | xargs) go run ./cmd'

unit_test_coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: run_db_container start_db stop_db create_db drop_db migrate_up_all migrate_down_all migrate_up_1 migrate_down_1 test unit_test_coverage
