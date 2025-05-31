"""
Утилиты для работы со временем
"""

from datetime import datetime, timedelta, timezone
from typing import Union
import re


class TimeUtils:
    """Утилиты для парсинга и работы со временем"""
    
    # Московская временная зона (UTC+3)
    MSK_TIMEZONE = timezone(timedelta(hours=3))
    
    @staticmethod
    def parse_datetime(time_str: str, default_timezone: timezone = None) -> datetime:
        """
        Парсит строку времени в различных форматах с поддержкой временных зон
        
        Поддерживаемые форматы:
        - "2025-06-01 14:30"
        - "2025-06-01 14:30:00"
        - "2025-06-01T14:30:00"
        - "2025-06-01 14:30:00+03:00" (с временной зоной)
        - "2025-06-01T14:30:00Z" (UTC)
        - "now" - текущее время
        
        Args:
            time_str: Строка времени для парсинга
            default_timezone: Временная зона по умолчанию (если не указана в строке).
                            По умолчанию использует MSK (UTC+3)
        """
        if time_str.lower() == "now":
            return datetime.now(default_timezone or TimeUtils.MSK_TIMEZONE)
        
        # Форматы для парсинга с временными зонами
        tz_formats = [
            "%Y-%m-%d %H:%M:%S%z",
            "%Y-%m-%dT%H:%M:%S%z",
            "%Y-%m-%d %H:%M%z",
            "%Y-%m-%dT%H:%M%z",
            "%Y-%m-%dT%H:%M:%SZ",  # UTC формат
            "%Y-%m-%d %H:%M:%SZ",  # UTC формат
        ]
        
        # Сначала пробуем форматы с временными зонами
        for fmt in tz_formats:
            try:
                if fmt.endswith('Z'):
                    # Заменяем Z на +00:00 для корректного парсинга UTC
                    time_str_modified = time_str.replace('Z', '+00:00')
                    fmt_modified = fmt.replace('Z', '%z')
                    return datetime.strptime(time_str_modified, fmt_modified)
                else:
                    return datetime.strptime(time_str, fmt)
            except ValueError:
                continue
        
        # Форматы без временных зон (naive datetime)
        naive_formats = [
            "%Y-%m-%d %H:%M",
            "%Y-%m-%d %H:%M:%S",
            "%Y-%m-%dT%H:%M:%S",
            "%Y-%m-%d",
            "%d.%m.%Y %H:%M",
            "%d.%m.%Y",
        ]
        
        # Если не нашли формат с временной зоной, пробуем naive форматы
        for fmt in naive_formats:
            try:
                naive_dt = datetime.strptime(time_str, fmt)
                # Применяем временную зону по умолчанию (МСК)
                timezone_to_use = default_timezone or TimeUtils.MSK_TIMEZONE
                return naive_dt.replace(tzinfo=timezone_to_use)
            except ValueError:
                continue
        
        raise ValueError(f"Не удалось распарсить время: {time_str}")
    
    @staticmethod
    def parse_datetime_msk(time_str: str) -> datetime:
        """
        Парсит строку времени и принудительно применяет МСК временную зону
        """
        return TimeUtils.parse_datetime(time_str, TimeUtils.MSK_TIMEZONE)
    
    @staticmethod
    def parse_datetime_utc(time_str: str) -> datetime:
        """
        Парсит строку времени и применяет UTC временную зону
        """
        return TimeUtils.parse_datetime(time_str, timezone.utc)
    
    @staticmethod
    def parse_duration(duration_str: str) -> timedelta:
        """
        Парсит строку длительности
        
        Поддерживаемые форматы:
        - "1h" - 1 час
        - "30m" - 30 минут
        - "45s" - 45 секунд
        - "2d" - 2 дня
        - "1h30m" - 1 час 30 минут
        - "2d4h" - 2 дня 4 часа
        """
        duration_str = duration_str.lower().strip()
        
        # Регулярные выражения для разных единиц
        patterns = {
            r'(\d+)d': 'days',
            r'(\d+)h': 'hours', 
            r'(\d+)m': 'minutes',
            r'(\d+)s': 'seconds'
        }
        
        kwargs = {}
        
        for pattern, unit in patterns.items():
            matches = re.findall(pattern, duration_str)
            if matches:
                # Берем последнее значение если есть несколько
                kwargs[unit] = int(matches[-1])
        
        if not kwargs:
            # Попробуем простые числа (предполагаем минуты)
            try:
                minutes = int(duration_str)
                return timedelta(minutes=minutes)
            except ValueError:
                raise ValueError(f"Не удалось распарсить длительность: {duration_str}")
        
        return timedelta(**kwargs)
    
    @staticmethod
    def format_duration(duration: timedelta) -> str:
        """Форматирует timedelta в читаемую строку"""
        total_seconds = int(duration.total_seconds())
        
        days = total_seconds // 86400
        hours = (total_seconds % 86400) // 3600
        minutes = (total_seconds % 3600) // 60
        seconds = total_seconds % 60
        
        parts = []
        if days > 0:
            parts.append(f"{days}d")
        if hours > 0:
            parts.append(f"{hours}h")
        if minutes > 0:
            parts.append(f"{minutes}m")
        if seconds > 0 or not parts:
            parts.append(f"{seconds}s")
        
        return " ".join(parts)
    
    @staticmethod
    def get_time_range_info(start_time: datetime, end_time: datetime) -> dict:
        """Возвращает информацию о временном диапазоне"""
        duration = end_time - start_time
        
        return {
            "start": start_time.strftime("%Y-%m-%d %H:%M:%S"),
            "end": end_time.strftime("%Y-%m-%d %H:%M:%S"),
            "duration": TimeUtils.format_duration(duration),
            "duration_seconds": duration.total_seconds(),
            "duration_hours": duration.total_seconds() / 3600,
            "duration_days": duration.days,
        }
    
    @staticmethod
    def split_time_range(start_time: datetime, end_time: datetime, 
                        chunk_duration: timedelta) -> list:
        """Разбивает временной диапазон на чанки заданной длительности"""
        chunks = []
        current_time = start_time
        
        while current_time < end_time:
            chunk_end = min(current_time + chunk_duration, end_time)
            chunks.append((current_time, chunk_end))
            current_time = chunk_end
        
        return chunks
    
    @staticmethod
    def estimate_data_points(start_time: datetime, end_time: datetime, 
                           avg_rr_ms: float = 800) -> int:
        """Оценивает количество точек данных для временного диапазона"""
        duration = end_time - start_time
        duration_ms = duration.total_seconds() * 1000
        return int(duration_ms / avg_rr_ms)
    
    @staticmethod
    def get_circadian_phase(timestamp: datetime) -> str:
        """Возвращает фазу циркадного ритма"""
        hour = timestamp.hour
        
        if 6 <= hour < 12:
            return "morning"
        elif 12 <= hour < 18:
            return "afternoon"
        elif 18 <= hour < 22:
            return "evening"
        else:
            return "night"
    
    @staticmethod
    def is_business_hours(timestamp: datetime) -> bool:
        """Проверяет, попадает ли время в рабочие часы"""
        return 9 <= timestamp.hour < 17 and timestamp.weekday() < 5 