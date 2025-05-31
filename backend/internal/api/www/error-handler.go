package www

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var appErr apperrors.AppError
	fmt.Println(err)

	// Проверяем, является ли ошибка нашей кастомной ошибкой
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPCode(), appErr.JSON())
		return
	}

	// Если ошибка не является AppError — вернуть 500
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "INTERNAL_SERVER_ERROR",
		"message": err.Error(),
	})
}
