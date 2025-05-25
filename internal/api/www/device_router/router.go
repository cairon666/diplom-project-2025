package device_router

import (
	"context"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/usecases/device_usecase"
	"github.com/gin-gonic/gin"
)

type DeviceUsecase interface {
	CreateDevice(ctx context.Context, dto device_usecase.CreateDeviceRequest) (device_usecase.CreateDeviceResponse, error)
	GetDeviceList(ctx context.Context) (device_usecase.GetDeviceListResponse, error)
	DeleteDevice(ctx context.Context, req device_usecase.DeleteDeviceRequest) error
}

type DeviceRouter struct {
	deviceUsecase   DeviceUsecase
	identityService *indentity.IdentityService
}

func NewDeviceRouter(deviceUsecase DeviceUsecase, identityService *indentity.IdentityService) *DeviceRouter {
	return &DeviceRouter{
		deviceUsecase:   deviceUsecase,
		identityService: identityService,
	}
}

func (r *DeviceRouter) Register(router gin.IRouter) {
	group := router.Group("/v1")
	group.Use(r.identityService.AuthMiddleware())

	group.GET("/user/devices", r.GetDeviceList)
	group.POST("/user/devices", r.CreateDevice)
	//group.GET("/user/external-apps/:id", nil)
	//group.PATCH("/user/external-apps/:id", nil)
	group.DELETE("/user/devices/:id", r.DeleteDevice)
}
