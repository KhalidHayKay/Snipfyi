package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"smply/config"
	"strings"
)

func main() {
	// 1. Force a manual confirmation
	fmt.Print("CRITICAL: This will delete ALL data. Type 'yes' to proceed: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	if strings.TrimSpace(response) != "yes" {
		fmt.Println("Aborted.")
		return
	}

	config.LoadEnv()
	if err := config.InitDB(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	tables := []string{"urls", "api_keys", "magic_tokens"}

	for _, table := range tables {
		// DROP TABLE removes the entire table structure
		// IF EXISTS prevents errors if the table is already gone
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)

		_, err := config.DB.Exec(ctx, query)
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("Successfully DROPPED table %s", table)
	}

}
