package external_apps_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/external_apps_usecase"
	"github.com/gin-gonic/gin"
)

type ExternalAppsUsecase interface {
	CreateExternalApp(ctx context.Context, dto external_apps_usecase.CreateExternalAppRequest) (external_apps_usecase.CreateExternalAppResponse, error)
	GetExternalAppList(ctx context.Context) (external_apps_usecase.GetExternalAppListResponse, error)
	DeleteExternalApp(ctx context.Context, dto external_apps_usecase.DeleteExternalAppRequest) error
}

type ExternalAppsRouter struct {
	externalAppsUsecase ExternalAppsUsecase
	identityService     *indentity.IdentityService
}

func NewExternalAppsRouter(externalAppsUsecase ExternalAppsUsecase, identityService *indentity.IdentityService) *ExternalAppsRouter {
	return &ExternalAppsRouter{
		externalAppsUsecase: externalAppsUsecase,
		identityService:     identityService,
	}
}

func (r *ExternalAppsRouter) Register(router gin.IRouter) {
	group := router.Group("/v1")
	group.Use(r.identityService.AuthMiddleware())

	group.GET("/user/external-apps", r.GetExternalAppList)
	group.POST("/user/external-apps", r.CreateExternalApp)
	group.DELETE("/user/external-apps/:id", r.DeleteExternalApp)

	// group.GET("/user/external-apps/:id", nil)
	// group.PATCH("/user/external-apps/:id", nil)
}
