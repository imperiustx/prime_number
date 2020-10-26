CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "roles" varchar[] NOT NULL,
  "password_hash" varchar NOT NULL,
  "date_created" timestamp NOT NULL DEFAULT (now()),
  "date_updated" timestamp
);

CREATE TABLE "prime_number_requests" (
  "request_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "send_number" bigint NOT NULL,
  "receive_number" bigint NOT NULL,
  "date_created" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "prime_number_requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "prime_number_requests" ("user_id");
