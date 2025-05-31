package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

type GetSettingResponse struct {
	Email       string `json:"email"`
	HasPassword bool   `json:"has_password"`
	HasTelegram bool   `json:"has_telegram"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

func (r *UserRouter) GetSetting(c *gin.Context) {
	resp, err := r.userUsecase.GetSetting(c.Request.Context())
	if err != nil {
		www.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetSettingResponse{
		Email:       resp.Email,
		HasPassword: resp.HasPassword,
		HasTelegram: resp.HasTelegram,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
	})
}
