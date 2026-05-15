package main

import (
	"context"
	"log"
	"smply/app/storage"
	"smply/config"
)

func main() {
	config.LoadEnv()

	pgsql, err := storage.InitPostgres()
	if err != nil {
		log.Fatal(err)
	}

	handler := NewCLIHandler(pgsql)

	command := NewCMD()

	command.Add("db:migrate up", handler.migrateUp, false)
	command.Add("db:migrate down", handler.migrateDown, true)
	command.Add("db:migrate status", handler.migrationStatus, false)
	command.Add("db:reset", handler.resetDB, true)

	ctx := context.Background()

	command.Run(ctx)
}
