package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
)

type UpdateSettingProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *UserRouter) UpdateSettingProfile(c *gin.Context) {
	var req UpdateSettingProfileRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)
		return
	}

	dto := user_usecase.NewUpdateSettingProfileRequest(req.FirstName, req.LastName)
	if err := r.userUsecase.UpdateSettingProfile(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
