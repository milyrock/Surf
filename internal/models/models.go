package models

import "time"

type RecordCategory string

const (
	Damage     RecordCategory = "Порча"
	Expiration RecordCategory = "Срок годности"
	Lunch      RecordCategory = "Ланч"
	Additional RecordCategory = "Еда"
)

type Record struct {
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Product   string         `db:"product"`
	Category  RecordCategory `db:"category"`
	Amount    int            `db:"amount"`
	CreatedAt time.Time      `db:"created_at"`
}
