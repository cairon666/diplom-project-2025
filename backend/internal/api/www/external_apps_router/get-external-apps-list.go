package external_apps_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExternalListItem struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Roles     []string  `json:"roles"`
}

type GetExternalAppListResponse struct {
	ExternalApps []ExternalListItem `json:"external_apps"`
}

func (r *ExternalAppsRouter) GetExternalAppList(ctx *gin.Context) {
	resp, err := r.externalAppsUsecase.GetExternalAppList(ctx.Request.Context())
	if err != nil {
		www.HandleError(ctx, err)
		return
	}

	externalApps := make([]ExternalListItem, 0, len(resp.ExternalApps))
	for _, externalApp := range resp.ExternalApps {
		externalApps = append(externalApps, ExternalListItem{
			ID:        externalApp.ID,
			OwnerID:   externalApp.OwnerID,
			Name:      externalApp.Name,
			CreatedAt: externalApp.CreatedAt,
			Roles:     externalApp.Roles,
		})
	}

	ctx.JSON(http.StatusOK, GetExternalAppListResponse{
		ExternalApps: externalApps,
	})
}
