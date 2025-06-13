package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
)

type UpdateSettingTelegramRequest struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int    `json:"auth_date"`
	Hash      string `json:"hash"`
}

func (r *UserRouter) UpdateSettingTelegram(c *gin.Context) {
	var req UpdateSettingTelegramRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)

		return
	}

	dto := user_usecase.NewUpdateSettingTelegramRequest(req.Id, req.FirstName, req.LastName, req.Username, req.PhotoURL, req.AuthDate, req.Hash)
	if err := r.userUsecase.UpdateSettingTelegram(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)

		return
	}

	c.Status(http.StatusOK)
}
