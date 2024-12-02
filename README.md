## Prerequisites

Before getting started, ensure you have done the following things:

1. Install `Go` on your machine
2. Install `Postgres` and `Redis` either directly on your machine or via Docker
3. Register for a _Mailtrap_ account for sending test emails

## Getting Started

To start this repo on your machine, do the following:

1.  Clone this repo
2.  Go to the repo directory on your machine
3.  Execute `go mod tidy && go mod vendor`
4.  Setup the environment variables in `./.env` file (see below for details)
5.  Start `Postgres` and `Redis` services
6.  Run `create_db` for creating the Postgres DB
7.  Run `make migrate_up_all` for running the DB migrations
8.  Finally, `go run ./cmd`

## Environment Variables

Setup the following environment variables in `./.env` file:

```
DB_USER="dbUser"
POSTGRES_PASSWORD="somepassword"
POSTGRES_URL="postgresql://root:secret@localhost:5433/fit_forge?sslmode=disable"
AUTH_SECRET_KEY="sup3rs3cr3t"
REDIS_DSN="localhost:6379"
EMAIL_HOST="https://sandbox.api.mailtrap.io"
MAILTRAP_API_KEY="API_KEY"
```
