package models

type RecordCategory string

const (
	Damage     RecordCategory = "Порча"
	Expiration RecordCategory = "Срок годности"
	Lunch      RecordCategory = "Ланч"
	Additional RecordCategory = "Еда"
)

type Record struct {
	ID       int `db:"id"`
	Name     string
	Product  string
	Category RecordCategory
	Amount   int
}
