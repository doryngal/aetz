package repository

import (
	models2 "binai.net/v2/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// allowedSortColumns – whitelist допустимых столбцов для сортировки.
var allowedSortColumns = map[string]bool{
	"id":               true,
	"advert_id":        true,
	"name":             true,
	"status":           true,
	"createdate":       true,
	"organizer":        true,
	"price":            true,
	"url":              true,
	"lottype":          true,
	"startdate":        true,
	"enddate":          true,
	"linkdownloadfile": true,
}

// LotRepository описывает интерфейс для работы с лотами.
type LotRepository interface {
	// FindRelevantLots возвращает список релевантных лотов по названию компании,
	// поисковому запросу и фильтрам, а также метаданные для пагинации.
	FindRelevantLots(companyName, searchQuery string, filters models2.Filters) ([]models2.Lot, models2.Metadata, error)
	// FindLotByID возвращает лот по его идентификатору.
	FindLotByID(id int) (*models2.Lot, error)
	// CountLots возвращает общее количество лотов для заданной компании.
	CountLots(companyName string) (int, error)
}

type lotRepository struct {
	db *sql.DB
}

// NewLotRepository создаёт новый экземпляр репозитория лотов.
func NewLotRepository(db *sql.DB) LotRepository {
	return &lotRepository{db: db}
}

// FindRelevantLots получает список релевантных лотов согласно заданным фильтрам.
func (r *lotRepository) FindRelevantLots(companyName, searchQuery string, filters models2.Filters) ([]models2.Lot, models2.Metadata, error) {
	// Устанавливаем значения по умолчанию для пагинации.
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10 // значение по умолчанию
	}

	// Получаем безопасные значения сортировки.
	//sortColumn, sortOrder := sanitizeSort(filters.Sort)
	query := fmt.Sprintf(`
        SELECT
            l.id, l.advert_id, l.name, l.status, l.organizer, l.price,
            l.url, l.lottype, l.startdate, l.enddate, l.linkdownloadfile
        FROM lots l`)
	// Формируем SQL-запрос.
	// Заметим, что ORDER BY формируется через Sprintf, но переменная sortColumn прошла валидацию.
	//query := fmt.Sprintf(`
	//    SELECT
	//        l.id, l.advert_id, l.name, l.status, l.createdate, l.organizer, l.price,
	//        l.url, l.lottype, l.startdate, l.enddate, l.linkdownloadfile
	//    FROM lots l
	//    JOIN relevant_lot r ON l.id = r.lot_id
	//    JOIN company c ON c.id = r.company_id
	//    WHERE c.name = $1
	//      AND (to_tsvector('simple', l.organizer || ' ' || l.name) @@ plainto_tsquery('simple', $2) OR $2 = '')
	//    ORDER BY %s %s
	//    LIMIT $3 OFFSET $4
	//`, sortColumn, sortOrder)

	// Создаём контекст с таймаутом для запроса.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, query)
	//rows, err := r.db.QueryContext(ctx, query, companyName, searchQuery, filters.PageSize, (filters.Page-1)*filters.PageSize)
	if err != nil {
		return nil, models2.Metadata{}, fmt.Errorf("querying relevant lots: %w", err)
	}
	defer rows.Close()

	var lots []models2.Lot
	for rows.Next() {
		var lot models2.Lot
		if err = rows.Scan(
			&lot.ID,
			&lot.AdvertID,
			&lot.Name,
			&lot.Status,
			&lot.Organizer,
			&lot.Price,
			&lot.URL,
			&lot.LotType,
			&lot.StartDate,
			&lot.EndDate,
			&lot.LinkDownloadFile,
		); err != nil {
			return nil, models2.Metadata{}, fmt.Errorf("scanning lot row: %w", err)
		}
		lots = append(lots, lot)
	}
	if err = rows.Err(); err != nil {
		return nil, models2.Metadata{}, fmt.Errorf("rows error: %w", err)
	}

	//totalRecords, err := r.CountLots(companyName)
	//if err != nil {
	//	return nil, models2.Metadata{}, fmt.Errorf("counting lots: %w", err)
	//}

	metadata := calculateMetadata(0, filters.Page, filters.PageSize)
	return lots, metadata, nil
}

// FindLotByID возвращает лот по заданному идентификатору.
func (r *lotRepository) FindLotByID(id int) (*models2.Lot, error) {
	query := `
        SELECT id, advert_id, name, status, createdate, organizer, price, url, 
               lottype, startdate, enddate, linkdownloadfile
        FROM lots
        WHERE id = $1
    `
	var lot models2.Lot
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lot.ID,
		&lot.AdvertID,
		&lot.Name,
		&lot.Status,
		&lot.Organizer,
		&lot.Price,
		&lot.URL,
		&lot.LotType,
		&lot.StartDate,
		&lot.EndDate,
		&lot.LinkDownloadFile,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models2.ErrNoRecord
		}
		return nil, fmt.Errorf("querying lot by id: %w", err)
	}
	return &lot, nil
}

// CountLots возвращает количество лотов для заданной компании.
func (r *lotRepository) CountLots(companyName string) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM lots l
        JOIN relevant_lots r ON l.id = r.lot_id
        JOIN company c ON c.id = r.company_id
        WHERE c.name = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	if err := r.db.QueryRowContext(ctx, query, companyName).Scan(&count); err != nil {
		return 0, fmt.Errorf("counting lots for company %s: %w", companyName, err)
	}
	return count, nil
}

// sanitizeSort проверяет и возвращает безопасные значения столбца сортировки и порядка.
func sanitizeSort(sort string) (column, order string) {
	// Значения по умолчанию.
	column = "id"
	order = "DESC"

	if sort == "" {
		return
	}

	// Определяем порядок сортировки.
	if strings.HasPrefix(sort, "-") {
		order = "DESC"
		sort = strings.TrimPrefix(sort, "-")
	} else {
		order = "ASC"
	}

	// Приводим к нижнему регистру для сравнения и проверяем whitelist.
	if _, ok := allowedSortColumns[strings.ToLower(sort)]; ok {
		column = sort
	}
	return
}

// calculateMetadata вычисляет данные для пагинации.
func calculateMetadata(totalRecords, page, pageSize int) models2.Metadata {
	if totalRecords == 0 {
		return models2.Metadata{}
	}

	totalPages := (totalRecords + pageSize - 1) / pageSize
	metadata := models2.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     totalPages,
		TotalRecords: totalRecords,
	}

	if page > 1 {
		metadata.PrevPage = page - 1
	}
	if page < totalPages {
		metadata.NextPage = page + 1
	}
	return metadata
}
