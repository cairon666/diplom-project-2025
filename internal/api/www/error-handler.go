package www

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var appErr *apperrors.BaseError
	fmt.Println(err)

	// Ищем первую ошибку типа AppError в цепочке
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPCode(), gin.H{
			"error":   appErr.Code(),
			"message": appErr.Error(), // показываем только бизнес-смысл
		})
		return
	}

	var appErr2 apperrors.AppError
	if errors.As(err, &appErr2) {
		c.JSON(appErr2.HTTPCode(), appErr2.JSON())
		return
	}

	// Если ошибка не является AppError — вернуть 500
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "INTERNAL_SERVER_ERROR",
		"message": err.Error(),
	})
}
