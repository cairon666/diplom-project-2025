#!/usr/bin/env python3
"""
🫀 Единый генератор R-R интервалов для тестирования системы

Одна команда для создания любых наборов данных:
- Хорошие и плохие данные
- Реалистичные медицинские сценарии  
- Автоматическая загрузка в API
- Валидация и отчеты

Примеры использования:
    # Создать реалистичный день для пользователя
    python generate_all_data.py --user-id UUID --realistic-day
    
    # Стресс-тест с плохими данными
    python generate_all_data.py --stress-test --bad-ratio 50
    
    # Специфическое медицинское состояние
    python generate_all_data.py --condition atrial_fibrillation --duration 2h
"""

import argparse
import sys
import os
from datetime import datetime, timedelta
from typing import Dict, List, Optional
import uuid
import json

# Добавляем путь к модулям
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from core.data_generator import UnifiedDataGenerator, RRDataPoint
from core.medical_conditions import MedicalConditions
from core.validators import DataValidator
from core.api_client import HealthAPIClient
from utils.time_utils import TimeUtils
from utils.reporting import ReportGenerator
from configs.config_loader import ConfigLoader


def main():
    parser = create_argument_parser()
    args = parser.parse_args()
    
    print("🫀 R-R Intervals Data Generator v2.0")
    print("=" * 50)
    
    # Инициализация компонентов
    config = ConfigLoader()
    generator = UnifiedDataGenerator(config)
    validator = DataValidator(config)
    api_client = HealthAPIClient(args.api_url, args.auth_token)
    reporter = ReportGenerator()
    
    try:
        # Определяем сценарий генерации
        scenario = determine_scenario(args)
        print(f"📋 Сценарий: {scenario['name']}")
        print(f"📅 Период: {scenario['start_time']} - {scenario['end_time']}")
        print(f"👤 Пользователь: {scenario['user_id']}")
        print()
        
        # Генерируем данные
        datasets = generate_data(generator, scenario, args)
        
        # Показываем примеры данных
        preview_data_samples(datasets)
        
        # Валидируем данные
        validation_results = validate_data(validator, datasets)
        
        # Сохраняем локально
        save_results = save_data_locally(datasets, scenario)
        
        # Загружаем в API если указан токен
        upload_results = None
        if args.upload and args.auth_token:
            upload_results = upload_to_api(api_client, datasets, args.batch_size)
        
        # Генерируем отчет
        generate_final_report(reporter, {
            'scenario': scenario,
            'datasets': datasets,
            'validation': validation_results,
            'save': save_results,
            'upload': upload_results
        })
        
        print("✅ Генерация данных завершена успешно!")
        
    except KeyboardInterrupt:
        print("\n⚠️ Операция прервана пользователем")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Ошибка: {e}")
        sys.exit(1)


def create_argument_parser():
    """Создает парсер командной строки с всеми опциями"""
    parser = argparse.ArgumentParser(
        description='Генератор R-R интервалов для тестирования системы',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Примеры использования:

  Реалистичный день:
    python generate_all_data.py --realistic-day --user-id "uuid-here" --upload

  Стресс-тестирование:
    python generate_all_data.py --stress-test --bad-ratio 30 --users 10

  Конкретное состояние:
    python generate_all_data.py --condition atrial_fibrillation --duration 4h

  Конкретный период:
    python generate_all_data.py --start "2025-06-01 00:00" --end "2025-06-01 23:59"

  Только хорошие данные (без невалидных):
    python generate_all_data.py --realistic-day --no-invalid-data --upload
        """
    )
    
    # Основные параметры
    parser.add_argument('--user-id', help='UUID пользователя (автогенерация если не указан)')
    parser.add_argument('--device-id', help='UUID устройства (автогенерация если не указан)')
    
    # Временные параметры
    time_group = parser.add_mutually_exclusive_group()
    time_group.add_argument('--realistic-day', action='store_true', 
                           help='Генерировать реалистичный день (24 часа)')
    time_group.add_argument('--start', help='Время начала (YYYY-MM-DD HH:MM)')
    time_group.add_argument('--duration', help='Длительность (1h, 30m, 2d)')
    
    parser.add_argument('--end', help='Время окончания (YYYY-MM-DD HH:MM)')
    
    # Сценарии генерации
    scenario_group = parser.add_mutually_exclusive_group()
    scenario_group.add_argument('--stress-test', action='store_true',
                               help='Стресс-тестирование с большим объемом данных')
    scenario_group.add_argument('--condition', choices=MedicalConditions.get_available_conditions(),
                               help='Генерировать конкретное медицинское состояние')
    scenario_group.add_argument('--custom-scenario', help='Путь к YAML файлу со сценарием')
    
    # Параметры качества данных
    parser.add_argument('--good-bad-ratio', default='70:30', 
                       help='Соотношение хороших к плохим данным (default: 70:30)')
    parser.add_argument('--bad-ratio', type=int, default=30,
                       help='Процент плохих данных (0-100, default: 30)')
    parser.add_argument('--no-invalid-data', action='store_true',
                       help='Отключить генерацию тестовых невалидных данных (только хорошие данные)')
    
    # Множественные пользователи
    parser.add_argument('--users', type=int, default=1,
                       help='Количество пользователей для генерации')
    parser.add_argument('--concurrent', action='store_true',
                       help='Параллельная генерация для множественных пользователей')
    
    # API параметры
    parser.add_argument('--upload', action='store_true',
                       help='Загрузить данные в API')
    parser.add_argument('--api-url', default='http://localhost:8080',
                       help='URL API сервера')
    parser.add_argument('--auth-token', 
                       help='JWT токен (или установите в переменной AUTH_TOKEN)')
    parser.add_argument('--batch-size', type=int, default=100,
                       help='Размер батча для загрузки в API')
    
    # Вывод и отчеты
    parser.add_argument('--output-dir', default='generated_data',
                       help='Папка для сохранения результатов')
    parser.add_argument('--report', action='store_true', default=True,
                       help='Генерировать подробный отчет')
    parser.add_argument('--plot', action='store_true',
                       help='Создать графики и визуализации')
    parser.add_argument('--quiet', action='store_true',
                       help='Минимальный вывод')
    
    # Отладка
    parser.add_argument('--dry-run', action='store_true',
                       help='Показать что будет сгенерировано без фактической генерации')
    parser.add_argument('--validate-only', action='store_true',
                       help='Только валидация существующих данных')
    
    return parser


def determine_scenario(args) -> Dict:
    """Определяет сценарий генерации на основе аргументов"""
    scenario = {
        'user_id': args.user_id or str(uuid.uuid4()),
        'device_id': args.device_id or str(uuid.uuid4()),
    }
    
    # ПРИОРИТЕТ 1: Медицинское состояние (высший приоритет)
    if args.condition:
        # Определяем временной диапазон для медицинского состояния
        if args.start:
            start_time = TimeUtils.parse_datetime_msk(args.start)
            if args.end:
                end_time = TimeUtils.parse_datetime_msk(args.end)
            elif args.duration:
                end_time = start_time + TimeUtils.parse_duration(args.duration)
            else:
                end_time = start_time + timedelta(hours=1)
        else:
            # Для медицинского состояния используем короткий период
            start_time = datetime.now(TimeUtils.MSK_TIMEZONE)
            duration = TimeUtils.parse_duration(args.duration) if args.duration else timedelta(hours=2)
            end_time = start_time + duration
        
        scenario.update({
            'name': f'Медицинское состояние: {args.condition}',
            'start_time': start_time,
            'end_time': end_time,
            'type': 'medical_condition',
            'condition': args.condition
        })
    
    # ПРИОРИТЕТ 2: Реалистичный день
    elif args.realistic_day:
        start_time = datetime.now(TimeUtils.MSK_TIMEZONE).replace(hour=0, minute=0, second=0, microsecond=0)
        end_time = start_time + timedelta(days=1)
        scenario.update({
            'name': 'Реалистичный день',
            'start_time': start_time,
            'end_time': end_time,
            'type': 'realistic_day'
        })
    
    # ПРИОРИТЕТ 3: Стресс-тест
    elif args.stress_test:
        # Стресс-тест: последние 24 часа
        end_time = datetime.now(TimeUtils.MSK_TIMEZONE)
        start_time = end_time - timedelta(days=1)
        
        scenario.update({
            'name': 'Стресс-тестирование',
            'start_time': start_time,
            'end_time': end_time,
            'type': 'stress_test'
        })
    
    # ПРИОРИТЕТ 4: Настраиваемый период времени
    elif args.start:
        start_time = TimeUtils.parse_datetime_msk(args.start)
        if args.end:
            end_time = TimeUtils.parse_datetime_msk(args.end)
        elif args.duration:
            end_time = start_time + TimeUtils.parse_duration(args.duration)
        else:
            end_time = start_time + timedelta(hours=1)  # default 1 hour
            
        scenario.update({
            'name': 'Настраиваемый период',
            'start_time': start_time,
            'end_time': end_time,
            'type': 'custom_period'
        })
    
    else:
        # Default: последний час
        end_time = datetime.now(TimeUtils.MSK_TIMEZONE)
        start_time = end_time - timedelta(hours=1)
        
        scenario.update({
            'name': 'Быстрый тест',
            'start_time': start_time,
            'end_time': end_time,
            'type': 'quick_test'
        })
    
    # Добавляем параметры качества данных
    scenario['good_bad_ratio'] = parse_ratio(args.good_bad_ratio)
    scenario['bad_ratio'] = args.bad_ratio
    scenario['users_count'] = args.users
    
    # Если отключены невалидные данные, переопределяем параметры
    if args.no_invalid_data:
        scenario['good_bad_ratio'] = (100, 0)
        scenario['bad_ratio'] = 0
    
    return scenario


def parse_ratio(ratio_str: str) -> tuple:
    """Парсит строку соотношения типа '70:30' в кортеж (70, 30)"""
    try:
        parts = ratio_str.split(':')
        return (int(parts[0]), int(parts[1]))
    except:
        return (70, 30)  # default


def generate_data(generator: UnifiedDataGenerator, scenario: Dict, args) -> Dict:
    """Генерирует данные согласно сценарию"""
    print("🔄 Генерация данных...")
    
    if args.dry_run:
        print("🔍 Режим предварительного просмотра:")
        generator.preview_generation(scenario)
        return {}
    
    return generator.generate_scenario_data(scenario)


def validate_data(validator: DataValidator, datasets: Dict) -> Dict:
    """Валидирует сгенерированные данные"""
    if not datasets:
        return {}
        
    print("✅ Валидация данных...")
    return validator.validate_all_datasets(datasets)


def save_data_locally(datasets: Dict, scenario: Dict) -> Dict:
    """Сохраняет данные локально"""
    if not datasets:
        return {}
        
    print("💾 Сохранение данных...")
    # Implementation will be in the actual generator
    return {'status': 'saved', 'files': []}


def upload_to_api(api_client: HealthAPIClient, datasets: Dict, batch_size: int) -> Dict:
    """Загружает данные в API"""
    if not datasets:
        return {}
        
    print("📤 Загрузка в API...")
    return api_client.upload_all_datasets(datasets, batch_size)


def generate_final_report(reporter: ReportGenerator, results: Dict):
    """Генерирует финальный отчет"""
    print("📊 Генерация отчета...")
    reporter.generate_comprehensive_report(results)


def preview_data_samples(datasets: Dict[str, List[RRDataPoint]], num_samples: int = 3):
    """Показывает примеры сгенерированных данных"""
    print(f"\n📋 Примеры сгенерированных данных:")
    
    for dataset_name, data_points in datasets.items():
        if not data_points:
            continue
            
        print(f"\n• {dataset_name.upper()}: (всего {len(data_points)} точек)")
        
        # Показываем несколько примеров
        samples = data_points[:num_samples] if len(data_points) >= num_samples else data_points
        
        for i, point in enumerate(samples, 1):
            bpm = round(60000 / point.rr_interval_ms) if point.rr_interval_ms > 0 else 0
            print(f"    {i}. R-R: {point.rr_interval_ms:.0f} мс → {bpm} BPM | Качество: {point.quality_score:.2f} | Источник: {point.source_condition}")
        
        # Показываем статистику
        rr_values = [p.rr_interval_ms for p in data_points]
        bpm_values = [60000/rr for rr in rr_values if rr > 0]
        
        print(f"    📊 Диапазон R-R: {min(rr_values):.0f}-{max(rr_values):.0f} мс")
        print(f"    💓 Диапазон BPM: {min(bpm_values):.0f}-{max(bpm_values):.0f}")
        print(f"    📈 Средний BPM: {sum(bpm_values)/len(bpm_values):.0f}")


if __name__ == "__main__":
    main() 