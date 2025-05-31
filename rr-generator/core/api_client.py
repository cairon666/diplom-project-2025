"""
API клиент для загрузки данных в Health API
"""

import requests
import json
from typing import List, Dict, Optional
from datetime import datetime, timedelta, timezone
import time
import os

from .data_generator import RRDataPoint


class HealthAPIClient:
    """Клиент для взаимодействия с Health API"""
    
    def __init__(self, api_url: str, auth_token: Optional[str] = None):
        self.api_url = api_url.rstrip('/')
        self.auth_token = auth_token or os.getenv('AUTH_TOKEN')
        self.session = requests.Session()
        
        if self.auth_token:
            self.session.headers.update({
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            })
    
    def upload_all_datasets(self, datasets: Dict[str, List[RRDataPoint]], 
                          batch_size: int = 100) -> Dict:
        """Загружает все наборы данных"""
        results = {}
        
        for dataset_name, data_points in datasets.items():
            print(f"📤 Загрузка {dataset_name}: {len(data_points)} точек...")
            
            try:
                upload_result = self.upload_rr_data(data_points, batch_size)
                results[dataset_name] = upload_result
                
                if upload_result['success']:
                    uploaded = upload_result['uploaded_count']
                    valid = upload_result.get('valid_count', uploaded)
                    print(f"✅ {dataset_name}: {uploaded} обработано, {valid} валидных")
                else:
                    print(f"❌ {dataset_name}: загрузка не удалась")
                
            except Exception as e:
                results[dataset_name] = {
                    'success': False,
                    'error': str(e),
                    'uploaded_count': 0,
                    'valid_count': 0
                }
                print(f"❌ {dataset_name}: ошибка - {e}")
        
        return results
    
    def upload_rr_data(self, data_points: List[RRDataPoint], 
                      batch_size: int = 100) -> Dict:
        """Загружает R-R данные пакетами"""
        if not data_points:
            return {'success': True, 'uploaded_count': 0, 'batches': 0}
        
        total_points = len(data_points)
        uploaded_count = 0
        valid_count = 0
        batch_count = 0
        errors = []
        
        # Разбиваем на пакеты
        for i in range(0, total_points, batch_size):
            batch = data_points[i:i + batch_size]
            batch_count += 1
            
            try:
                # Конвертируем в формат API
                api_data = self._convert_to_api_format(batch)
                
                # Отправляем пакет
                response = self._send_batch(api_data)
                
                if response['success']:
                    processed = response.get('processed_count', len(batch))
                    valid = response.get('valid_count', len(batch))
                    uploaded_count += processed
                    valid_count += valid
                    print(f"  Пакет {batch_count}: ✅ {processed} обработано, {valid} валидных")
                else:
                    errors.append(f"Пакет {batch_count}: {response.get('error', 'Unknown error')}")
                    print(f"  Пакет {batch_count}: ❌ {response.get('error')}")
                
                # Небольшая пауза между пакетами
                time.sleep(0.1)
                
            except Exception as e:
                error_msg = f"Пакет {batch_count}: {str(e)}"
                errors.append(error_msg)
                print(f"  Пакет {batch_count}: ❌ {e}")
        
        return {
            'success': len(errors) == 0,
            'uploaded_count': uploaded_count,
            'valid_count': valid_count,
            'total_count': total_points,
            'batches': batch_count,
            'errors': errors,
            'upload_rate': uploaded_count / total_points if total_points > 0 else 0
        }
    
    def _convert_to_api_format(self, data_points: List[RRDataPoint]) -> Dict:
        """Конвертирует данные в формат API"""
        if not data_points:
            return {"device_id": "", "intervals": []}
        
        # Берем device_id из первой точки (все точки должны иметь одинаковый device_id)
        device_id = data_points[0].device_id
        
        return {
            "device_id": device_id,
            "intervals": [
                {
                    "user_id": dp.user_id,
                    "device_id": dp.device_id,
                    "timestamp": self._format_timestamp_for_api(dp.timestamp),
                    "rr_interval_ms": int(round(dp.rr_interval_ms)),  # Конвертируем в int
                    "quality_score": dp.quality_score,
                    "is_anomaly": dp.is_anomaly,
                    "source": dp.source_condition
                }
                for dp in data_points
            ]
        }
    
    def _format_timestamp_for_api(self, timestamp: datetime) -> str:
        """Форматирует timestamp для API"""
        # Если timestamp уже в UTC, используем его как есть
        if timestamp.tzinfo is None:
            # Naive datetime - считаем что это уже UTC
            return timestamp.isoformat() + "Z"
        elif timestamp.utcoffset() == timedelta(0):
            # Уже в UTC
            return timestamp.replace(tzinfo=None).isoformat() + "Z"
        else:
            # Конвертируем в UTC
            utc_timestamp = timestamp.astimezone(timezone.utc)
            return utc_timestamp.replace(tzinfo=None).isoformat() + "Z"
    
    def _send_batch(self, data: Dict) -> Dict:
        """Отправляет один пакет данных"""
        try:
            response = self.session.post(
                f"{self.api_url}/v1/rr-intervals/batch",
                json=data,
                timeout=30
            )
            
            if response.status_code == 201:
                # Успешная загрузка - возвращаем только краткую статистику
                response_data = response.json()
                processed = response_data.get('processed_count', 0)
                valid = response_data.get('valid_count', 0)
                return {
                    'success': True, 
                    'processed_count': processed,
                    'valid_count': valid
                }
            elif response.status_code == 200:
                response_data = response.json()
                processed = response_data.get('processed_count', 0)
                valid = response_data.get('valid_count', 0)
                return {
                    'success': True, 
                    'processed_count': processed,
                    'valid_count': valid
                }
            elif response.status_code == 401:
                return {'success': False, 'error': 'Не авторизован - проверьте токен'}
            elif response.status_code == 413:
                return {'success': False, 'error': 'Пакет слишком большой'}
            else:
                return {
                    'success': False, 
                    'error': f'HTTP {response.status_code}: {response.text[:200]}...'  # Ограничиваем длину ошибки
                }
                
        except requests.exceptions.Timeout:
            return {'success': False, 'error': 'Timeout при отправке'}
        except requests.exceptions.ConnectionError:
            return {'success': False, 'error': 'Ошибка соединения с сервером'}
        except Exception as e:
            return {'success': False, 'error': str(e)}
    
    def test_connection(self) -> Dict:
        """Тестирует соединение с API"""
        try:
            response = self.session.get(f"{self.api_url}/health", timeout=10)
            
            if response.status_code == 200:
                return {
                    'success': True, 
                    'message': 'Соединение успешно',
                    'server_info': response.json()
                }
            else:
                return {
                    'success': False,
                    'error': f'Сервер вернул статус {response.status_code}'
                }
                
        except Exception as e:
            return {
                'success': False,
                'error': f'Ошибка соединения: {str(e)}'
            }
    
    def get_user_data_summary(self, user_id: str) -> Dict:
        """Получает сводку данных пользователя"""
        try:
            response = self.session.get(
                f"{self.api_url}/v1/users/{user_id}/rr-intervals/summary",
                timeout=10
            )
            
            if response.status_code == 200:
                return {'success': True, 'data': response.json()}
            else:
                return {'success': False, 'error': f'HTTP {response.status_code}'}
                
        except Exception as e:
            return {'success': False, 'error': str(e)}
    
    def delete_test_data(self, user_ids: List[str]) -> Dict:
        """Удаляет тестовые данные"""
        results = {}
        
        for user_id in user_ids:
            try:
                response = self.session.delete(
                    f"{self.api_url}/v1/users/{user_id}/rr-intervals",
                    timeout=30
                )
                
                if response.status_code in [200, 204]:
                    results[user_id] = {'success': True, 'message': 'Данные удалены'}
                else:
                    results[user_id] = {
                        'success': False, 
                        'error': f'HTTP {response.status_code}'
                    }
                    
            except Exception as e:
                results[user_id] = {'success': False, 'error': str(e)}
        
        return results 