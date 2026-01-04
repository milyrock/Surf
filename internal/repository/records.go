package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/milyrock/Surf/internal/models"
)

const (
	insertRecord     = `INSERT INTO surf.records (username, product, category, amount) values ($1,$2,$3,$4) returning id, created_at`
	selectRecords    = `SELECT id, username as name, product, category, amount, created_at FROM surf.records ORDER BY id DESC`
	selectByUser     = `SELECT id, username as name, product, category, amount, created_at FROM surf.records WHERE username = $1 ORDER BY id DESC`
	selectByDate     = `SELECT id, username as name, product, category, amount, created_at FROM surf.records WHERE DATE(created_at) = $1 ORDER BY id DESC`
	selectByUserDate = `SELECT id, username as name, product, category, amount, created_at FROM surf.records WHERE username = $1 AND DATE(created_at) = $2 ORDER BY id DESC`
)

type RecordRepository struct {
	db *sqlx.DB
}

func NewRecordRepository(db *sqlx.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func mapCategoryToDBEnum(category models.RecordCategory) string {
	switch category {
	case models.Damage:
		return "порча"
	case models.Expiration:
		return "срок_годности"
	case models.Lunch:
		return "ланч"
	case models.Additional:
		return "еда"
	default:
		return strings.ToLower(string(category))
	}
}

func (r *RecordRepository) Create(rec *models.Record) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	categoryStr := mapCategoryToDBEnum(rec.Category)

	err = tx.QueryRowx(insertRecord, rec.Name, rec.Product, categoryStr, rec.Amount).Scan(&rec.ID, &rec.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func mapDBEnumToCategory(categoryStr string) models.RecordCategory {
	switch categoryStr {
	case "порча":
		return models.Damage
	case "срок_годности":
		return models.Expiration
	case "ланч":
		return models.Lunch
	case "еда":
		return models.Additional
	default:
		return models.RecordCategory(categoryStr)
	}
}

func (r *RecordRepository) GetAll() ([]*models.Record, error) {
	rows, err := r.db.Queryx(selectRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	return r.scanRecords(rows)
}

func (r *RecordRepository) GetByUserID(username string) ([]*models.Record, error) {
	rows, err := r.db.Queryx(selectByUser, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query records by user: %w", err)
	}
	defer rows.Close()

	return r.scanRecords(rows)
}

func (r *RecordRepository) GetByDate(date time.Time) ([]*models.Record, error) {
	dateStr := date.Format("2006-01-02")
	rows, err := r.db.Queryx(selectByDate, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query records by date: %w", err)
	}
	defer rows.Close()

	return r.scanRecords(rows)
}

func (r *RecordRepository) GetByUserIDAndDate(username string, date time.Time) ([]*models.Record, error) {
	dateStr := date.Format("2006-01-02")
	rows, err := r.db.Queryx(selectByUserDate, username, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query records by user and date: %w", err)
	}
	defer rows.Close()

	return r.scanRecords(rows)
}

func (r *RecordRepository) scanRecords(rows *sqlx.Rows) ([]*models.Record, error) {
	var records []*models.Record
	for rows.Next() {
		var rec models.Record
		var categoryStr string

		err := rows.Scan(&rec.ID, &rec.Name, &rec.Product, &categoryStr, &rec.Amount, &rec.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		rec.Category = mapDBEnumToCategory(categoryStr)
		records = append(records, &rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating records: %w", err)
	}

	return records, nil
}
