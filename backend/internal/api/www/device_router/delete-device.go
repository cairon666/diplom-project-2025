package device_router

import (
	"net/http"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/usecases/device_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (dr *DeviceRouter) DeleteDevice(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		www.HandleError(c, apperrors.InvalidParams())
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		www.HandleError(c, apperrors.InvalidParams())
	}

	dto := device_usecase.NewDeleteDeviceRequest(id)
	if err := dr.deviceUsecase.DeleteDevice(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
