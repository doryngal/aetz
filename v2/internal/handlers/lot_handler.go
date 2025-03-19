package handlers

import (
	models2 "binai.net/v2/internal/models"
	"binai.net/v2/internal/repository"
	"binai.net/v2/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LotHandler отвечает за обработку HTTP-запросов, связанных с лотами.
type LotHandler struct {
	lotService services.LotService
}

// NewLotHandler создаёт новый экземпляр обработчика лотов и регистрирует маршруты.
func NewLotHandler(repo repository.LotRepository) LotHandler {
	service := services.NewLotService(repo)
	return LotHandler{
		lotService: service,
	}
}

// GetLotList обрабатывает запрос на получение списка лотов.
// Он извлекает параметры запроса, формирует фильтры и возвращает данные с метаданными для пагинации.
func (h *LotHandler) GetLotList(c *gin.Context) {
	// Извлекаем параметры запроса.
	companyName := c.Query("company")
	searchQuery := c.Query("search")

	// Извлечение параметров пагинации и сортировки.
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size parameter"})
		return
	}

	// Формируем фильтры.
	filters := models2.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     c.DefaultQuery("sort", "-createdate"),
	}

	// Вызываем сервис для получения данных.
	lots, metadata, err := h.lotService.GetLotList(companyName, searchQuery, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем ответ в формате JSON.
	c.JSON(http.StatusOK, gin.H{
		"data":     lots,
		"metadata": metadata,
	})
}

// GetLotByID обрабатывает запрос на получение одного лота по его идентификатору.
func (h *LotHandler) GetLotByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lot id"})
		return
	}

	lot, err := h.lotService.GetLotByID(id)
	if err != nil {
		// Если лот не найден, возвращаем 404.
		if err == models2.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lot)
}
