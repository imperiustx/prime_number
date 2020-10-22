CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "roles" varchar[] NOT NULL,
  "password_hash" bytea NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "requests" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "send_number" int NOT NULL,
  "receive_number" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "requests" ("user_id");
