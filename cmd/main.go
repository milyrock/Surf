package main

import (
	"log"

	"github.com/milyrock/Surf/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
