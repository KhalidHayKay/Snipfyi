package main

import (
	"context"
	"log"
	"smply/config"
	"smply/internal/storage"
)

func main() {
	config.LoadEnv()

	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	handler := NewCLIHandler(db)

	command := NewCMD()

	command.Add("db:migrate up", handler.migrateUp, false)
	command.Add("db:migrate down", handler.migrateDown, true)
	command.Add("db:migrate status", handler.migrationStatus, false)
	command.Add("db:reset", handler.resetDB, true)

	ctx := context.Background()

	command.Run(ctx)
}
