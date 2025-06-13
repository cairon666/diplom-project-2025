package device_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/usecases/device_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateDeviceRequest struct {
	DeviceName string `json:"device_name"`
}

type CreateDeviceResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	DeviceName string    `json:"device_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func (dr *DeviceRouter) CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.BindJSON(&req); err != nil {
		www.HandleError(c, err)

		return
	}

	dto := device_usecase.NewCreateDeviceRequest(req.DeviceName)
	resp, err := dr.deviceUsecase.CreateDevice(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, CreateDeviceResponse{
		ID:         resp.ID,
		UserID:     resp.UserID,
		DeviceName: resp.DeviceName,
		CreatedAt:  resp.CreatedAt,
	})
}
