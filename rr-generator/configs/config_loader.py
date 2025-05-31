"""
Загрузчик конфигурации для генератора
"""

import os
from typing import Dict, Any


class ConfigLoader:
    """Загружает и управляет конфигурацией"""
    
    def __init__(self, config_file: str = None):
        self.config = self._load_default_config()
        
        if config_file and os.path.exists(config_file):
            self._load_from_file(config_file)
    
    def _load_default_config(self) -> Dict[str, Any]:
        """Загружает конфигурацию по умолчанию"""
        return {
            "generation": {
                "default_sampling_rate": 1000,  # Hz
                "default_duration_hours": 1,
                "noise_level": 0.02,
                "hrv_enabled": True
            },
            "validation": {
                "min_rr_ms": 200,
                "max_rr_ms": 3000,
                "max_outlier_ratio": 0.05,
                "min_quality_score": 0.5
            },
            "api": {
                "default_batch_size": 100,
                "max_retry_attempts": 3,
                "timeout_seconds": 30
            },
            "output": {
                "default_format": "json",
                "compression": False,
                "include_metadata": True
            }
        }
    
    def _load_from_file(self, config_file: str):
        """Загружает конфигурацию из файла"""
        # Реализация для YAML/JSON файлов
        pass
    
    def get(self, key: str, default=None):
        """Получает значение конфигурации"""
        keys = key.split('.')
        value = self.config
        
        for k in keys:
            if isinstance(value, dict) and k in value:
                value = value[k]
            else:
                return default
        
        return value 