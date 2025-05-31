"""
Валидация сгенерированных R-R данных
"""

from typing import Dict, List, Tuple, Optional
from dataclasses import dataclass
from datetime import datetime, timedelta
import statistics

from .data_generator import RRDataPoint


@dataclass
class ValidationResult:
    """Результат валидации"""
    is_valid: bool
    score: float  # 0.0 - 1.0
    issues: List[str]
    statistics: Dict
    recommendations: List[str]


class DataValidator:
    """Валидация качества сгенерированных данных"""
    
    def __init__(self, config):
        self.config = config
        self.thresholds = self._load_validation_thresholds()
    
    def validate_all_datasets(self, datasets: Dict[str, List[RRDataPoint]]) -> Dict:
        """Валидирует все наборы данных"""
        results = {}
        
        for dataset_name, data_points in datasets.items():
            results[dataset_name] = self.validate_dataset(data_points)
        
        # Общая сводка
        results['summary'] = self._create_validation_summary(results)
        
        return results
    
    def validate_dataset(self, data_points: List[RRDataPoint]) -> ValidationResult:
        """Валидирует один набор данных"""
        if not data_points:
            return ValidationResult(
                is_valid=False,
                score=0.0,
                issues=["Пустой набор данных"],
                statistics={},
                recommendations=["Проверьте параметры генерации"]
            )
        
        # Выполняем различные проверки
        issues = []
        statistics = self._calculate_statistics(data_points)
        
        # 1. Проверка диапазона R-R интервалов
        rr_issues = self._validate_rr_range(data_points)
        issues.extend(rr_issues)
        
        # 2. Проверка временной последовательности
        time_issues = self._validate_time_sequence(data_points)
        issues.extend(time_issues)
        
        # 3. Проверка качества данных
        quality_issues = self._validate_data_quality(data_points)
        issues.extend(quality_issues)
        
        # 4. Проверка медицинской достоверности
        medical_issues = self._validate_medical_plausibility(data_points)
        issues.extend(medical_issues)
        
        # Рассчитываем общий score
        score = self._calculate_validation_score(data_points, issues)
        
        # Генерируем рекомендации
        recommendations = self._generate_recommendations(issues, statistics)
        
        return ValidationResult(
            is_valid=len(issues) == 0,
            score=score,
            issues=issues,
            statistics=statistics,
            recommendations=recommendations
        )
    
    def _validate_rr_range(self, data_points: List[RRDataPoint]) -> List[str]:
        """Проверяет диапазон R-R интервалов"""
        issues = []
        rr_values = [dp.rr_interval_ms for dp in data_points]
        
        # Физиологические пределы
        min_rr = min(rr_values)
        max_rr = max(rr_values)
        
        if min_rr < 200:  # > 300 BPM
            issues.append(f"Слишком короткие R-R интервалы: {min_rr:.1f}ms")
        
        if max_rr > 3000:  # < 20 BPM
            issues.append(f"Слишком длинные R-R интервалы: {max_rr:.1f}ms")
        
        # Проверка выбросов
        mean_rr = statistics.mean(rr_values)
        std_rr = statistics.stdev(rr_values) if len(rr_values) > 1 else 0
        
        outliers = [rr for rr in rr_values if abs(rr - mean_rr) > 3 * std_rr]
        if outliers:
            outlier_ratio = len(outliers) / len(rr_values)
            if outlier_ratio > 0.05:  # Более 5% выбросов
                issues.append(f"Много выбросов: {outlier_ratio:.1%}")
        
        return issues
    
    def _validate_time_sequence(self, data_points: List[RRDataPoint]) -> List[str]:
        """Проверяет временную последовательность"""
        issues = []
        
        if len(data_points) < 2:
            return issues
        
        # Проверяем монотонность времени
        for i in range(1, len(data_points)):
            if data_points[i].timestamp <= data_points[i-1].timestamp:
                issues.append(f"Нарушена временная последовательность в позиции {i}")
                break
        
        # Проверяем соответствие временных меток и R-R интервалов
        time_gaps = []
        for i in range(1, len(data_points)):
            time_diff = (data_points[i].timestamp - data_points[i-1].timestamp).total_seconds() * 1000
            expected_gap = data_points[i-1].rr_interval_ms
            
            # Допускаем 10% расхождение
            if abs(time_diff - expected_gap) > expected_gap * 0.1:
                time_gaps.append(i)
        
        if time_gaps:
            gap_ratio = len(time_gaps) / len(data_points)
            if gap_ratio > 0.1:  # Более 10% несоответствий
                issues.append(f"Несоответствие времени и R-R интервалов: {gap_ratio:.1%}")
        
        return issues
    
    def _validate_data_quality(self, data_points: List[RRDataPoint]) -> List[str]:
        """Проверяет качество данных"""
        issues = []
        
        # Проверяем средний quality_score
        quality_scores = [dp.quality_score for dp in data_points]
        avg_quality = statistics.mean(quality_scores)
        
        if avg_quality < 0.5:
            issues.append(f"Низкое среднее качество данных: {avg_quality:.2f}")
        
        # Проверяем долю аномалий
        anomalies = [dp for dp in data_points if dp.is_anomaly]
        anomaly_ratio = len(anomalies) / len(data_points)
        
        if anomaly_ratio > 0.5:
            issues.append(f"Слишком много аномалий: {anomaly_ratio:.1%}")
        
        return issues
    
    def _validate_medical_plausibility(self, data_points: List[RRDataPoint]) -> List[str]:
        """Проверяет медицинскую достоверность"""
        issues = []
        
        rr_values = [dp.rr_interval_ms for dp in data_points]
        
        # Вариабельность сердечного ритма
        if len(rr_values) > 10:
            std_rr = statistics.stdev(rr_values)
            mean_rr = statistics.mean(rr_values)
            coefficient_of_variation = std_rr / mean_rr
            
            # Нормальная вариабельность 2-15%
            if coefficient_of_variation < 0.01:
                issues.append("Слишком низкая вариабельность (возможно артефакт)")
            elif coefficient_of_variation > 0.3:
                issues.append("Слишком высокая вариабельность")
        
        # Проверка на фибрилляцию предсердий
        consecutive_differences = []
        for i in range(1, len(rr_values)):
            diff = abs(rr_values[i] - rr_values[i-1])
            consecutive_differences.append(diff)
        
        if consecutive_differences:
            mean_diff = statistics.mean(consecutive_differences)
            if mean_diff > 100:  # Более 100ms между соседними интервалами
                issues.append("Возможная фибрилляция предсердий")
        
        return issues
    
    def _calculate_statistics(self, data_points: List[RRDataPoint]) -> Dict:
        """Рассчитывает статистики набора данных"""
        if not data_points:
            return {}
        
        rr_values = [dp.rr_interval_ms for dp in data_points]
        quality_scores = [dp.quality_score for dp in data_points]
        
        # Базовые статистики
        stats = {
            "count": len(data_points),
            "duration_minutes": (data_points[-1].timestamp - data_points[0].timestamp).total_seconds() / 60,
            "rr_mean": statistics.mean(rr_values),
            "rr_std": statistics.stdev(rr_values) if len(rr_values) > 1 else 0,
            "rr_min": min(rr_values),
            "rr_max": max(rr_values),
            "quality_mean": statistics.mean(quality_scores),
            "anomaly_count": sum(1 for dp in data_points if dp.is_anomaly),
        }
        
        # Производные метрики
        stats["hr_mean"] = 60000 / stats["rr_mean"]
        stats["hr_range"] = (60000 / stats["rr_max"], 60000 / stats["rr_min"])
        stats["anomaly_ratio"] = stats["anomaly_count"] / stats["count"]
        stats["coefficient_of_variation"] = stats["rr_std"] / stats["rr_mean"]
        
        # HRV метрики
        if len(rr_values) > 2:
            consecutive_diffs = [abs(rr_values[i] - rr_values[i-1]) for i in range(1, len(rr_values))]
            stats["rmssd"] = (statistics.mean([d**2 for d in consecutive_diffs]))**0.5
            stats["pnn50"] = sum(1 for d in consecutive_diffs if d > 50) / len(consecutive_diffs)
        
        return stats
    
    def _calculate_validation_score(self, data_points: List[RRDataPoint], issues: List[str]) -> float:
        """Рассчитывает общий балл валидации"""
        if not data_points:
            return 0.0
        
        # Базовый балл
        score = 1.0
        
        # Штрафы за проблемы
        for issue in issues:
            if "Слишком" in issue or "выбросов" in issue:
                score -= 0.2
            elif "Нарушена" in issue or "Несоответствие" in issue:
                score -= 0.3
            else:
                score -= 0.1
        
        # Бонус за качество
        avg_quality = statistics.mean([dp.quality_score for dp in data_points])
        score *= avg_quality
        
        return max(0.0, min(1.0, score))
    
    def _generate_recommendations(self, issues: List[str], statistics: Dict) -> List[str]:
        """Генерирует рекомендации по улучшению"""
        recommendations = []
        
        if "выбросов" in str(issues):
            recommendations.append("Используйте фильтрацию выбросов")
        
        if "Низкое" in str(issues) and "качество" in str(issues):
            recommendations.append("Увеличьте качество генерируемых данных")
        
        if "вариабельность" in str(issues):
            recommendations.append("Проверьте параметры HRV генерации")
        
        if not recommendations:
            recommendations.append("Данные выглядят корректно")
        
        return recommendations
    
    def _create_validation_summary(self, results: Dict) -> Dict:
        """Создает общую сводку валидации"""
        dataset_results = {k: v for k, v in results.items() if k != 'summary'}
        
        if not dataset_results:
            return {"status": "no_data"}
        
        total_datasets = len(dataset_results)
        valid_datasets = sum(1 for r in dataset_results.values() if r.is_valid)
        avg_score = statistics.mean([r.score for r in dataset_results.values()])
        
        all_issues = []
        for r in dataset_results.values():
            all_issues.extend(r.issues)
        
        return {
            "total_datasets": total_datasets,
            "valid_datasets": valid_datasets,
            "validation_rate": valid_datasets / total_datasets,
            "average_score": avg_score,
            "total_issues": len(all_issues),
            "common_issues": list(set(all_issues)),
            "status": "passed" if valid_datasets == total_datasets else "warnings"
        }
    
    def _load_validation_thresholds(self) -> Dict:
        """Загружает пороговые значения для валидации"""
        return {
            "min_rr_ms": 200,
            "max_rr_ms": 3000,
            "max_outlier_ratio": 0.05,
            "min_quality_score": 0.5,
            "max_anomaly_ratio": 0.5,
            "min_coefficient_variation": 0.01,
            "max_coefficient_variation": 0.3
        } 