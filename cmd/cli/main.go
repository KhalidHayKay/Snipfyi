package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"smply/config"
	"smply/internal/storage"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected command: db:drop | db:migrate | db:seed")
	}

	config.LoadEnv()
	if err := storage.InitDB(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "db:drop":
		clearDB()
	case "db:migrate":
		migrateDB(ctx)
	case "db:seed":
		seedDB()
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func seedDB() {
	fmt.Println("No seed logic implemented yet.")
	// your seed logic here
}

func migrateDB(ctx context.Context) {
	_, err := storage.DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			original TEXT NOT NULL,
			alias TEXT UNIQUE NOT NULL,
			visited INTEGER DEFAULT 0,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_visited TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS api_keys (
			id BIGSERIAL PRIMARY KEY,
			owner_email TEXT NOT NULL UNIQUE,
			key_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,
			last_used_at TIMESTAMP NULL,
			revoked_at TIMESTAMP NULL
		);

		CREATE TABLE IF NOT EXISTS magic_tokens (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL,
			token_hash TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			used_at TIMESTAMP NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS click_events	(
			id BIGSERIAL PRIMARY KEY,
			url_id INTEGER NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
			referer TEXT,
			user_agent TEXT,
			ip_address TEXT,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database tables created successfully.")
}

func clearDB() {
	fmt.Print("CRITICAL: This will delete ALL data. Type 'yes' to proceed: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	if strings.TrimSpace(response) != "yes" {
		fmt.Println("Aborted.")
		return
	}

	ctx := context.Background()

	tables := []string{"urls", "api_keys", "magic_tokens"}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)

		_, err := storage.DB.Exec(ctx, query)
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("Successfully DROPPED table %s", table)
	}
}
