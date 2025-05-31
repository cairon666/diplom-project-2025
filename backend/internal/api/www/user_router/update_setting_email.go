package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
)

type UpdateSettingEmailRequest struct {
	Email string `json:"email"`
}

func (r *UserRouter) UpdateSettingEmail(c *gin.Context) {
	var req UpdateSettingEmailRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)
		return
	}

	dto := user_usecase.NewUpdateSettingEmailRequest(req.Email)
	if err := r.userUsecase.UpdateSettingEmail(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
