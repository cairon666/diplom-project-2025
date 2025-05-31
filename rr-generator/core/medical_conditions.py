"""
Медицинские состояния и патологические сердечные ритмы
"""

from enum import Enum
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass


class HeartRhythmState(Enum):
    """Состояния сердечного ритма"""
    NORMAL = "normal"
    TACHYCARDIA = "tachycardia"          # Тахикардия
    BRADYCARDIA = "bradycardia"          # Брадикардия 
    ATRIAL_FIBRILLATION = "afib"         # Мерцательная аритмия
    IRREGULAR = "irregular"              # Нерегулярный ритм
    PREMATURE_BEATS = "premature"        # Преждевременные сокращения


@dataclass
class ConditionParameters:
    """Параметры медицинского состояния"""
    name: str
    description: str
    rhythm_state: HeartRhythmState
    hr_range: Tuple[int, int]  # Диапазон ЧСС
    rr_variability: float      # Коэффициент вариабельности R-R
    is_anomaly: bool = True
    severity: str = "moderate"  # mild, moderate, severe
    medical_code: Optional[str] = None


class MedicalConditions:
    """
    Медицинские состояния для генерации реалистичных патологических данных
    """
    
    def __init__(self):
        self.conditions = self._initialize_conditions()
    
    def _initialize_conditions(self) -> Dict[str, ConditionParameters]:
        """Инициализирует медицинские состояния"""
        return {
            # Аритмии
            "atrial_fibrillation": ConditionParameters(
                name="Мерцательная аритмия",
                description="Нерегулярный и часто быстрый сердечный ритм",
                rhythm_state=HeartRhythmState.ATRIAL_FIBRILLATION,
                hr_range=(80, 150),
                rr_variability=0.3,  # Высокая вариабельность
                is_anomaly=True,
                severity="moderate",
                medical_code="I48.0"
            ),
            
            "atrial_flutter": ConditionParameters(
                name="Трепетание предсердий",
                description="Быстрый но регулярный ритм предсердий",
                rhythm_state=HeartRhythmState.TACHYCARDIA,
                hr_range=(120, 200),
                rr_variability=0.05,  # Низкая вариабельность
                is_anomaly=True,
                severity="moderate",
                medical_code="I48.3"
            ),
            
            "supraventricular_tachycardia": ConditionParameters(
                name="Наджелудочковая тахикардия",
                description="Быстрый ритм из верхних камер сердца",
                rhythm_state=HeartRhythmState.TACHYCARDIA,
                hr_range=(150, 220),
                rr_variability=0.02,  # Очень регулярно
                is_anomaly=True,
                severity="severe",
                medical_code="I47.1"
            ),
            
            # Блокады
            "sinus_bradycardia": ConditionParameters(
                name="Синусовая брадикардия",
                description="Медленный но регулярный сердечный ритм",
                rhythm_state=HeartRhythmState.BRADYCARDIA,
                hr_range=(35, 59),
                rr_variability=0.04,  # Нормальная вариабельность
                is_anomaly=True,
                severity="mild",
                medical_code="R00.1"
            ),
            
            "complete_heart_block": ConditionParameters(
                name="Полная атриовентрикулярная блокада",
                description="Полная блокада проведения между предсердиями и желудочками",
                rhythm_state=HeartRhythmState.BRADYCARDIA,
                hr_range=(25, 45),
                rr_variability=0.08,
                is_anomaly=True,
                severity="severe",
                medical_code="I44.2"
            ),
            
            # Экстрасистолы
            "premature_ventricular_contractions": ConditionParameters(
                name="Желудочковые экстрасистолы",
                description="Преждевременные сокращения желудочков",
                rhythm_state=HeartRhythmState.PREMATURE_BEATS,
                hr_range=(60, 100),
                rr_variability=0.2,  # Высокая из-за экстрасистол
                is_anomaly=True,
                severity="mild",
                medical_code="I49.3"
            ),
            
            "premature_atrial_contractions": ConditionParameters(
                name="Предсердные экстрасистолы",
                description="Преждевременные сокращения предсердий",
                rhythm_state=HeartRhythmState.PREMATURE_BEATS,
                hr_range=(65, 95),
                rr_variability=0.15,
                is_anomaly=True,
                severity="mild",
                medical_code="I49.1"
            ),
            
            # Нормальные вариации
            "sinus_arrhythmia": ConditionParameters(
                name="Синусовая аритмия",
                description="Нормальная вариация ритма при дыхании",
                rhythm_state=HeartRhythmState.NORMAL,
                hr_range=(60, 100),
                rr_variability=0.08,
                is_anomaly=False,
                severity="mild",
                medical_code=None
            ),
            
            "athletic_heart": ConditionParameters(
                name="Спортивное сердце",
                description="Адаптация сердца к тренировкам",
                rhythm_state=HeartRhythmState.BRADYCARDIA,
                hr_range=(45, 65),
                rr_variability=0.12,  # Высокая HRV у спортсменов
                is_anomaly=False,
                severity="mild",
                medical_code=None
            ),
            
            # Тестовые состояния для стресс-тестирования
            "exercise_stress": ConditionParameters(
                name="Физическая нагрузка",
                description="Сердечный ритм при физической активности",
                rhythm_state=HeartRhythmState.TACHYCARDIA,
                hr_range=(100, 180),
                rr_variability=0.03,
                is_anomaly=False,
                severity="mild",
                medical_code=None
            ),
            
            "emotional_stress": ConditionParameters(
                name="Эмоциональный стресс",
                description="Ритм при психологическом стрессе",
                rhythm_state=HeartRhythmState.TACHYCARDIA,
                hr_range=(85, 120),
                rr_variability=0.06,
                is_anomaly=False,
                severity="mild",
                medical_code=None
            ),
            
            "sleep_state": ConditionParameters(
                name="Сон",
                description="Сердечный ритм во время сна",
                rhythm_state=HeartRhythmState.BRADYCARDIA,
                hr_range=(50, 70),
                rr_variability=0.1,
                is_anomaly=False,
                severity="mild",
                medical_code=None
            )
        }
    
    def get_condition_parameters(self, condition_name: str) -> Optional[ConditionParameters]:
        """Возвращает параметры медицинского состояния"""
        return self.conditions.get(condition_name)
    
    def get_available_conditions(self) -> List[str]:
        """Возвращает список доступных медицинских состояний"""
        return list(self.conditions.keys())
    
    def get_conditions_by_severity(self, severity: str) -> List[str]:
        """Возвращает состояния по уровню тяжести"""
        return [
            name for name, params in self.conditions.items()
            if params.severity == severity
        ]
    
    def get_anomalous_conditions(self) -> List[str]:
        """Возвращает только аномальные (патологические) состояния"""
        return [
            name for name, params in self.conditions.items()
            if params.is_anomaly
        ]
    
    def get_normal_variations(self) -> List[str]:
        """Возвращает нормальные вариации ритма"""
        return [
            name for name, params in self.conditions.items()
            if not params.is_anomaly
        ]
    
    def get_conditions_by_rhythm_state(self, rhythm_state: HeartRhythmState) -> List[str]:
        """Возвращает состояния по типу ритма"""
        return [
            name for name, params in self.conditions.items()
            if params.rhythm_state == rhythm_state
        ]
    
    def get_condition_info(self, condition_name: str) -> Dict:
        """Возвращает подробную информацию о состоянии"""
        params = self.get_condition_parameters(condition_name)
        if not params:
            return {}
        
        return {
            "name": params.name,
            "description": params.description,
            "rhythm_type": params.rhythm_state.value,
            "heart_rate_range": f"{params.hr_range[0]}-{params.hr_range[1]} BPM",
            "variability": f"{params.rr_variability:.1%}",
            "is_pathological": params.is_anomaly,
            "severity": params.severity,
            "medical_code": params.medical_code or "N/A"
        }
    
    def suggest_test_scenarios(self) -> Dict[str, List[str]]:
        """Предлагает тестовые сценарии для разных целей"""
        return {
            "basic_testing": [
                "sinus_arrhythmia",
                "sinus_bradycardia", 
                "exercise_stress"
            ],
            "pathology_detection": [
                "atrial_fibrillation",
                "premature_ventricular_contractions",
                "supraventricular_tachycardia"
            ],
            "stress_testing": [
                "atrial_fibrillation",
                "complete_heart_block",
                "exercise_stress",
                "emotional_stress"
            ],
            "hrv_analysis": [
                "athletic_heart",
                "sinus_arrhythmia",
                "sleep_state"
            ],
            "algorithm_validation": [
                "atrial_flutter",
                "premature_atrial_contractions",
                "sinus_bradycardia"
            ]
        }
    
    @staticmethod
    def get_available_conditions() -> List[str]:
        """Статический метод для получения списка состояний (для argparse)"""
        temp_instance = MedicalConditions()
        return list(temp_instance.conditions.keys()) 