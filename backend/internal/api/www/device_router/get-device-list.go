package device_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceListItem struct {
	ID         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	DeviceName string    `json:"device_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetDeviceListResponse struct {
	Devices []DeviceListItem `json:"devices"`
}

func (dr *DeviceRouter) GetDeviceList(c *gin.Context) {
	resp, err := dr.deviceUsecase.GetDeviceList(c.Request.Context())
	if err != nil {
		www.HandleError(c, err)
		return
	}

	devices := make([]DeviceListItem, 0, len(resp.Devices))
	for _, device := range resp.Devices {
		devices = append(devices, DeviceListItem{
			ID:         device.ID,
			UserId:     device.UserID,
			DeviceName: device.DeviceName,
			CreatedAt:  device.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, GetDeviceListResponse{
		Devices: devices,
	})
}
