package models

import (
	"encoding/json"
	"time"
)

// Lot - структура для хранения данных лота.
type Lot struct {
	ID               int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	LotID            string          `json:"lot_id"`
	AdvertID         string          `json:"advert_id"`
	Name             string          `json:"name"`
	Status           string          `json:"status"`
	CreateDate       time.Time       `json:"createdate"`
	Organizer        string          `json:"organizer"`
	Price            string          `json:"price"`
	URL              string          `json:"url"`
	LotType          string          `json:"lottype"`
	StartDate        time.Time       `json:"startdate"`
	EndDate          time.Time       `json:"enddate"`
	LinkDownloadFile string          `json:"linkdownloadfile"`
	TechnicalSpec    int64           `json:"technicalspec"`
	LocalFilePaths   json.RawMessage `json:"local_file_paths"`
}

type RelevantLot struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	LotID     string `json:"lot_id"`
	CompanyID int64  `json:"company_id"`
}

// Metadata - структура для работы с пагинацией.
type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
	PrevPage     int `json:"prev_page"`
	NextPage     int `json:"next_page"`
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	Price        string
	Regions      string
	StartDate    string
	EndDate      string
	SortSafelist []string
}
