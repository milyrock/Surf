package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/milyrock/Surf/internal/models"
)

const (
	CategoryDamage     = "damage"
	CategoryExpiration = "expiration"
	CategoryLunch      = "lunch"
	CategoryAdditional = "additional"
	CategoryFood       = "food"

	CategoryDamageRU      = "порча"
	CategoryExpirationRU  = "срок"
	CategoryExpirationRU2 = "годности"
	CategoryLunchRU       = "ланч"
	CategoryAdditionalRU  = "еда"
)

func (s *Service) Add(update tgbotapi.Update) (*models.Record, error) {
	if update.Message == nil {
		return nil, fmt.Errorf("message is nil")
	}

	args := strings.Fields(update.Message.CommandArguments())
	if len(args) < 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /add <продукт> <тип> <количество>\nПример: /add молоко порча 5\n\nТипы: порча, срок годности, ланч, еда")
		s.bot.Send(msg)
		return nil, fmt.Errorf("insufficient arguments: expected product, type, and amount")
	}

	product := args[0]
	var categoryStr string
	var amountStr string

	if len(args) >= 4 && strings.ToLower(args[1]) == CategoryExpirationRU && strings.ToLower(args[2]) == CategoryExpirationRU2 {
		categoryStr = strings.ToLower(args[1] + " " + args[2])
		amountStr = args[3]
	} else {
		categoryStr = strings.ToLower(args[1])
		amountStr = args[2]
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Неверное количество: %s. Количество должно быть числом.", amountStr))
		s.bot.Send(msg)
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	var category models.RecordCategory
	switch categoryStr {
	case CategoryDamage, CategoryDamageRU:
		category = models.Damage
	case CategoryExpiration, CategoryExpirationRU, CategoryExpirationRU2, "срок годности":
		category = models.Expiration
	case CategoryLunch, CategoryLunchRU:
		category = models.Lunch
	case CategoryAdditional, CategoryAdditionalRU, CategoryFood:
		category = models.Additional
	default:
		validCategories := fmt.Sprintf("%s, %s, %s, %s", CategoryDamageRU, "срок годности", CategoryLunchRU, CategoryAdditionalRU)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Неверный тип: %s. Доступные типы: %s", categoryStr, validCategories))
		s.bot.Send(msg)
		return nil, fmt.Errorf("invalid category: %s", categoryStr)
	}

	username := update.Message.From.UserName
	if username == "" {
		username = fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)
		if strings.TrimSpace(username) == "" {
			username = fmt.Sprintf("user_%d", update.Message.From.ID)
		}
	}

	rec := &models.Record{
		Name:     username,
		Product:  product,
		Category: category,
		Amount:   amount,
	}

	err = s.recordRepo.Create(rec)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при создании записи: %v", err))
		s.bot.Send(msg)
		return nil, fmt.Errorf("couldn't create record: %w", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("✅ Запись успешно создана!\nПродукт: %s\nТип: %s\nКоличество: %d", product, category, amount))
	s.bot.Send(msg)

	return rec, nil
}

func (s *Service) List(update tgbotapi.Update) error {
	if update.Message == nil {
		return fmt.Errorf("message is nil")
	}

	records, err := s.recordRepo.GetAll()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при получении записей: %v", err))
		s.bot.Send(msg)
		return fmt.Errorf("couldn't get records: %w", err)
	}

	if len(records) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📋 Записей пока нет. Используйте /add для добавления.")
		s.bot.Send(msg)
		return nil
	}

	var builder strings.Builder
	builder.WriteString("📋 Список записей:\n\n")

	for i, rec := range records {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Product))
		builder.WriteString(fmt.Sprintf("   Тип: %s\n", rec.Category))
		builder.WriteString(fmt.Sprintf("   Количество: %d\n", rec.Amount))
		builder.WriteString(fmt.Sprintf("   Пользователь: %s\n", rec.Name))
		builder.WriteString(fmt.Sprintf("   Дата: %s\n", rec.CreatedAt.Format("02.01.2006 15:04")))
		if i < len(records)-1 {
			builder.WriteString("\n")
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	s.bot.Send(msg)

	return nil
}

func (s *Service) GetByUserID(update tgbotapi.Update) error {
	if update.Message == nil {
		return fmt.Errorf("message is nil")
	}

	args := strings.Fields(update.Message.CommandArguments())
	if len(args) < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /user <username>\nПример: /user ilyrock")
		s.bot.Send(msg)
		return fmt.Errorf("insufficient arguments: expected username")
	}

	username := args[0]
	records, err := s.recordRepo.GetByUserID(username)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при получении записей: %v", err))
		s.bot.Send(msg)
		return fmt.Errorf("couldn't get records by user: %w", err)
	}

	if len(records) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("📋 Записей для пользователя '%s' не найдено.", username))
		s.bot.Send(msg)
		return nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("📋 Записи пользователя %s:\n\n", username))

	for i, rec := range records {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Product))
		builder.WriteString(fmt.Sprintf("   Тип: %s\n", rec.Category))
		builder.WriteString(fmt.Sprintf("   Количество: %d\n", rec.Amount))
		builder.WriteString(fmt.Sprintf("   Дата: %s\n", rec.CreatedAt.Format("02.01.2006 15:04")))
		if i < len(records)-1 {
			builder.WriteString("\n")
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	s.bot.Send(msg)

	return nil
}

func (s *Service) GetByDate(update tgbotapi.Update) error {
	if update.Message == nil {
		return fmt.Errorf("message is nil")
	}

	args := strings.Fields(update.Message.CommandArguments())
	if len(args) < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /date <дата>\nПример: /date 2026-01-04\nФормат: YYYY-MM-DD")
		s.bot.Send(msg)
		return fmt.Errorf("insufficient arguments: expected date")
	}

	dateStr := args[0]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Неверный формат даты: %s. Используйте формат YYYY-MM-DD\nПример: 2026-01-04", dateStr))
		s.bot.Send(msg)
		return fmt.Errorf("invalid date format: %w", err)
	}

	records, err := s.recordRepo.GetByDate(date)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при получении записей: %v", err))
		s.bot.Send(msg)
		return fmt.Errorf("couldn't get records by date: %w", err)
	}

	if len(records) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("📋 Записей за %s не найдено.", date.Format("02.01.2006")))
		s.bot.Send(msg)
		return nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("📋 Записи за %s:\n\n", date.Format("02.01.2006")))

	for i, rec := range records {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec.Product))
		builder.WriteString(fmt.Sprintf("   Тип: %s\n", rec.Category))
		builder.WriteString(fmt.Sprintf("   Количество: %d\n", rec.Amount))
		builder.WriteString(fmt.Sprintf("   Пользователь: %s\n", rec.Name))
		builder.WriteString(fmt.Sprintf("   Время: %s\n", rec.CreatedAt.Format("15:04")))
		if i < len(records)-1 {
			builder.WriteString("\n")
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	s.bot.Send(msg)

	return nil
}
