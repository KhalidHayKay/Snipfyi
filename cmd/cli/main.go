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
	input := strings.Join(os.Args[1:], " ")

	cmd, ok := commandMap[input]
	if !ok {
		log.Fatalf("Unknown command: %s", input)
	}

	config.LoadEnv()
	if err := storage.InitDB(); err != nil {
		log.Fatal(err)
	}

	if cmd.Destructive {
		fmt.Print("⚠️  This is a destructive operation. Type 'yes' to proceed: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')

		if strings.TrimSpace(response) != "yes" {
			fmt.Println("Aborted.")
			return
		}
	}

	ctx := context.Background()

	cmd.Run(ctx)
}
