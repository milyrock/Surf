package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotConfig struct {
	Bot        *tgbotapi.BotAPI
	RecordRepo RecordRepository
	//GoogleSheetsService GoogleSheetsAPI
}

type BotService struct {
	bot        *tgbotapi.BotAPI
	recordRepo RecordRepository
	//googleSheetsService GoogleSheetsAPI
}

func NewBotService(cfg *BotConfig) *BotService {
	return &BotService{
		bot:        cfg.Bot,
		recordRepo: cfg.RecordRepo,
		//googleSheetsService:   cfg.GoogleSheetsService,
	}
}
