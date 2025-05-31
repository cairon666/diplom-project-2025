package auth_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Id          string `json:"id"`
}

func (r *AuthRouter) LoginRoute(c *gin.Context) {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)
		return
	}

	dtoReq := auth_usecase.NewLoginRequest(req.Email, req.Password)
	resp, err := r.authUsecase.Login(c, dtoReq)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	r.jwtService.SetRefreshCookie(c, resp.RefreshToken)
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: resp.AccessToken,
		Id:          resp.Id,
	})
}
