package rr_intervals_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

// GetDifferentialHistogramRequestDTO представляет HTTP запрос на получение дифференциальной гистограммы
type GetDifferentialHistogramRequestDTO struct {
	From      string `form:"from" binding:"required"`      // RFC3339 формат
	To        string `form:"to" binding:"required"`        // RFC3339 формат
	BinsCount int    `form:"bins_count,omitempty"`         // 0 = авто, 10-25
}

// GetDifferentialHistogramResponseDTO представляет HTTP ответ с дифференциальной гистограммой
type GetDifferentialHistogramResponseDTO struct {
	*DifferentialHistogramDataDTO
}

func (r *RRIntervalsRouter) GetDifferentialHistogram(c *gin.Context) {
	var reqDTO GetDifferentialHistogramRequestDTO
	if err := c.ShouldBindQuery(&reqDTO); err != nil {
		www.HandleError(c, err)
		return
	}

	// Парсим временные параметры
	from, err := time.Parse(time.RFC3339, reqDTO.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid from time format (expected RFC3339)",
		})
		return
	}

	to, err := time.Parse(time.RFC3339, reqDTO.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid to time format (expected RFC3339)",
		})
		return
	}

	// Валидируем bins count
	binsCount := reqDTO.BinsCount
	if binsCount < 0 || binsCount > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bins_count must be between 0 and 50",
		})
		return
	}

	// Выполняем usecase
	histogram, err := r.rrIntervalsUsecase.GetRRDifferentialHistogram(c.Request.Context(), from, to, binsCount)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	// Преобразуем usecase ответ в HTTP DTO
	responseDTO := GetDifferentialHistogramResponseDTO{
		DifferentialHistogramDataDTO: toDifferentialHistogramDataDTO(histogram),
	}

	c.JSON(http.StatusOK, responseDTO)
} 