package bot

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// handleRRIntervalsCommand обрабатывает команду /rrintervals
func (b *Bot) handleRRIntervalsCommand(chatID, userID int64, args string) {
	user, exists := b.storage.GetUser(userID)
	if !exists || user.APIKey == "" {
		b.sendMessage(chatID, "❌ API ключ не установлен. Используйте /connect для подключения.")
		return
	}

	from, to, err := b.parseDateRange(args)
	if err != nil {
		b.sendMessage(chatID, "❌ Неверный формат даты. Используйте: /rr 2024-01-01 2024-01-02")
		return
	}

	response, err := b.client.GetRRIntervals(user.APIKey, from, to)
	if err != nil {
		b.logger.Error("Failed to get R-R intervals", zap.Error(err))
		b.sendMessage(chatID, "❌ Ошибка при получении R-R интервалов. Проверьте подключение.")
		return
	}

	if len(response.RRIntervals) == 0 {
		b.sendMessage(chatID, "📊 Нет данных о R-R интервалах за указанный период.")
		return
	}

	message := fmt.Sprintf(`💓 *R-R интервалы*
📅 Период: %s - %s
📊 Всего записей: %d

Первые несколько значений:`,
		from.Format("02.01.2006"),
		to.Format("02.01.2006"),
		response.TotalCount,
	)

	// Показываем первые 10 значений
	count := len(response.RRIntervals)
	if count > 10 {
		count = 10
	}

	for i := 0; i < count; i++ {
		rr := response.RRIntervals[i]
		message += fmt.Sprintf("\n• %d мс (%s)", rr.RRValue, rr.Timestamp.Format("15:04:05"))
	}

	if len(response.RRIntervals) > 10 {
		message += fmt.Sprintf("\n... и еще %d записей", len(response.RRIntervals)-10)
	}

	message += "\n\nИспользуйте /rrstats для статистического анализа."

	b.sendMessage(chatID, message)
}

// handleRRStatisticsCommand обрабатывает команду /rrstats
func (b *Bot) handleRRStatisticsCommand(chatID, userID int64, args string) {
	user, exists := b.storage.GetUser(userID)
	if !exists || user.APIKey == "" {
		b.sendMessage(chatID, "❌ API ключ не установлен. Используйте /connect для подключения.")
		return
	}

	from, to, err := b.parseDateRange(args)
	if err != nil {
		b.sendMessage(chatID, "❌ Неверный формат даты. Используйте: /rrstats 2024-01-01 2024-01-02")
		return
	}

	response, err := b.client.GetRRStatistics(user.APIKey, from, to, true, true, 20)
	if err != nil {
		b.logger.Error("Failed to get R-R statistics", zap.Error(err))
		b.sendMessage(chatID, "❌ Ошибка при получении статистики R-R интервалов.")
		return
	}

	if response.Summary == nil || response.Summary.Count == 0 {
		b.sendMessage(chatID, "📊 Нет данных о R-R интервалах за указанный период.")
		return
	}

	message := fmt.Sprintf(`📈 *Статистика R-R интервалов*
📅 Период: %s - %s

📊 *Основные показатели:*
• Количество: %d измерений
• Среднее: %.1f мс
• Стандартное отклонение: %.1f мс
• Минимум: %d мс
• Максимум: %d мс`,
		from.Format("02.01.2006"),
		to.Format("02.01.2006"),
		response.Summary.Count,
		response.Summary.Mean,
		response.Summary.StdDev,
		response.Summary.Min,
		response.Summary.Max,
	)

	// Добавляем HRV метрики если доступны
	if response.HRVMetrics != nil {
		message += fmt.Sprintf(`

💓 *HRV Метрики:*
• RMSSD: %.1f мс
• SDNN: %.1f мс
• pNN50: %.1f%%
• Треугольный индекс: %.1f
• LF/HF: %.2f`,
			response.HRVMetrics.RMSSD,
			response.HRVMetrics.SDNN,
			response.HRVMetrics.PNN50,
			response.HRVMetrics.TriangularIndex,
			response.HRVMetrics.LFHFRatio,
		)
	}

	message += "\n\nИспользуйте /rranalyze для более подробного анализа."

	b.sendMessage(chatID, message)
}

// handleRRAnalysisCommand обрабатывает команду /rranalyze
func (b *Bot) handleRRAnalysisCommand(chatID, userID int64, args string) {
	user, exists := b.storage.GetUser(userID)
	if !exists || user.APIKey == "" {
		b.sendMessage(chatID, "❌ API ключ не установлен. Используйте /connect для подключения.")
		return
	}

	from, to, err := b.parseDateRange(args)
	if err != nil {
		b.sendMessage(chatID, "❌ Неверный формат даты. Используйте: /rranalyze 2024-01-01 2024-01-02")
		return
	}

	// Получаем статистику с полным анализом
	statsResponse, err := b.client.GetRRStatistics(user.APIKey, from, to, true, true, 20)
	if err != nil {
		b.logger.Error("Failed to get R-R statistics", zap.Error(err))
		b.sendMessage(chatID, "❌ Ошибка при получении статистики R-R интервалов.")
		return
	}

	if statsResponse.Summary == nil || statsResponse.Summary.Count == 0 {
		b.sendMessage(chatID, "📊 Нет данных о R-R интервалах за указанный период.")
		return
	}

	// Получаем данные скаттерплота
	scatterResponse, err := b.client.GetRRScatterplot(user.APIKey, from, to)
	if err != nil {
		b.logger.Warn("Failed to get scatterplot data", zap.Error(err))
	}

	message := fmt.Sprintf(`🔬 *Полный анализ R-R интервалов*
📅 Период: %s - %s

📊 *Основная статистика:*
• Измерений: %d
• Среднее: %.1f мс (%.1f BPM)
• Вариативность: %.1f мс
• Диапазон: %d - %d мс`,
		from.Format("02.01.2006"),
		to.Format("02.01.2006"),
		statsResponse.Summary.Count,
		statsResponse.Summary.Mean,
		60000/statsResponse.Summary.Mean, // преобразуем в BPM
		statsResponse.Summary.StdDev,
		statsResponse.Summary.Min,
		statsResponse.Summary.Max,
	)

	// HRV анализ
	if statsResponse.HRVMetrics != nil {
		hrv := statsResponse.HRVMetrics
		message += fmt.Sprintf(`

💗 *Анализ вариабельности пульса (HRV):*

🔸 *Временные показатели:*
• RMSSD: %.1f мс (краткосрочная вариабельность)
• SDNN: %.1f мс (общая вариабельность)
• pNN50: %.1f%% (парасимпатическая активность)

🔸 *Геометрические показатели:*
• Треугольный индекс: %.1f
• TINN: %.1f мс

🔸 *Частотные показатели:*
• VLF: %.1f мс² (очень низкие частоты)
• LF: %.1f мс² (низкие частоты)
• HF: %.1f мс² (высокие частоты)
• LF/HF: %.2f (симпато/парасимпатический баланс)
• Общая мощность: %.1f мс²`,
			hrv.RMSSD,
			hrv.SDNN,
			hrv.PNN50,
			hrv.TriangularIndex,
			hrv.TINN,
			hrv.VLFPower,
			hrv.LFPower,
			hrv.HFPower,
			hrv.LFHFRatio,
			hrv.TotalPower,
		)

		// Интерпретация результатов
		message += "\n\n🩺 *Интерпретация:*"
		
		if hrv.RMSSD < 20 {
			message += "\n• ⚠️ Низкая краткосрочная вариабельность"
		} else if hrv.RMSSD > 50 {
			message += "\n• ✅ Высокая краткосрочная вариабельность"
		} else {
			message += "\n• ✅ Нормальная краткосрочная вариабельность"
		}

		if hrv.LFHFRatio > 2.5 {
			message += "\n• ⚠️ Повышенный симпатический тонус (стресс)"
		} else if hrv.LFHFRatio < 1.5 {
			message += "\n• ✅ Хороший парасимпатический тонус (восстановление)"
		} else {
			message += "\n• ✅ Сбалансированная активность НС"
		}
	}

	// Анализ скаттерплота
	if scatterResponse != nil && scatterResponse.Statistics != nil {
		stats := scatterResponse.Statistics
		message += fmt.Sprintf(`

📊 *Анализ диаграммы Пуанкаре:*
• SD1: %.1f мс (быстрая изменчивость)
• SD2: %.1f мс (медленная изменчивость) 
• SD1/SD2: %.2f (соотношение)
• CSI: %.1f (индекс сердечного стресса)
• CVI: %.1f (индекс сердечной активности)`,
			stats.SD1,
			stats.SD2,
			stats.SD1SD2Ratio,
			stats.CSI,
			stats.CVI,
		)
	}

	// Гистограмма
	if statsResponse.Histogram != nil && len(statsResponse.Histogram.Bins) > 0 {
		message += "\n\n📈 *Распределение R-R интервалов:*"
		// Находим пик гистограммы
		maxCount := int64(0)
		peakBin := statsResponse.Histogram.Bins[0]
		for _, bin := range statsResponse.Histogram.Bins {
			if bin.Count > maxCount {
				maxCount = bin.Count
				peakBin = bin
			}
		}
		message += fmt.Sprintf("\n• Пик распределения: %d-%d мс (%.1f%%)",
			peakBin.RangeStart, peakBin.RangeEnd, peakBin.Frequency*100)
	}

	b.sendMessage(chatID, message)
}

// parseDateRange парсит диапазон дат из аргументов
func (b *Bot) parseDateRange(args string) (time.Time, time.Time, error) {
	args = strings.TrimSpace(args)
	
	// Если аргументы пустые, используем сегодняшний день
	if args == "" {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		return start, end, nil
	}

	parts := strings.Fields(args)
	
	// Если одна дата, используем её как начало и конец дня
	if len(parts) == 1 {
		date, err := time.Parse("2006-01-02", parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
		return start, end, nil
	}

	// Если две даты
	if len(parts) == 2 {
		from, err := time.Parse("2006-01-02", parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		to, err := time.Parse("2006-01-02", parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		
		// Устанавливаем время начала дня для from и конца дня для to
		start := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		end := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())
		
		// Проверяем правильность порядка дат
		if start.After(end) {
			start, end = end, start
		}
		
		return start, end, nil
	}

	return time.Time{}, time.Time{}, fmt.Errorf("неверный формат аргументов")
} 