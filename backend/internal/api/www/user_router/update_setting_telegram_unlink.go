package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

func (r *UserRouter) UpdateSettingTelegramUnlink(c *gin.Context) {
	if err := r.userUsecase.UpdateSettingTelegramUnlink(c.Request.Context()); err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusOK)
}
