package main

import (
	"context"
	"smply/internal/storage"
)

func getSchemas(ctx context.Context) (schemas []Schema, err error) {
	rows, err := storage.DB.Query(ctx, `
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
