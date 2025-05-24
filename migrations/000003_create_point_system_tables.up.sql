CREATE TABLE "points" (
  "user_id" integer PRIMARY KEY REFERENCES "users" ("id") ON DELETE CASCADE,
  "total_points" integer DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "point_transactions" (
  "id" uuid PRIMARY KEY,
  "transaction_type" varchar(50) CHECK (transaction_type IN ('earn', 'spend', 'expire', 'subscription_deduction')),
  "points" integer NOT NULL,
  "reason" text,
  "user_id" integer REFERENCES "users" ("id") ON DELETE CASCADE,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
