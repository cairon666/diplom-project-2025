package health_router

import (
	"net/http"
	"time"

	"github.com/cairon666/vkr-backend/internal/api/www"
	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/usecases/health_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Step struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	DeviceID  uuid.UUID `json:"device_id"`
	StepCount int64     `json:"step_count"`
	CreatedAt time.Time `json:"created_at"`
}

type GetStepsResponse struct {
	Steps []Step `json:"steps"`
}

func (hr *HealthRouter) GetSteps(c *gin.Context) {
	fromQ := c.Query("from")
	toQ := c.Query("to")

	from, err := time.Parse("2006-01-02", fromQ)
	if err != nil {
		www.HandleError(c, apperrors.ErrInvalidParams)
		return
	}

	to, err := time.Parse("2006-01-02", toQ)
	if err != nil {
		www.HandleError(c, apperrors.ErrInvalidParams)
		return
	}

	dto := health_usecase.NewGetStepsRequest(from, to)
	resp, err := hr.healthUsecase.GetSteps(c.Request.Context(), dto)
	if err != nil {
		www.HandleError(c, err)
		return
	}

	steps := make([]Step, len(resp.Steps))
	for i, v := range resp.Steps {
		steps[i] = Step{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			StepCount: v.StepCount,
			CreatedAt: v.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, GetStepsResponse{
		Steps: steps,
	})
}

type CreateStepsRequest struct {
	Steps []Step `json:"steps"`
}

func (hr *HealthRouter) CreateSteps(c *gin.Context) {
	var req CreateStepsRequest
	if err := c.BindJSON(&req); err != nil {
		www.HandleError(c, err)
		return
	}

	steps := make([]models.Step, len(req.Steps))
	for i, v := range req.Steps {
		steps[i] = models.Step{
			ID:        v.ID,
			UserID:    v.UserID,
			DeviceID:  v.DeviceID,
			StepCount: v.StepCount,
			CreatedAt: v.CreatedAt,
		}
	}

	dto := health_usecase.NewCreateStepsRequest(steps)
	if err := hr.healthUsecase.CreateSteps(c.Request.Context(), dto); err != nil {
		www.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
