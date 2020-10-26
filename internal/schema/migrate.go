package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
CREATE TABLE "users" (
	"user_id" UUID PRIMARY KEY,
	"name" varchar NOT NULL,
	"email" varchar UNIQUE NOT NULL,
	"roles" varchar[] NOT NULL,
	"password_hash" bytea NOT NULL,
	"date_created" timestamp NOT NULL,
	"date_updated" timestamp NOT NULL
);`,
	},
	{
		Version:     2,
		Description: "Add requests",
		Script: `
CREATE TABLE prime_number_requests (
	request_id UUID PRIMARY KEY,
	user_id UUID,
	send_number BIGINT,
	receive_number BIGINT,
	date_created TIMESTAMP,

	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {

	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
