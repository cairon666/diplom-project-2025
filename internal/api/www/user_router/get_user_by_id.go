package user_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
)

type GetUserByIdResponseDTO struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *UserRouter) GetUserByIdRouter(c *gin.Context) {
	user, err := r.userUsecase.GetUserById(c.Request.Context())
	if err != nil {
		www.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetUserByIdResponseDTO{
		Id:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}
