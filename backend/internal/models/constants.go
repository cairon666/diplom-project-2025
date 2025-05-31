package models

// TrendDirection представляет направление тренда
type TrendDirection string

const (
	TrendDirectionUp     TrendDirection = "up"     // Восходящий тренд
	TrendDirectionDown   TrendDirection = "down"   // Нисходящий тренд
	TrendDirectionStable TrendDirection = "stable" // Стабильный тренд
)

// String возвращает строковое представление направления тренда
func (td TrendDirection) String() string {
	return string(td)
}

// OverallTrend представляет общий тренд
type OverallTrend string

const (
	OverallTrendIncreasing OverallTrend = "increasing" // Возрастающий тренд
	OverallTrendDecreasing OverallTrend = "decreasing" // Убывающий тренд
	OverallTrendStable     OverallTrend = "stable"     // Стабильный тренд
)

// String возвращает строковое представление общего тренда
func (ot OverallTrend) String() string {
	return string(ot)
}
