-- simply-bank schema (PostgreSQL 18+)
-- PKs are UUIDv7: time-orderable, globally unique, server-generated.
-- Money is numeric (arbitrary precision), mapped to decimal.Decimal in Go.

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT uuidv7(),
  "owner" varchar NOT NULL,
  "balance" numeric NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" uuid PRIMARY KEY DEFAULT uuidv7(),
  "account_id" uuid NOT NULL,
  "amount" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" uuid PRIMARY KEY DEFAULT uuidv7(),
  "from_account_id" uuid NOT NULL,
  "to_account_id" uuid NOT NULL,
  "amount" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE RESTRICT;

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE RESTRICT;

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE RESTRICT;
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE RESTRICT;

CREATE INDEX ON "accounts" ("owner");
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
