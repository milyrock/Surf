package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/milyrock/Surf/internal/models"
)

type RecordRepository interface {
	Create(rec *models.Record) error
}

type Service struct {
	recordRepo RecordRepository
	bot        *tgbotapi.BotAPI
	// googleSheetsService GoogleSheetsAPI
}

type Config struct {
	RecordRepo RecordRepository
	Bot        *tgbotapi.BotAPI
	// GoogleSheetsService GoogleSheetsAPI
}

func NewService(config *Config) *Service {
	return &Service{recordRepo: config.RecordRepo, bot: config.Bot}
}
