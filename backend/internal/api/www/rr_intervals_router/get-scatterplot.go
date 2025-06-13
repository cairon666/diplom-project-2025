package rr_intervals_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

// GetScatterplotRequestDTO представляет HTTP запрос на получение скаттерограммы.
type GetScatterplotRequestDTO struct {
	From string `binding:"required" form:"from"` // RFC3339 формат
	To   string `binding:"required" form:"to"`   // RFC3339 формат
}

// GetScatterplotResponseDTO представляет HTTP ответ со скаттерограммой.
type GetScatterplotResponseDTO struct {
	*ScatterplotDataDTO
}

func (r *RRIntervalsRouter) GetScatterplot(c *gin.Context) {
	var reqDTO GetScatterplotRequestDTO
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

	// Выполняем usecase
	scatterplot, err := r.rrIntervalsUsecase.GetRRScatterplot(c.Request.Context(), from, to)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	// Преобразуем usecase ответ в HTTP DTO
	responseDTO := GetScatterplotResponseDTO{
		ScatterplotDataDTO: toScatterplotDataDTO(scatterplot),
	}

	c.JSON(http.StatusOK, responseDTO)
}
