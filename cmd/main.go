package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if err := HandleBotCommands(update); err != nil {
			fmt.Errorf("cannot handle bot command : %v", err)
		}
	}

}

func HandleBotCommands(update tgbotapi.Update) error {
	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}

	switch update.Message.Command() {
	case "add":
		return HandleAddItemCommand(update.Message)
	default:
		return nil
	}
}

func HandleAddItemCommand(message *tgbotapi.Message) error {
	fmt.Println("we are here")
	fmt.Println(message.Text)
	return nil
}
