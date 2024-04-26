include .env

validate_env:
	@if [ -z "${DB_USER}" ]; then echo "Error: DB_USER is not defined in .env"; exit 1; fi
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

migrate_up: validate_env
	migrate -path migrations -database "postgresql://${DB_USER}:${POSTGRES_PASSWORD}@localhost:5433/fit_forge?sslmode=disable" -verbose up

migrate_down: validate_env
	migrate -path migrations -database "postgresql://${DB_USER}:${POSTGRES_PASSWORD}@localhost:5433/fit_forge?sslmode=disable" -verbose down

.PHONY: run_db_container start_db stop_db create_db drop_db migrate_up migrate_down 
