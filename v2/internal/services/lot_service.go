package services

import (
	"binai.net/v2/internal/models"
	"binai.net/v2/internal/repository"
	"fmt"
)

// LotService описывает бизнес-логику для работы с лотами.
type LotService interface {
	// GetLotList возвращает список лотов и метаданные пагинации для заданной компании и поискового запроса.
	GetLotList(companyName, searchQuery string, filters models.Filters) ([]models.Lot, models.Metadata, error)
	// GetLotByID возвращает лот по его идентификатору.
	GetLotByID(id int) (*models.Lot, error)
}

type lotService struct {
	repo repository.LotRepository
}

// NewLotService создаёт новый экземпляр сервиса для лотов.
func NewLotService(repo repository.LotRepository) LotService {
	return &lotService{
		repo: repo,
	}
}

// GetLotList получает список лотов по заданным параметрам.
// Если параметры пагинации не заданы корректно, устанавливаются значения по умолчанию.
// Если companyName не указан, возвращается ошибка.
func (s *lotService) GetLotList(companyName, searchQuery string, filters models.Filters) ([]models.Lot, models.Metadata, error) {
	// Проверяем обязательный параметр companyName.
	if companyName == "" {
		return nil, models.Metadata{}, fmt.Errorf("company name is required")
	}

	// Устанавливаем значения по умолчанию для пагинации.
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}

	// Дополнительная бизнес-логика, например, санитация поискового запроса или логирование, может быть добавлена здесь.

	// Передаём запрос в репозиторий.
	lots, metadata, err := s.repo.FindRelevantLots(companyName, searchQuery, filters)
	if err != nil {
		return nil, models.Metadata{}, fmt.Errorf("failed to retrieve lots: %w", err)
	}

	return lots, metadata, nil
}

// GetLotByID получает лот по его идентификатору.
// Если идентификатор меньше единицы, возвращается ошибка.
func (s *lotService) GetLotByID(id int) (*models.Lot, error) {
	if id < 1 {
		return nil, fmt.Errorf("invalid lot id: %d", id)
	}

	lot, err := s.repo.FindLotByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve lot: %w", err)
	}

	return lot, nil
}
