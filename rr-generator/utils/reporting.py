"""
Генератор отчетов
"""

from typing import Dict, Any
from datetime import datetime
import json


class ReportGenerator:
    """Генерирует отчеты о процессе генерации данных"""
    
    def generate_comprehensive_report(self, results: Dict[str, Any]):
        """Генерирует подробный отчет"""
        print("\n" + "="*60)
        print("📊 ОТЧЕТ О ГЕНЕРАЦИИ ДАННЫХ")
        print("="*60)
        
        scenario = results.get('scenario', {})
        datasets = results.get('datasets', {})
        validation = results.get('validation', {})
        upload = results.get('upload')
        
        # Сценарий
        print(f"\n📋 Сценарий: {scenario.get('name', 'Unknown')}")
        print(f"👤 Пользователь: {scenario.get('user_id', 'N/A')}")
        print(f"📅 Период: {scenario.get('start_time')} - {scenario.get('end_time')}")
        
        # Сгенерированные данные
        print(f"\n📦 Сгенерированные наборы данных:")
        total_points = 0
        for dataset_name, data_points in datasets.items():
            count = len(data_points) if data_points else 0
            total_points += count
            print(f"  • {dataset_name}: {count:,} точек")
        print(f"  Всего: {total_points:,} точек")
        
        # Валидация
        if validation:
            print(f"\n✅ Валидация:")
            summary = validation.get('summary', {})
            print(f"  • Валидных наборов: {summary.get('valid_datasets', 0)}/{summary.get('total_datasets', 0)}")
            print(f"  • Средний балл: {summary.get('average_score', 0):.2f}")
            print(f"  • Общих проблем: {summary.get('total_issues', 0)}")
        
        # Загрузка
        if upload:
            print(f"\n📤 Загрузка в API:")
            total_uploaded = sum(result.get('uploaded_count', 0) for result in upload.values() if isinstance(result, dict))
            print(f"  • Загружено точек: {total_uploaded:,}")
            
            successful_uploads = sum(1 for result in upload.values() if isinstance(result, dict) and result.get('success'))
            print(f"  • Успешных загрузок: {successful_uploads}/{len(upload)}")
        
        print(f"\n⏰ Отчет сгенерирован: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("="*60) 