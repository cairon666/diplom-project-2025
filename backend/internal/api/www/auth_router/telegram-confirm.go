package auth_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TelegramConfirmRequest struct {
	TempId    string `json:"temp_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TelegramConfirmResponse struct {
	AccessToken string `json:"access_token"`
	Id          string `json:"id"`
}

func (r *AuthRouter) TelegramConfirmRoute(c *gin.Context) {
	var req TelegramConfirmRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)

		return
	}

	tempId := uuid.MustParse(req.TempId)
	reqDTO := auth_usecase.NewCompleteRegisterRequest(tempId, req.Email, req.FirstName, req.LastName)
	resp, err := r.authUsecase.TelegramConfirm(c, reqDTO)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	r.jwtService.SetRefreshCookie(c, resp.RefreshToken)
	c.JSON(http.StatusOK, TelegramConfirmResponse{
		AccessToken: resp.AccessToken,
		Id:          resp.UserId.String(),
	})
}
