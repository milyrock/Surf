package main

import (
	"fmt"
	"log"

	"github.com/milyrock/Surf/internal/app"
	"github.com/milyrock/Surf/internal/config"
)

func main() {
	fmt.Println("yo")
	cfg, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	fmt.Println(cfg)

	fmt.Printf("%#v", cfg)

	bot, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init telegram bot: %v", err)
	}

	fmt.Println(bot)
}
