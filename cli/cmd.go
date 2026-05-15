package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type CMD struct {
	commands map[string]Command
}

func NewCMD() *CMD {
	commands := make(map[string]Command)

	commands["list"] = Command{
		Run: func(ctx context.Context) {
			PrintCommands(commands)
		},
		Destructive: false,
	}

	return &CMD{
		commands: commands,
	}
}

func (c *CMD) Add(name string, handler func(context.Context), destructive bool) {
	c.commands[name] = Command{
		Run:         handler,
		Destructive: destructive,
	}
}

func (c *CMD) Run(ctx context.Context) {
	input := strings.Join(os.Args[1:], " ")

	cmd, ok := c.commands[input]
	if !ok {
		log.Fatalf("Unknown command: %s", input)
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

	cmd.Run(ctx)
}

func PrintCommands(commands map[string]Command) {
	fmt.Println("Available commands:")

	// Commands to exclude
	excluded := map[string]bool{
		"list": true,
		// add more here
		// "debug": true,
	}

	// Collect allowed command names
	keys := make([]string, 0, len(commands))
	for k := range commands {
		if !excluded[k] {
			keys = append(keys, k)
		}
	}

	// Sort alphabetically
	sort.Strings(keys)

	// Print commands
	for _, cmd := range keys {
		fmt.Printf("- %s\n", cmd)
	}
}
