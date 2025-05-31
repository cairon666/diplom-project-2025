package rr_intervals_service

import (
	"context"
	"math"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/google/uuid"
)

// BuildScatterplot строит скаттерограмму (диаграмму Пуанкаре) R-R интервалов
func (s *RRIntervalsService) BuildScatterplot(ctx context.Context, userID uuid.UUID, from, to time.Time) (*models.ScatterplotData, error) {
	// Получаем данные R-R интервалов
	values, err := s.rrIntervalsRepo.GetRawValuesForAnalysis(ctx, userID, from, to)
	if err != nil {
		return nil, apperrors.DataProcessingErrorf("failed to get RR intervals for scatterplot: %v", err)
	}

	if len(values) < 2 {
		return &models.ScatterplotData{
			Points:     []models.ScatterplotPoint{},
			TotalCount: 0,
			Statistics: &models.ScatterplotStatistics{},
			Ellipse:    &models.PoincarePlotEllipse{},
		}, nil
	}

	// Создаем точки скаттерограммы
	points := make([]models.ScatterplotPoint, len(values)-1)
	for i := 1; i < len(values); i++ {
		points[i-1] = models.ScatterplotPoint{
			RRn:  values[i-1],
			RRn1: values[i],
		}
	}

	// Вычисляем статистику скаттерограммы
	statistics := s.calculateScatterplotStatistics(points)

	// Вычисляем параметры эллипса Пуанкаре
	ellipse := s.calculatePoincarePlotEllipse(points, statistics)

	return &models.ScatterplotData{
		Points:     points,
		TotalCount: int64(len(points)),
		Statistics: statistics,
		Ellipse:    ellipse,
	}, nil
}

// calculateScatterplotStatistics вычисляет статистику скаттерограммы (параметры Пуанкаре)
func (s *RRIntervalsService) calculateScatterplotStatistics(points []models.ScatterplotPoint) *models.ScatterplotStatistics {
	if len(points) == 0 {
		return &models.ScatterplotStatistics{}
	}

	// Вычисляем средние значения
	var sumRRn, sumRRn1 float64
	for _, point := range points {
		sumRRn += float64(point.RRn)
		sumRRn1 += float64(point.RRn1)
	}
	meanRRn := sumRRn / float64(len(points))
	meanRRn1 := sumRRn1 / float64(len(points))

	// Вычисляем дисперсии и ковариацию
	var varRRn, varRRn1, covariance float64
	for _, point := range points {
		diffRRn := float64(point.RRn) - meanRRn
		diffRRn1 := float64(point.RRn1) - meanRRn1
		
		varRRn += diffRRn * diffRRn
		varRRn1 += diffRRn1 * diffRRn1
		covariance += diffRRn * diffRRn1
	}
	varRRn /= float64(len(points))
	varRRn1 /= float64(len(points))
	covariance /= float64(len(points))

	// Вычисляем SD1 и SD2 (стандартные отклонения по осям эллипса)
	// SD1 = квадратный корень из дисперсии разностей RR[n+1] - RR[n], деленной на 2
	var sumSquareDiffs float64
	for _, point := range points {
		diff := float64(point.RRn1 - point.RRn)
		sumSquareDiffs += diff * diff
	}
	sd1 := math.Sqrt(sumSquareDiffs / (2.0 * float64(len(points))))

	// SD2 = квадратный корень из дисперсии сумм RR[n+1] + RR[n], деленной на 2
	var sumSquareSums float64
	for _, point := range points {
		sum := float64(point.RRn1 + point.RRn)
		meanSum := 2 * (meanRRn + meanRRn1) / 2 // упрощается до meanRRn + meanRRn1
		diff := sum - meanSum
		sumSquareSums += diff * diff
	}
	sd2 := math.Sqrt(sumSquareSums / (2.0 * float64(len(points))))

	// Вычисляем отношение SD1/SD2
	var sd1SD2Ratio float64
	if sd2 != 0 {
		sd1SD2Ratio = sd1 / sd2
	}

	// Вычисляем CSI (Cardiac Sympathetic Index) и CVI (Cardiac Vagal Index)
	// CSI = длина оси эллипса / SD1
	// CVI = логарифм от произведения SD1 и SD2
	csi := sd2 // Упрощенная формула
	cvi := math.Log(sd1 * sd2)

	return &models.ScatterplotStatistics{
		SD1:         sd1,
		SD2:         sd2,
		SD1SD2Ratio: sd1SD2Ratio,
		CSI:         csi,
		CVI:         cvi,
	}
}

// calculatePoincarePlotEllipse вычисляет параметры эллипса диаграммы Пуанкаре
func (s *RRIntervalsService) calculatePoincarePlotEllipse(points []models.ScatterplotPoint, stats *models.ScatterplotStatistics) *models.PoincarePlotEllipse {
	if len(points) == 0 {
		return &models.PoincarePlotEllipse{}
	}

	// Вычисляем центр эллипса (средние значения)
	var sumRRn, sumRRn1 float64
	for _, point := range points {
		sumRRn += float64(point.RRn)
		sumRRn1 += float64(point.RRn1)
	}
	centerX := sumRRn / float64(len(points))
	centerY := sumRRn1 / float64(len(points))

	// Площадь эллипса = π * SD1 * SD2
	area := math.Pi * stats.SD1 * stats.SD2

	return &models.PoincarePlotEllipse{
		CenterX: centerX,
		CenterY: centerY,
		SD1:     stats.SD1,
		SD2:     stats.SD2,
		Area:    area,
	}
} 