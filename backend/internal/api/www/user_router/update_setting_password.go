package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
)

type UpdateSettingPasswordRequest struct {
	Password string `json:"password"`
}

func (r *UserRouter) UpdateSettingPassword(c *gin.Context) {
	var req UpdateSettingPasswordRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)

		return
	}

	dto := user_usecase.NewUpdateSettingPasswordRequest(req.Password)
	if err := r.userUsecase.UpdateSettingPassword(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusOK)
}
