## Prerequisites

Before getting started, ensure you have done the following things:

1. Install `Go`
2. Install `Docker`
3. Register for a _Mailtrap_ account for sending test emails

## Getting Started

To start this repo on your machine, do the following:

1.  Clone this repo
2.  Go to the repo directory on your machine
3.  Execute `go mod tidy && go mod vendor`
4.  Setup the environment variables in `./.env` file (see below for details)
5.  Start `Docker` on your machine
6.  Run `docker compose up -d`
7.  Run `make migrate_up_all` for running the DB migrations
8.  Finally, `go run ./cmd`

## Environment Variables

Setup the following environment variables in `./.env` file:

```
ROOT_DIR="/Users/folder/to/project"
DB_USER="dbUser"
POSTGRES_USER=someusername
POSTGRES_PASSWORD="somepassword"
POSTGRES_MULTIPLE_DATABASES=main_dv,main_db_test
POSTGRES_URL="postgresql://someusername:somepassword@localhost:5433/fit_forge?sslmode=disable"
POSTGRES_TEST_DB_URL="postgresql://root:password@localhost:5433/fit_forge_test?sslmode=disable"
AUTH_SECRET_KEY="sup3rs3cr3t"
REDIS_DSN="localhost:6379"
REDIS_PASSWORD="secret"
RABBITMQ_DEFAULT_USER="username"
RABBITMQ_DEFAULT_PASS="password"
RABBITMQ_URL="amqp://username:password@localhost:5672/"
EMAIL_HOST="https://sandbox.api.mailtrap.io"
MAILTRAP_API_KEY="API_KEY"
SUBSCRIPTION_INTERVAL_MINUTES=1
PUBSUB_PROJECT_ID=local-project
PUBSUB_EMULATOR_HOST=localhost:8085
```
