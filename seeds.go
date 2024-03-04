package main

import (
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/seeders/seeds"
	"log"
	"os/exec"
)

func init() {
	err := config.LoadENV()
	if err != nil {
		log.Fatalf("Error loading environment variables")
	} else {
		fmt.Println("Environment variables loaded!")
	}

	config.ConnectDb()
}

func main() {
	// Example: run "ls" command
	output, err := runCommand("make", "goose-refresh")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Command Output:")
	fmt.Println(output)

	fmt.Println("Running seeds...")
	for _, seed := range seeds.All() {
		if err := seed.Run(); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
	fmt.Println("Seeds ran successfully!")

}

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command: %v", err)
	}

	return string(output), nil
}
