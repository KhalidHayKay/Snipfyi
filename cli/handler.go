package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CLIHandler struct {
	pgsql *pgxpool.Pool
}

func NewCLIHandler(pgsql *pgxpool.Pool) *CLIHandler {
	return &CLIHandler{pgsql}
}

func (h *CLIHandler) migrateUp(ctx context.Context) {
	_, err := h.pgsql.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("Schema migration table creation failed: %s", err)
	}

	var isUpToDate bool = true

	for _, m := range migrations {
		var exists bool

		err := h.pgsql.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name=$1)`,
			m.Name,
		).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			continue
		}

		_, err = h.pgsql.Exec(ctx, m.Up)
		if err != nil {
			log.Fatalf("Migration failed: %s\n%v", m.Name, err)
		}

		_, err = h.pgsql.Exec(ctx,
			`INSERT INTO schema_migrations (name) VALUES ($1)`,
			m.Name,
		)
		if err != nil {
			log.Fatal(err)
		}

		isUpToDate = false

		log.Printf("Applied migration: %s", m.Name)
	}

	if isUpToDate {
		log.Println("Migration up to date")
	}
}

func (h *CLIHandler) migrateDown(ctx context.Context) {
	schemas, err := h.getSchemas(ctx)
	if err != nil {
		log.Fatalf("Schema migrations fetch error: %s", err)
	}

	if len(schemas) < 1 {
		log.Println("No migrations to rollback.")
		return
	}

	var migrationMap = func() map[string]Migration {
		m := make(map[string]Migration)
		for _, mig := range migrations {
			m[mig.Name] = mig
		}
		return m
	}()

	for _, s := range schemas {
		target, ok := migrationMap[s.Name]

		if !ok {
			log.Fatalf("Migration %s not found in code", target.Name)
		}

		// run DOWN
		_, err = h.pgsql.Exec(ctx, target.Down)
		if err != nil {
			log.Fatalf("Rollback failed: %s\n%v", target.Name, err)
		}

		// remove from schema_migrations
		_, err = h.pgsql.Exec(ctx,
			`DELETE FROM schema_migrations WHERE name=$1`,
			target.Name,
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Rolled back migration: %s", target.Name)
	}
}

func (h *CLIHandler) migrationStatus(ctx context.Context) {
	schemas, err := h.getSchemas(ctx)
	if err != nil {
		log.Fatalf("Schema migrations fetch err: %s", err)
	}

	appliedMap := make(map[string]Schema)
	for _, s := range schemas {
		appliedMap[s.Name] = s
	}

	for _, m := range migrations {
		if _, exists := appliedMap[m.Name]; exists {
			log.Printf("✅ %s..........APPLIED", m.Name)
		} else {
			log.Printf("⏳ %s..........PENDING", m.Name)
		}
	}
}

func (h *CLIHandler) resetDB(ctx context.Context) {
	log.Println("Migrating down...")
	h.migrateDown(ctx)
	log.Println("Migrating up...")
	h.migrateUp(ctx)
}

func (h *CLIHandler) seedDB() {
	fmt.Println("No seed logic implemented yet.")
}

func (h *CLIHandler) getSchemas(ctx context.Context) (schemas []Schema, err error) {
	rows, err := h.pgsql.Query(ctx, `
		SELECT id, name, applied_at
		FROM schema_migrations
		ORDER BY id DESC
		LIMIT 5
	`)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var schema Schema
		if err = rows.Scan(
			&schema.Id,
			&schema.Name,
			&schema.AppliedAt,
		); err != nil {
			return
		}
		schemas = append(schemas, schema)
	}

	err = rows.Err()

	return
}
