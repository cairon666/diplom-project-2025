import { useState } from 'react';
import { toast } from 'sonner';
import { LuTrash2, LuPlus } from 'react-icons/lu';
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

import { useGetExternalAppsListQuery, useDeleteExternalAppMutation, ExternalAppListItem } from '../api';
import { CreateExternalAppModal } from './CreateExternalAppModal';

export function ExternalAppsList() {
    const { data, isLoading, error } = useGetExternalAppsListQuery();
    const [deleteExternalApp, { isLoading: isDeleting }] = useDeleteExternalAppMutation();
    const [deletingAppId, setDeletingAppId] = useState<string | null>(null);

    const handleDeleteApp = async (app: ExternalAppListItem) => {
        setDeletingAppId(app.id);
        try {
            await deleteExternalApp({ id: app.id }).unwrap();
            toast.success('Внешнее приложение успешно удалено');
        } catch (error) {
            toast.error('Произошла ошибка при удалении приложения');
        } finally {
            setDeletingAppId(null);
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
            <div className="flex-1 p-6">
                <div className="max-w-7xl mx-auto">
                    <div className="flex items-center justify-center p-8">
                        <LoadingSpinner />
                        <span className="ml-2">Загрузка приложений...</span>
                    </div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex-1 p-6">
                <div className="max-w-7xl mx-auto">
                    <div className="p-4 bg-red-100 border border-red-300 rounded">
                        <h3 className="font-bold text-red-800">Ошибка загрузки внешних приложений</h3>
                        <p className="text-red-700">
                            {error && typeof error === 'object' && 'data' in error 
                                ? JSON.stringify(error.data) 
                                : 'Неизвестная ошибка'}
                        </p>
                        <p className="text-sm text-red-600 mt-2">
                            Проверьте консоль разработчика для дополнительной информации.
                        </p>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="flex-1 p-6">
            <div className="max-w-7xl mx-auto">
                <div className="space-y-4">
                    <div className="flex flex-col">
                        <div className="flex items-center justify-between mb-2">
                            <h2 className="text-2xl font-bold">Внешние приложения</h2>
                            <CreateExternalAppModal>
                                <Button>
                                    <LuPlus className="h-4 w-4 mr-2" />
                                    Создать приложение
                                </Button>
                            </CreateExternalAppModal>
                        </div>
                        <p className="text-muted-foreground mb-6">
                            Управление подключенными внешними приложениями и их API ключами.
                        </p>
                    </div>

                    {!data?.external_apps || data.external_apps.length === 0 ? (
                        <div className="text-center p-8">
                            <h3 className="text-lg font-medium mb-2">Нет внешних приложений</h3>
                            <p className="text-muted-foreground mb-4">
                                Создайте первое внешнее приложение для интеграции с API.
                            </p>
                            <CreateExternalAppModal>
                                <Button>
                                    <LuPlus className="h-4 w-4 mr-2" />
                                    Создать приложение
                                </Button>
                            </CreateExternalAppModal>
                        </div>
                    ) : (
                        <div className="grid gap-4">
                            {data.external_apps.map((app) => (
                                <Card key={app.id}>
                                    <CardHeader className="pb-3">
                                        <div className="flex items-center justify-between">
                                            <CardTitle className="text-lg">{app.name}</CardTitle>
                                            <AlertDialog>
                                                <AlertDialogTrigger asChild>
                                                    <Button 
                                                        variant="destructive" 
                                                        size="icon"
                                                        disabled={deletingAppId === app.id}
                                                    >
                                                        <LuTrash2 className="h-4 w-4" />
                                                    </Button>
                                                </AlertDialogTrigger>
                                                <AlertDialogOverlay className="bg-black/80" />
                                                <AlertDialogContent className="bg-white border border-gray-200 shadow-xl">
                                                    <AlertDialogHeader>
                                                        <AlertDialogTitle className="text-gray-900">Удалить приложение?</AlertDialogTitle>
                                                        <AlertDialogDescription className="text-gray-600">
                                                            Вы уверены, что хотите удалить приложение "{app.name}"? 
                                                            Это действие нельзя отменить, и API ключ перестанет работать.
                                                        </AlertDialogDescription>
                                                    </AlertDialogHeader>
                                                    <AlertDialogFooter>
                                                        <AlertDialogCancel className="bg-gray-100 text-gray-900 hover:bg-gray-200">Отмена</AlertDialogCancel>
                                                        <AlertDialogAction
                                                            onClick={() => handleDeleteApp(app)}
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
                                        <div className="space-y-2">
                                            <div className="text-sm text-muted-foreground">
                                                <p>ID: {app.id}</p>
                                                <p>Создано: {formatDate(app.created_at)}</p>
                                            </div>
                                            {app.roles && app.roles.length > 0 && (
                                                <div>
                                                    <p className="text-sm font-medium mb-1">Разрешения:</p>
                                                    <div className="flex flex-wrap gap-1">
                                                        {app.roles.map((role) => (
                                                            <span
                                                                key={role}
                                                                className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                                                            >
                                                                {role === 'external_app_reader' ? 'Чтение' : 
                                                                 role === 'external_app_writer' ? 'Запись' : role}
                                                            </span>
                                                        ))}
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
} 