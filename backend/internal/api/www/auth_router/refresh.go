package auth_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
)

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (r *AuthRouter) RefreshRoute(c *gin.Context) {
	refreshToken, err := r.jwtService.GetRefreshCookie(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)

		return
	}

	reqDTO := auth_usecase.NewRefreshTokenRequest(refreshToken)
	resp, err := r.authUsecase.RefreshToken(c, reqDTO)
	if err != nil {
		c.Status(http.StatusUnauthorized)

		return
	}

	r.jwtService.SetRefreshCookie(c, resp.RefreshToken)
	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: resp.AccessToken,
	})
}
