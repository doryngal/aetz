package models

import (
	"math"
	"strings"
	// "binai.net/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
	Price        string
	Regions      string
	StartDate    string
	EndDate      string
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) price() string {
	return f.Price
}

// func ValidateFilters(v *validator.Validator, f Filters) {
// 	v.Check(f.Page > 0, "page", "must be greater than zero")
// 	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
// 	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
// 	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

// 	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
// }

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`

	// Добавляем поля для предыдущей и следующей страницы
	PrevPage int `json:"prev_page,omitempty"`
	NextPage int `json:"next_page,omitempty"`
}

// func calculateMetadata(totalRecords, page, pageSize int) Metadata {
// 	if totalRecords == 0 {
// 		return Metadata{}
// 	}

// 	return Metadata{
// 		CurrentPage:  page,
// 		PageSize:     pageSize,
// 		FirstPage:    1,
// 		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
// 		TotalRecords: totalRecords,
// 	}
// }

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	metadata := Metadata{
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
