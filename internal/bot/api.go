package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/milyrock/Surf/internal/models"
)

const (
	CommandAdd  = "add"
	CommandList = "list"
	CommandUser = "user"
	CommandDate = "date"
)

type Service interface {
	Add(update tgbotapi.Update) (*models.Record, error)
	List(update tgbotapi.Update) error
	GetByUserID(update tgbotapi.Update) error
	GetByDate(update tgbotapi.Update) error
}

type BotAPIConfig struct {
	Service Service
	Api     *tgbotapi.BotAPI
}

type Bot struct {
	service Service
	api     *tgbotapi.BotAPI
}

func NewAPI(config BotAPIConfig) *Bot {
	return &Bot{service: config.Service, api: config.Api}
}

func (b *Bot) Run() error {
	updates := b.api.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case CommandAdd:
			_, err := b.service.Add(update)
			if err != nil {
				b.reply(update, fmt.Sprintf("Ошибка: %v", err))
			}
		case CommandList:
			if err := b.service.List(update); err != nil {
				b.reply(update, fmt.Sprintf("Ошибка: %v", err))
			}
		case CommandUser:
			if err := b.service.GetByUserID(update); err != nil {
				b.reply(update, fmt.Sprintf("Ошибка: %v", err))
			}
		case CommandDate:
			if err := b.service.GetByDate(update); err != nil {
				b.reply(update, fmt.Sprintf("Ошибка: %v", err))
			}
		default:
			b.reply(update, "Неизвестная команда. Используйте:\n/add <продукт> <тип> <количество>\n/list - все записи\n/user <username> - записи пользователя\n/date <YYYY-MM-DD> - записи за дату")
		}
	}

	return nil
}

func (b *Bot) reply(update tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	b.api.Send(msg)
}
