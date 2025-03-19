package models

type Company struct {
	ID   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	BIN  string `json:"bin"`
}
