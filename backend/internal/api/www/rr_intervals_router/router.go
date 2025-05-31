package rr_intervals_router

import (
	"context"
	"time"

	"github.com/cairon666/vkr-backend/internal/indentity"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/usecases/rr_intervals_usecase"
	"github.com/gin-gonic/gin"
)

// TimeRangeDTO представляет временной диапазон для HTTP ответа
type TimeRangeDTO struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type RRIntervalsUsecase interface {
	CreateBatchRRIntervals(ctx context.Context, req rr_intervals_usecase.CreateBatchRRIntervalsRequest) (rr_intervals_usecase.CreateBatchRRIntervalsResponse, error)
	GetRRIntervals(ctx context.Context, req rr_intervals_usecase.GetRRIntervalsRequest) (rr_intervals_usecase.GetRRIntervalsResponse, error)
	
	// Аналитические методы
	GetRRStatistics(ctx context.Context, req rr_intervals_usecase.GetRRStatisticsRequest) (rr_intervals_usecase.GetRRStatisticsResponse, error)
	GetAggregatedRRData(ctx context.Context, from, to time.Time, intervalMinutes int) ([]models.AggregatedRRData, error)
	
	// Отдельные компоненты анализа
	GetRRHistogram(ctx context.Context, from, to time.Time, binsCount int) (*models.RRHistogramData, error)
	GetRRDifferentialHistogram(ctx context.Context, from, to time.Time, binsCount int) (*models.DifferentialHistogramData, error)
	GetRRScatterplot(ctx context.Context, from, to time.Time) (*models.ScatterplotData, error)
	GetRRTrends(ctx context.Context, from, to time.Time, windowSize int) (*models.RRTrendAnalysis, error)
	GetHRVMetrics(ctx context.Context, from, to time.Time) (*models.HRVMetrics, error)
	
	// Комплексный анализ (новый оптимизированный метод)
	GetCompleteAnalysis(ctx context.Context, req rr_intervals_usecase.GetCompleteAnalysisRequest) (rr_intervals_usecase.GetCompleteAnalysisResponse, error)
}

type RRIntervalsRouter struct {
	rrIntervalsUsecase RRIntervalsUsecase
	identityService    *indentity.IdentityService
}

func NewRRIntervalsRouter(rrIntervalsUsecase RRIntervalsUsecase, identityService *indentity.IdentityService) *RRIntervalsRouter {
	return &RRIntervalsRouter{
		rrIntervalsUsecase: rrIntervalsUsecase,
		identityService:    identityService,
	}
}

func (r *RRIntervalsRouter) Register(router gin.IRouter) {
	group := router.Group("/v1")
	group.Use(r.identityService.AuthMiddleware())

	// Базовые операции с R-R интервалами
	group.POST("/rr-intervals/batch", r.CreateBatchRRIntervals)
	group.GET("/rr-intervals", r.GetRRIntervals)
	
	// Аналитические endpoints
	analyticsGroup := group.Group("/rr-intervals/analytics")
	
	// Статистика R-R интервалов
	analyticsGroup.GET("/statistics", r.GetRRStatistics)
	
	// Агрегированные данные
	analyticsGroup.GET("/aggregated", r.GetAggregatedData)
	
	// Отдельные компоненты анализа
	analyticsGroup.GET("/histogram", r.GetHistogram)
	analyticsGroup.GET("/differential-histogram", r.GetDifferentialHistogram)
	analyticsGroup.GET("/scatterplot", r.GetScatterplot)
	analyticsGroup.GET("/trends", r.GetTrends)
	analyticsGroup.GET("/hrv", r.GetHRVMetrics)
	
	// Комплексный анализ (новый оптимизированный endpoint)
	analyticsGroup.GET("/complete", r.GetCompleteAnalysis)
} 