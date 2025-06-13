package external_apps_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/external_apps_usecase"
	"github.com/gin-gonic/gin"
)

type CreateExternalAppRequest struct {
	Name  string   `binding:"required" json:"name"`
	Roles []string `binding:"required" json:"roles"`
}

type CreateExternalAppResponse struct {
	ApiKey        string `json:"api_key"`
	IdExternalApp string `json:"id_external_app"`
}

func (r *ExternalAppsRouter) CreateExternalApp(ctx *gin.Context) {
	var req CreateExternalAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		www.HandleError(ctx, err)

		return
	}

	dto := external_apps_usecase.NewCreateExternalAppRequest(req.Name, req.Roles)
	resp, err := r.externalAppsUsecase.CreateExternalApp(ctx.Request.Context(), dto)
	if err != nil {
		www.HandleError(ctx, err)

		return
	}

	ctx.JSON(http.StatusCreated, CreateExternalAppResponse{
		ApiKey:        resp.ApiKey,
		IdExternalApp: resp.IdExternalApp.String(),
	})
}
