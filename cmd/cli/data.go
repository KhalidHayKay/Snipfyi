package main

var commandMap = map[string]Command{
	"db:migrate up": {
		Run:         migrateUp,
		Destructive: false,
	},
	"db:migrate down": {
		Run:         migrateDown,
		Destructive: true,
	},
	"db:migrate status": {
		Run:         migrationStatus,
		Destructive: false,
	},
	"db:reset": {
		Run:         resetDB,
		Destructive: true,
	},
}

var migrations = []Migration{
	{
		Name: "000001_create_urls_table",
		Up: `CREATE TABLE urls (
			id SERIAL PRIMARY KEY,
			original TEXT NOT NULL,
			alias TEXT UNIQUE NOT NULL,
			visited INTEGER DEFAULT 0,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_visited TIMESTAMP
		);`,
		Down: `DROP TABLE IF EXISTS urls CASCADE;`,
	},
	{
		Name: "000002_create_api_keys_table",
		Up: `CREATE TABLE api_keys (
			id BIGSERIAL PRIMARY KEY,
			owner_email TEXT NOT NULL UNIQUE,
			key_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			last_used_at TIMESTAMP NULL,
			revoked_at TIMESTAMP NULL
		);`,
		Down: `DROP TABLE IF EXISTS api_keys CASCADE;`,
	},
	{
		Name: "000003_create_magic_tokens_table",
		Up: `CREATE TABLE magic_tokens (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			used_at TIMESTAMP NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		Down: `DROP TABLE IF EXISTS magic_tokens CASCADE;`,
	},
	{
		Name: "000004_create_click_events_table",
		Up: `CREATE TABLE click_events (
			id BIGSERIAL PRIMARY KEY,
			url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
			referer TEXT,
			user_agent TEXT,
			ip_address TEXT,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		Down: `DROP TABLE IF EXISTS click_events CASCADE;`,
	},
}
