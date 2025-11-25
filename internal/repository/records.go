package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/milyrock/Surf/internal/models"
)

const (
	insertRecord = `INSERT INTO surf.records (username, product, category, amount) values ($1,$2,$3,$4) returning id`
)

type RecordRepository struct {
	db *sqlx.DB
}

func NewRecordRepository(db *sqlx.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func (r *RecordRepository) Create(rec *models.Record) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = tx.QueryRowx(insertRecord, rec.Name, rec.Product, rec.Category, rec.Amount).Scan(&rec.ID)
	if err != nil{
		fmt.Println(err)
	}
	
	return nil
}
