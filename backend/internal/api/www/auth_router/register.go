package auth_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/auth_usecase"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
}

func (r *AuthRouter) RegisterRoute(c *gin.Context) {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		www.HandleError(c, err)
		return
	}

	dtoReq := auth_usecase.NewRegisterRequest(req.Email, req.Password, req.FirstName, req.SecondName)
	if err := r.authUsecase.Register(c, dtoReq); err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
