import { useState } from 'react';
import { toast } from 'sonner';
import { LuTrash2 } from 'react-icons/lu';
import { LoadingSpinner } from 'ui/loading-spinner';
import { Button } from 'ui/button';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
    AlertDialogOverlay,
} from 'ui/alert-dialog';
import { Card, CardContent, CardHeader, CardTitle } from 'ui/card';

import { useGetDeviceListQuery, useDeleteDeviceMutation, DeviceListItem } from '../api';

export function DevicesList() {
    const { data, isLoading } = useGetDeviceListQuery();
    const [deleteDevice, { isLoading: isDeleting }] = useDeleteDeviceMutation();
    const [deletingDeviceId, setDeletingDeviceId] = useState<string | null>(null);

    const handleDeleteDevice = async (device: DeviceListItem) => {
        setDeletingDeviceId(device.id);
        try {
            await deleteDevice({ id: device.id }).unwrap();
            toast.success('Устройство успешно удалено');
        } catch (error) {
            toast.error('Произошла ошибка при удалении устройства');
        } finally {
            setDeletingDeviceId(null);
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center p-8">
                <LoadingSpinner />
                <span className="ml-2">Загрузка устройств...</span>
            </div>
        );
    }

    if (!data?.devices || data.devices.length === 0) {
        return (
            <div className="text-center p-8">
                <h3 className="text-lg font-medium mb-2">Нет подключенных устройств</h3>
                <p className="text-muted-foreground">
                    Подключите устройство через мобильное приложение для отображения здесь
                </p>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex flex-col">
                <h2 className="text-2xl font-bold mb-2">Устройства</h2>
                <p className="text-muted-foreground mb-6">
                    Управление подключенными устройствами. Добавление новых устройств доступно только через мобильное приложение.
                </p>
            </div>

            <div className="grid gap-4">
                {data.devices.map((device) => (
                    <Card key={device.id}>
                        <CardHeader className="pb-3">
                            <div className="flex items-center justify-between">
                                <CardTitle className="text-lg">{device.device_name}</CardTitle>
                                <AlertDialog>
                                    <AlertDialogTrigger asChild>
                                        <Button 
                                            variant="destructive" 
                                            size="icon"
                                            disabled={deletingDeviceId === device.id}
                                        >
                                            <LuTrash2 className="h-4 w-4" />
                                        </Button>
                                    </AlertDialogTrigger>
                                    <AlertDialogOverlay className="bg-black/80" />
                                    <AlertDialogContent className="bg-white border border-gray-200 shadow-xl">
                                        <AlertDialogHeader>
                                            <AlertDialogTitle className="text-gray-900">Удалить устройство?</AlertDialogTitle>
                                            <AlertDialogDescription className="text-gray-600">
                                                Вы уверены, что хотите удалить устройство "{device.device_name}"? 
                                                Это действие нельзя отменить.
                                            </AlertDialogDescription>
                                        </AlertDialogHeader>
                                        <AlertDialogFooter>
                                            <AlertDialogCancel className="bg-gray-100 text-gray-900 hover:bg-gray-200">Отмена</AlertDialogCancel>
                                            <AlertDialogAction
                                                onClick={() => handleDeleteDevice(device)}
                                                className="bg-red-600 text-white hover:bg-red-700"
                                            >
                                                Удалить
                                            </AlertDialogAction>
                                        </AlertDialogFooter>
                                    </AlertDialogContent>
                                </AlertDialog>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="text-sm text-muted-foreground">
                                <p>ID: {device.id}</p>
                                <p>Добавлено: {formatDate(device.created_at)}</p>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
} 