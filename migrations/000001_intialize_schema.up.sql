CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "password" char(60) NOT NULL,
  "email_verified_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "is_completed" bool NOT NULL DEFAULT false,
  "user_id" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "todos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
