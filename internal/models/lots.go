package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type LotModelInterface interface {
	GetRelevantLotList(companyName, searchStr string, filters Filters) ([]Lot, Metadata, error)
	GetRelevantLotById(id int) (*Lot, error)
	GetLots(companyName string) ([]Lot, error)
}

type LotModel struct {
	DB *sql.DB
}

type Lot struct {
	ID               int            `json:"id"`
	AdvertID         string         `json:"advert_id"`
	Name             string         `json:"name"`
	Status           string         `json:"status"`
	CreateDate       time.Time      `json:"create_date"`
	Organizer        string         `json:"organizer"`
	Price            string         `json:"price"`
	URL              string         `json:"url"`
	LotType          string         `json:"lot_type"`
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
	LinkDownloadFile string         `json:"link_download_file"`
	DeliveryPlace    string         `json:"delivery_place"`
	ParentLotLink    sql.NullString `json:"parent_lot"`
	TechnicalSpec    int            `json:"technical_spec"`
	Now              time.Time
}

func (l *LotModel) GetRelevantLotList(companyName, searchStr string, filters Filters) ([]Lot, Metadata, error) {
	fmt.Println(searchStr, filters)
	validColumns := map[string]bool{
		"id": true, "advert_id": true, "name": true, "status": true,
		"createdate": true, "organizer": true, "price": true, "url": true,
		"lottype": true, "startdate": true, "enddate": true, "linkdownloadfile": true,
	}

	sortColumn := filters.sortColumn()
	sortDirection := filters.sortDirection()

	if !validColumns[sortColumn] {
		sortColumn = "id" // дефолтное поле сортировки
	}

	if sortDirection != "ASC" && sortDirection != "DESC" {
		sortDirection = "ASC" // дефолтное направление
	}

	query := fmt.Sprintf(`
SELECT DISTINCT
    l.id, l.advert_id, l.name, l.status, l.createdate, l.organizer, 
    l.price, l.url, l.lottype, l.startdate, l.enddate, l.linkdownloadfile
FROM public.lots l
JOIN public.relevant_lots r ON l.lot_id = r.lot_id;`)

	// select count(*) from lots l JOIN relevant_lot r ON l.id = r.lot_id where (LOWER(l.status) LIKE '%%(прием заявок)%%');
	// AND (LOWER(l.status) LIKE '%%(прием заявок)%%')

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// fmt.Println("filter.StartDate", filters.StartDate, "filter.EndDate", filters.EndDate)
	//args := []any{companyName, searchStr, filters.limit(), filters.offset(), filters.price(), filters.StartDate}
	// args := []any{companyName, searchStr, filters.limit(), filters.offset(), filters.price()}
	// args := []any{companyName}

	rows, err := l.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("empty Metadata", err)
		return nil, Metadata{
			CurrentPage:  0,
			PageSize:     0,
			FirstPage:    0,
			LastPage:     0,
			TotalRecords: 0,
			PrevPage:     0,
			NextPage:     0,
		}, err
	}

	defer rows.Close()

	fmt.Println(companyName)
	// totalRecords, _ := l.GetLots(companyName)
	lots := []Lot{}
	totalRecords := 0
	for rows.Next() {
		var lot Lot
		totalRecords++

		err := rows.Scan(
			&lot.ID,
			&lot.AdvertID,
			&lot.Name,

			&lot.Status,
			&lot.CreateDate,
			&lot.Organizer,

			&lot.Price,
			&lot.URL,
			&lot.LotType,

			&lot.StartDate,
			&lot.EndDate,
			&lot.LinkDownloadFile,
		)
		if err != nil {
			fmt.Println("totalRecords", err)
			return nil, Metadata{}, err
		}

		parentAdvert_id := l.ParentAdvertId(lot.ParentLotLink)

		if parentAdvert_id.Valid {
			// Если parentAdvert_id валидный, соединяем его с базовым URL
			lot.ParentLotLink = sql.NullString{
				String: "https://goszakup.gov.kz/ru/announce/index/" + parentAdvert_id.String,
				Valid:  true, // Устанавливаем Valid в true, так как строка не пустая
			}
		} else {
			// Если parentAdvert_id не валидный, можно использовать пустую строку или другую логику
			lot.ParentLotLink = sql.NullString{
				String: "",    // Здесь можно использовать пустую строку
				Valid:  false, // Устанавливаем Valid в false, так как значение NULL
			}
		}

		lots = append(lots, lot)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("ERR: ", err)
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return lots, metadata, nil
}

func (l *LotModel) ParentAdvertId(parentLot_id sql.NullString) sql.NullString {
	query := `
		SELECT advert_id
		FROM lots

	`
	var advertID string

	// Выполнение запроса
	err := l.DB.QueryRow(query, parentLot_id).Scan(&advertID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если строки с таким id нет, можно вернуть пустую строку или обработать по-другому
			fmt.Println("No matching lot found for parent_lot_id:", parentLot_id)
			return sql.NullString{String: "", Valid: false}
		}
		// Ошибка запроса
		fmt.Println("Error fetching advert_id:", err)
		return sql.NullString{String: "", Valid: false}
	}

	// Возвращаем найденный advert_id
	return sql.NullString{String: advertID, Valid: true}
}

func (l *LotModel) GetRelevantLotById(id int) (*Lot, error) {
	query := `
		SELECT
			l.id,
			l.advert_id,
			l.name,
			l.status,
			l.createdate,
			l.organizer,
			l.price,
			l.url,
			l.lottype,
			l.startdate,
			l.enddate,
			l.linkdownloadfile
		FROM lots l WHERE l.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	lot := &Lot{}

	row := l.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&lot.ID,
		&lot.AdvertID,
		&lot.Name,
		&lot.Status,
		&lot.CreateDate,
		&lot.Organizer,
		&lot.Price,
		&lot.URL,
		&lot.LotType,
		&lot.StartDate,
		&lot.EndDate,
		&lot.LinkDownloadFile,
	)
	if err != nil {
		return nil, err
	}

	return lot, nil
}

func (l *LotModel) GetLots(companyName string) ([]Lot, error) {
	query := `
    SELECT DISTINCT
			l.id,
			l.advert_id,
			l.name,
			l.status,
			l.createdate,
			l.organizer,
			l.price,
			l.url,
			l.lottype,
			l.startdate,
			l.enddate,
			l.linkdownloadfile
		FROM lots l
		JOIN relevant_lots r ON l.id = r.lot_id
		JOIN company c ON c.id = r.company_id
		WHERE c.name = $1
		AND (LOWER(l.status) LIKE '%%(прием заявок)%%')
		ORDER BY l.id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := l.DB.QueryContext(ctx, query, companyName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := []Lot{}

	for rows.Next() {
		var lot Lot

		err := rows.Scan(
			&lot.ID,
			&lot.AdvertID,
			&lot.Name,
			&lot.Status,
			&lot.CreateDate,
			&lot.Organizer,
			&lot.Price,
			&lot.URL,
			&lot.LotType,
			&lot.StartDate,
			&lot.EndDate,
			&lot.LinkDownloadFile,
		)
		if err != nil {
			return nil, err
		}

		lots = append(lots, lot)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("ERR: ", err)
		return nil, err
	}

	return lots, nil
}
