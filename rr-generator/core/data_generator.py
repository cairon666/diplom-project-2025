"""
Основной генератор R-R интервалов с медицинской точностью
"""

import numpy as np
import random
from datetime import datetime, timedelta
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass
import math

from .medical_conditions import MedicalConditions, HeartRhythmState
from utils.time_utils import TimeUtils


@dataclass
class RRDataPoint:
    """Одна точка данных R-R интервала"""
    timestamp: datetime
    rr_interval_ms: float
    user_id: str
    device_id: str
    quality_score: float = 1.0  # 0.0 = плохие данные, 1.0 = идеальные
    is_anomaly: bool = False
    source_condition: Optional[str] = None


class UnifiedDataGenerator:
    """
    Единый генератор всех типов R-R данных
    Может создавать как хорошие, так и плохие данные
    """
    
    def __init__(self, config):
        self.config = config
        self.medical_conditions = MedicalConditions()
        self.noise_generators = self._initialize_noise_generators()
        
    def generate_scenario_data(self, scenario: Dict) -> Dict[str, List[RRDataPoint]]:
        """
        Генерирует данные для полного сценария
        
        Returns:
            Dict с ключами: 'normal', 'bad', 'medical_conditions'
        """
        total_duration = scenario['end_time'] - scenario['start_time']
        good_ratio, bad_ratio = scenario['good_bad_ratio']
        
        datasets = {}
        
        # Приоритет: если задано медицинское состояние
        if scenario.get('condition'):
            datasets['medical_condition'] = self._generate_medical_condition(
                scenario['condition'],
                scenario['start_time'],
                total_duration,
                scenario['user_id'],
                scenario['device_id']
            )
            # Если также нужны плохие данные, добавляем их
            if bad_ratio > 0:
                # Генерируем смешанные данные с медицинским состоянием и плохими данными
                datasets['mixed'] = self._generate_mixed_data_realistic(
                    scenario['start_time'],
                    total_duration,
                    good_ratio,
                    bad_ratio,
                    scenario['user_id'],
                    scenario['device_id']
                )
        # Если нет медицинского состояния, генерируем обычные данные
        elif bad_ratio == 0:
            # Если нужно только хорошие данные
            datasets['normal'] = self._generate_normal_segment(
                scenario['start_time'],
                total_duration,
                scenario['user_id'],
                scenario['device_id']
            )
        else:
            # Генерируем смешанные данные с равномерным распределением плохих участков
            datasets['mixed'] = self._generate_mixed_data_realistic(
                scenario['start_time'],
                total_duration,
                good_ratio,
                bad_ratio,
                scenario['user_id'],
                scenario['device_id']
            )
        
        # Стресс-тест
        if scenario.get('type') == 'stress_test':
            datasets.update(self._generate_stress_test_data(scenario))
        
        return datasets
    
    def _generate_normal_segment(self, start_time: datetime, duration: timedelta, 
                               user_id: str, device_id: str) -> List[RRDataPoint]:
        """Генерирует нормальные (хорошие) R-R интервалы"""
        
        points = []
        current_time = start_time
        end_time = start_time + duration
        
        # Базовые параметры здорового сердца
        base_hr = random.randint(60, 80)  # Здоровый пульс в покое
        base_rr = 60000 / base_hr  # мс
        
        while current_time < end_time:
            # Добавляем естественную вариабельность (HRV)
            hrv_variation = np.random.normal(0, base_rr * 0.04)  # ~4% стандартное отклонение
            
            # Добавляем циркадные ритмы
            circadian_factor = self._get_circadian_factor(current_time)
            
            # Добавляем дыхательную аритмию
            respiratory_variation = self._get_respiratory_variation(current_time)
            
            rr_interval = base_rr + hrv_variation + circadian_factor + respiratory_variation
            
            # Ограничиваем в разумных пределах
            rr_interval = max(400, min(1500, rr_interval))  # 40-150 BPM
            
            points.append(RRDataPoint(
                timestamp=current_time,
                rr_interval_ms=rr_interval,
                user_id=user_id,
                device_id=device_id,
                quality_score=random.uniform(0.9, 1.0),  # Высокое качество
                is_anomaly=False,
                source_condition="normal"
            ))
            
            # Переходим к следующей точке
            current_time += timedelta(milliseconds=rr_interval)
            
        return points
    
    def _generate_bad_segment(self, start_time: datetime, duration: timedelta,
                            user_id: str, device_id: str) -> List[RRDataPoint]:
        """Генерирует плохие/ошибочные данные"""
        
        points = []
        current_time = start_time
        end_time = start_time + duration
        
        bad_data_types = [
            'missing_beats',      # Пропущенные удары
            'artifacts',          # Артефакты движения
            'noise',             # Электрические помехи
            'sensor_drift',      # Дрейф сенсора
            'double_detection',  # Двойное детектирование
            'extreme_outliers'   # Экстремальные выбросы
        ]
        
        while current_time < end_time:
            bad_type = random.choice(bad_data_types)
            
            if bad_type == 'missing_beats':
                # Пропускаем несколько ударов
                skip_duration = random.randint(2000, 5000)  # 2-5 секунд
                current_time += timedelta(milliseconds=skip_duration)
                continue
                
            elif bad_type == 'artifacts':
                # Генерируем артефакты движения
                rr_interval = random.uniform(100, 2500)  # Очень широкий диапазон
                quality_score = random.uniform(0.0, 0.3)
                
            elif bad_type == 'noise':
                # Электрические помехи
                base_rr = 800
                noise = np.random.normal(0, base_rr * 0.5)  # Сильный шум
                rr_interval = max(50, base_rr + noise)
                quality_score = random.uniform(0.1, 0.4)
                
            elif bad_type == 'sensor_drift':
                # Постепенный дрейф сенсора
                drift_factor = 1 + (current_time - start_time).total_seconds() * 0.0001
                rr_interval = 800 * drift_factor
                quality_score = random.uniform(0.3, 0.6)
                
            elif bad_type == 'double_detection':
                # Двойное детектирование R-пика
                base_rr = 800
                rr_interval = base_rr / 2  # Половина от нормального
                quality_score = random.uniform(0.2, 0.5)
                
            else:  # extreme_outliers
                # Экстремальные выбросы
                if random.random() < 0.5:
                    rr_interval = random.uniform(50, 200)  # Очень быстро
                else:
                    rr_interval = random.uniform(3000, 10000)  # Очень медленно
                quality_score = random.uniform(0.0, 0.1)
            
            points.append(RRDataPoint(
                timestamp=current_time,
                rr_interval_ms=rr_interval,
                user_id=user_id,
                device_id=device_id,
                quality_score=quality_score,
                is_anomaly=True,
                source_condition=f"bad_data_{bad_type}"
            ))
            
            current_time += timedelta(milliseconds=max(100, rr_interval))
            
        return points
    
    def _generate_medical_condition(self, condition: str, start_time: datetime, 
                                  duration: timedelta, user_id: str, device_id: str) -> List[RRDataPoint]:
        """Генерирует данные для конкретного медицинского состояния"""
        
        condition_params = self.medical_conditions.get_condition_parameters(condition)
        if not condition_params:
            raise ValueError(f"Неизвестное медицинское состояние: {condition}")
        
        points = []
        current_time = start_time
        end_time = start_time + duration
        
        # Получаем состояние ритма для данного условия
        rhythm_state = condition_params.rhythm_state
        
        while current_time < end_time:
            rr_interval = self._generate_pathological_rr(rhythm_state, current_time, condition_params)
            
            points.append(RRDataPoint(
                timestamp=current_time,
                rr_interval_ms=rr_interval,
                user_id=user_id,
                device_id=device_id,
                quality_score=random.uniform(0.7, 0.9),  # Хорошее качество, но патология
                is_anomaly=condition_params.is_anomaly,
                source_condition=condition
            ))
            
            current_time += timedelta(milliseconds=rr_interval)
            
        return points
    
    def _generate_stress_test_data(self, scenario: Dict) -> Dict[str, List[RRDataPoint]]:
        """Генерирует большие объемы данных для стресс-тестирования"""
        
        datasets = {}
        users_count = scenario.get('users_count', 10)
        
        # Генерируем данные для множественных пользователей
        for i in range(users_count):
            user_id = f"stress_user_{i:03d}"
            device_id = f"device_{i:03d}"
            
            # Смешиваем разные типы данных
            normal_data = self._generate_normal_segment(
                scenario['start_time'],
                timedelta(hours=8),  # 8 часов нормальных данных
                user_id,
                device_id
            )
            
            bad_data = self._generate_bad_segment(
                scenario['start_time'] + timedelta(hours=8),
                timedelta(hours=2),  # 2 часа плохих данных
                user_id,
                device_id
            )
            
            datasets[f'stress_normal_{i}'] = normal_data
            datasets[f'stress_bad_{i}'] = bad_data
            
        return datasets
    
    def _get_circadian_factor(self, timestamp: datetime) -> float:
        """Возвращает циркадную модуляцию ЧСС"""
        hour = timestamp.hour + timestamp.minute / 60.0
        
        # Минимум около 4 утра, максимум около 16:00
        circadian_phase = (hour - 4) * 2 * math.pi / 24
        factor = -50 * math.cos(circadian_phase)  # ±50 мс
        
        return factor
    
    def _get_respiratory_variation(self, timestamp: datetime) -> float:
        """Генерирует дыхательную синусовую аритмию"""
        # Дыхательный цикл ~15 дыханий в минуту
        respiratory_phase = timestamp.second * 2 * math.pi / 4  # 4-секундный цикл
        variation = 20 * math.sin(respiratory_phase)  # ±20 мс
        
        return variation
    
    def _generate_pathological_rr(self, rhythm_state: HeartRhythmState, timestamp: datetime, condition_params=None) -> float:
        """Генерирует патологический R-R интервал"""
        
        if rhythm_state == HeartRhythmState.ATRIAL_FIBRILLATION:
            # Мерцательная аритмия - очень нерегулярный ритм
            if condition_params:
                hr_min, hr_max = condition_params.hr_range
                base_rr = random.uniform(60000/hr_max, 60000/hr_min)
            else:
                base_rr = random.uniform(400, 1200)  # Fallback
            irregular_factor = random.uniform(0.5, 2.0)  # Очень нерегулярно
            result = base_rr * irregular_factor
            return result
            
        elif rhythm_state == HeartRhythmState.TACHYCARDIA:
            # Тахикардия - быстрый но регулярный ритм
            if condition_params:
                hr_min, hr_max = condition_params.hr_range
                hr = random.uniform(hr_min, hr_max)
            else:
                hr = random.uniform(100, 180)  # Fallback
            return 60000 / hr
            
        elif rhythm_state == HeartRhythmState.BRADYCARDIA:
            # Брадикардия - медленный ритм (используем специфические параметры!)
            if condition_params:
                hr_min, hr_max = condition_params.hr_range
                hr = random.uniform(hr_min, hr_max)
            else:
                hr = random.uniform(40, 59)  # Fallback
            return 60000 / hr
            
        elif rhythm_state == HeartRhythmState.IRREGULAR:
            # Нерегулярный ритм
            base_rr = 800
            irregularity = random.uniform(0.6, 1.6)
            return base_rr * irregularity
            
        elif rhythm_state == HeartRhythmState.PREMATURE_BEATS:
            # Преждевременные сокращения - более выраженные
            if random.random() < 0.35:  # 35% экстрасистол (было 20%)
                return random.uniform(250, 500)  # Очень короткий интервал (экстрасистола)
            else:
                if condition_params:
                    hr_min, hr_max = condition_params.hr_range
                    hr = random.uniform(hr_min, hr_max)
                    base_rr = 60000 / hr
                else:
                    base_rr = random.uniform(600, 1000)
                
                # Компенсаторная пауза после экстрасистолы (чаще и длиннее)
                if random.random() < 0.25:  # 25% длинных пауз (было 10%)
                    return base_rr * random.uniform(1.8, 3.0)  # Более длинные паузы (было 1.5-2.5)
                return base_rr
            
        else:  # NORMAL
            # Нормальный ритм с естественной вариабельностью
            if condition_params:
                hr_min, hr_max = condition_params.hr_range
                hr = random.uniform(hr_min, hr_max)
                base_rr = 60000 / hr
            else:
                base_rr = random.uniform(600, 1000)  # 60-100 BPM
            hrv = np.random.normal(0, base_rr * 0.04)
            return base_rr + hrv
    
    def _initialize_noise_generators(self) -> Dict:
        """Инициализирует генераторы различных типов шума"""
        return {
            'gaussian': lambda sigma: np.random.normal(0, sigma),
            'uniform': lambda range_val: random.uniform(-range_val, range_val),
            'spike': lambda: random.uniform(100, 2000) if random.random() < 0.01 else 0,
            'drift': lambda t: t * 0.001  # Линейный дрейф
        }
    
    def preview_generation(self, scenario: Dict):
        """Показывает предварительный просмотр того, что будет сгенерировано"""
        duration = scenario['end_time'] - scenario['start_time']
        
        print(f"📋 Предварительный просмотр генерации:")
        print(f"   Период: {duration}")
        print(f"   Тип: {scenario.get('type', 'unknown')}")
        print(f"   Пользователей: {scenario.get('users_count', 1)}")
        
        # Приблизительная оценка количества точек
        avg_rr = 800  # мс
        estimated_points = int(duration.total_seconds() * 1000 / avg_rr)
        print(f"   Ожидаемое количество точек: ~{estimated_points:,}")
        
        if scenario.get('condition'):
            print(f"   Медицинское состояние: {scenario['condition']}")
        
        good_ratio, bad_ratio = scenario.get('good_bad_ratio', (70, 30))
        if bad_ratio == 0:
            print(f"   Качество данных: только хорошие данные (невалидные отключены)")
        else:
            print(f"   Соотношение хорошие:плохие = {good_ratio}:{bad_ratio}")
    
    def _generate_mixed_data_realistic(self, start_time: datetime, total_duration: timedelta,
                                     good_ratio: int, bad_ratio: int, user_id: str, device_id: str) -> List[RRDataPoint]:
        """
        Генерирует смешанные данные с реалистичным распределением плохих участков по всему периоду
        """
        points = []
        current_time = start_time
        end_time = start_time + total_duration
        
        # Определяем количество и длительность плохих сегментов
        total_seconds = total_duration.total_seconds()
        bad_total_seconds = total_seconds * (bad_ratio / 100)
        
        # Создаем 3-5 коротких плохих сегментов вместо одного длинного
        num_bad_segments = random.randint(3, 5)
        bad_segment_duration = bad_total_seconds / num_bad_segments
        
        # Генерируем случайные позиции для плохих сегментов
        bad_segments = []
        for i in range(num_bad_segments):
            # Случайное время начала плохого сегмента
            segment_start_offset = random.uniform(0, total_seconds - bad_segment_duration)
            segment_start = start_time + timedelta(seconds=segment_start_offset)
            segment_end = segment_start + timedelta(seconds=bad_segment_duration)
            bad_segments.append((segment_start, segment_end))
        
        # Сортируем сегменты по времени
        bad_segments.sort(key=lambda x: x[0])
        
        # Генерируем данные, переключаясь между хорошими и плохими участками
        while current_time < end_time:
            # Проверяем, находимся ли мы в плохом сегменте
            in_bad_segment = any(start <= current_time < end for start, end in bad_segments)
            
            if in_bad_segment:
                # Генерируем плохие данные
                data_point = self._generate_single_bad_point(current_time, user_id, device_id)
            else:
                # Генерируем хорошие данные
                data_point = self._generate_single_normal_point(current_time, user_id, device_id)
            
            points.append(data_point)
            current_time += timedelta(milliseconds=data_point.rr_interval_ms)
            
        return points
    
    def _generate_single_normal_point(self, timestamp: datetime, user_id: str, device_id: str) -> RRDataPoint:
        """Генерирует одну точку нормальных данных"""
        # Базовые параметры здорового сердца
        base_hr = random.randint(60, 80)  # Здоровый пульс в покое
        base_rr = 60000 / base_hr  # мс
        
        # Добавляем естественную вариабельность (HRV)
        hrv_variation = np.random.normal(0, base_rr * 0.04)  # ~4% стандартное отклонение
        
        # Добавляем циркадные ритмы
        circadian_factor = self._get_circadian_factor(timestamp)
        
        # Добавляем дыхательную аритмию
        respiratory_variation = self._get_respiratory_variation(timestamp)
        
        rr_interval = base_rr + hrv_variation + circadian_factor + respiratory_variation
        
        # Ограничиваем в разумных пределах
        rr_interval = max(400, min(1500, rr_interval))  # 40-150 BPM
        
        return RRDataPoint(
            timestamp=timestamp,
            rr_interval_ms=rr_interval,
            user_id=user_id,
            device_id=device_id,
            quality_score=random.uniform(0.9, 1.0),  # Высокое качество
            is_anomaly=False,
            source_condition="normal"
        )
    
    def _generate_single_bad_point(self, timestamp: datetime, user_id: str, device_id: str) -> RRDataPoint:
        """Генерирует одну точку плохих данных"""
        bad_data_types = [
            'artifacts',          # Артефакты движения
            'noise',             # Электрические помехи
            'double_detection',  # Двойное детектирование
            'extreme_outliers'   # Экстремальные выбросы (реже)
        ]
        
        # Экстремальные выбросы происходят реже
        if random.random() < 0.9:  # 90% - не экстремальные
            bad_type = random.choice(bad_data_types[:-1])  # Исключаем extreme_outliers
        else:
            bad_type = 'extreme_outliers'
        
        if bad_type == 'artifacts':
            # Артефакты движения - умеренные выбросы
            rr_interval = random.uniform(300, 1800)  # Более узкий диапазон
            quality_score = random.uniform(0.4, 0.7)
            
        elif bad_type == 'noise':
            # Электрические помехи
            base_rr = 800
            noise = np.random.normal(0, base_rr * 0.2)  # Меньше шума
            rr_interval = max(200, base_rr + noise)
            quality_score = random.uniform(0.3, 0.6)
            
        elif bad_type == 'double_detection':
            # Двойное детектирование R-пика
            base_rr = 800
            rr_interval = base_rr / 2  # Половина от нормального
            quality_score = random.uniform(0.2, 0.5)
            
        else:  # extreme_outliers - редко
            # Экстремальные выбросы (теперь реже и менее экстремальные)
            if random.random() < 0.5:
                rr_interval = random.uniform(100, 300)  # Быстро, но не экстремально
            else:
                rr_interval = random.uniform(1500, 2500)  # Медленно, но не экстремально
            quality_score = random.uniform(0.0, 0.2)
        
        return RRDataPoint(
            timestamp=timestamp,
            rr_interval_ms=rr_interval,
            user_id=user_id,
            device_id=device_id,
            quality_score=quality_score,
            is_anomaly=True,
            source_condition=f"bad_data_{bad_type}"
        ) 