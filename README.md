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
8.  Run `make generate_jwks` to generate the necessary secrets
9.  Finally, `make run` to run the Go app with `.env` file

## Manual Deployment to GCP

This project is deployed to **GCP Artifact Registry**, **Cloud Run** and **API Gateway** using a combination of `gcloud` CLI and GCP Console. Before start deploying the app, do the following step:

- Authenticate a user using `gcloud auth application-default login --no-launch-browser` (using `--no-launch-browser` somehow that works)

### 1. Create a new repository on Artifact Registry

1. Login to GCP Console (https://console.cloud.google.com/)
2. Select "Artifact Registry" menu
3. Click "+ Create repository" button
4. Name a repository (e.g. `my-app-prod-repo`)
5. Ensure the format is `Docker` and use `us-central1` as the region or other regions that are within the Free tier
6. Click 'Create'

### 2. Dockerize the app and push the image to Artifact Registry

1. Go back to the terminal and dockerize the app using `docker buildx build -t us-central1-docker.pkg.dev/{project-id}/{repository-name}/{app-name}:1.0.0 .` command (note that the command ends with `.`)
2. Once the Docker image is created, we push it using `docker push us-central1-docker.pkg.dev/{project-id}/{repository-name}/{app-name}:1.0.0` command

_Note: `project-id` is the GCP project ID, `repository-name` is what you created on Step 1.4, `app-name` is the name of the app (can be anything)_

### 3. Create service and deploy Cloud Run

1. Go back to GCP Console and select or find "Cloud Run" menu
2. Click "Services" > "+ Deploy Container"
3. Choose "Artifact Registry | Docker Hub"
4. Select the Docker image you just pushed as the "Container image URL"
5. Pick a name for the Cloud Run service (e.g. `core-api`)
6. Choose `us-central1` as the region or any other region as long as it's the same region you choose for the Artifact Registry's repository
7. Setup the corresponding environment variables on "Containers, Volumes, Networking, Security" > "Containers" > "Variables & Secrets"
8. Click "Create"

### 4. Create an API Gateway

1. Create an OpenAPI specification named `openapi.yaml` in the repository (all of routes and their HTTP methods need to be defined in API Gateway)
2. Go to GCP Console and select "API Gateway" menu
3. Create an API Gateway using the `openapi.yaml` file

That's the gist of deploying this repository on GCP.

## Environment Variables

Setup the following environment variables in `./.env` file:

```
ENV=localhost
ROOT_DIR=/Users/folder/to/project
DB_USER=dbUser
POSTGRES_USER=someusername
POSTGRES_PASSWORD=somepassword
POSTGRES_MULTIPLE_DATABASES=main_dv,main_db_test
POSTGRES_URL=postgresql://someusername:somepassword@localhost:5433/fit_forge?sslmode=disable
POSTGRES_TEST_DB_URL=postgresql://root:password@localhost:5433/fit_forge_test?sslmode=disable
AUTH_SECRET_KEY=sup3rs3cr3t
REDIS_DSN=localhost:6379
REDIS_USERNAME=username
REDIS_PASSWORD=secret
EMAIL_HOST=https://sandbox.api.mailtrap.io
MAILTRAP_API_KEY=API_KEY
SUBSCRIPTION_INTERVAL_MINUTES=1
PUBSUB_PROJECT_ID=local-project
PUBSUB_EMULATOR_HOST=localhost:8085
JWT_ISSUER_CLAIM=your-app-url
JWT_AUDIENCE_CLAIM=your-app-audience
JWK_KEY_ID=your-app-jwk-key-id
GCP_SECRET_DIR=your-secrets-dir-on-gcp
FRONTEND_URL=http://localhost:3000
```
