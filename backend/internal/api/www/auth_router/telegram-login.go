package auth_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
)

type TelegramLoginRequest struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int    `json:"auth_date"`
	Hash      string `json:"hash"`
}

type TelegramLoginResponse struct {
	AccessToken string `json:"access_token"`
	Id          string `json:"id"`
}

func (r *AuthRouter) TelegramLoginRoute(c *gin.Context) {
	var req TelegramLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		www.HandleError(c, err)

		return
	}

	reqDTO := auth_usecase.NewTelegramLoginRequest(req.Id, req.FirstName, req.LastName, req.Username, req.PhotoURL, req.AuthDate, req.Hash)
	resp, err := r.authUsecase.TelegramLogin(c, reqDTO)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	r.jwtService.SetRefreshCookie(c, resp.RefreshToken)
	c.JSON(http.StatusOK, TelegramLoginResponse{
		AccessToken: resp.AccessToken,
		Id:          resp.UserId.String(),
	})
}
