package main

import (
	"fmt"
	"main/internal/config"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("There is no enough args: go run cmd/main.go <config/path.yaml> <.env>")
		os.Exit(1)
	}

	_, err := config.Load(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

}
