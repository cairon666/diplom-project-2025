package external_apps_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/usecases/external_apps_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *ExternalAppsRouter) DeleteExternalApp(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		www.HandleError(ctx, apperrors.InvalidParams())

		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		www.HandleError(ctx, apperrors.InvalidParams())

		return
	}

	dto := external_apps_usecase.NewDeleteExternalAppRequest(id)
	if err := r.externalAppsUsecase.DeleteExternalApp(ctx.Request.Context(), dto); err != nil {
		www.HandleError(ctx, err)

		return
	}

	ctx.Status(http.StatusNoContent)
}
