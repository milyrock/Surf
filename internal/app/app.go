package app

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/milyrock/Surf/internal/bot"
	"github.com/milyrock/Surf/internal/config"
	"github.com/milyrock/Surf/internal/repository"
	"github.com/milyrock/Surf/internal/service"
)

func Run() error {
	cfg, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		return err
	}

	botAPI, err := initBot(cfg)
	if err != nil {
		return err
	}

	db, err := InitDB(cfg.Database.Postgres.ConnConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	recordRepo := repository.NewRecordRepository(db)

	serviceConfig := &service.Config{
		RecordRepo: recordRepo,
		Bot:        botAPI,
	}
	botService := service.NewService(serviceConfig)

	botConfig := bot.BotAPIConfig{
		Service: botService,
		Api:     botAPI,
	}
	botAPIHandler := bot.NewAPI(botConfig)

	log.Println("Bot is running...")
	return botAPIHandler.Run()
}

func initBot(cfg *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}
