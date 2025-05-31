import { useCallback } from 'react';
import { toast } from 'sonner';
import { ApiErrorCode, isBaseError } from 'src/store';
import { useAppDispatch } from 'src/store';
import { Button } from 'ui/button';

import { 
    userSettingApi, 
    useSettingQuery, 
    useUpdateTelegramSettingMutation, 
    useUnlinkTelegramSettingMutation 
} from '../api';
import { TelegramInfo } from '../types';
import { TelegramWidget } from './TelegramWidget';

export function TelegramSetting() {
    const dispatch = useAppDispatch();
    const { data } = useSettingQuery();
    const [updateTelegram, { isLoading: isUpdating }] = useUpdateTelegramSettingMutation();
    const [unlinkTelegram, { isLoading: isUnlinking }] = useUnlinkTelegramSettingMutation();

    const handleTelegramAuth = useCallback((user: TelegramInfo) => {
        updateTelegram({
            id: user.id,
            first_name: user.first_name,
            last_name: user.last_name,
            username: user.username,
            photo_url: user.photo_url,
            auth_date: user.auth_date,
            hash: user.hash
        })
            .unwrap()
            .then(() => {
                dispatch(
                    userSettingApi.util.updateQueryData(
                        'setting',
                        undefined,
                        (draft) => {
                            draft.has_telegram = true;
                        },
                    )
                );
                toast.success('Telegram успешно привязан к аккаунту');
            })
            .catch((error: unknown) => {
                if (isBaseError(error)) {
                    if (error.data.error === ApiErrorCode.CodeInvalidTelegramHash) {
                        toast.error('Хэш Telegram недействителен, пожалуйста, попробуйте снова');
                    } else if (error.data.error === ApiErrorCode.CodeProviderAlreadyConnected) {
                        toast.error('Этот Telegram аккаунт уже привязан к другому пользователю');
                    } else if (error.data.error === ApiErrorCode.CodeProviderAccountAlreadyLinked) {
                        toast.error('Telegram аккаунт уже привязан');
                    } else {
                        toast.error('Произошла ошибка при привязке Telegram');
                    }
                } else {
                    toast.error('Произошла ошибка при привязке Telegram');
                }
            });
    }, [updateTelegram, dispatch]);

    const handleUnlinkTelegram = () => {
        unlinkTelegram()
            .unwrap()
            .then(() => {
                dispatch(
                    userSettingApi.util.updateQueryData(
                        'setting',
                        undefined,
                        (draft) => {
                            draft.has_telegram = false;
                        },
                    )
                );
                toast.success('Telegram отвязан от аккаунта');
            })
            .catch(() => {
                toast.error('Произошла ошибка при отвязке Telegram');
            });
    };

    if (!data) {
        return null;
    }

    return (
        <div className="flex flex-col gap-4 w-full">
            <div className="flex items-center justify-between">
                <div className="flex flex-col">
                    <h3 className="text-lg font-medium">Telegram</h3>
                    <p className="text-sm text-muted-foreground">
                        {data.has_telegram 
                            ? 'Telegram аккаунт привязан к вашему профилю' 
                            : 'Telegram аккаунт не привязан'
                        }
                    </p>
                </div>
                
                {data.has_telegram ? (
                    <Button 
                        variant="destructive" 
                        onClick={handleUnlinkTelegram}
                        disabled={isUnlinking}
                    >
                        {isUnlinking ? 'Отвязываем...' : 'Отвязать Telegram'}
                    </Button>
                ) : (
                    <TelegramWidget
                        onAuth={handleTelegramAuth}
                        isLoading={isUpdating}
                        buttonText="Привязать Telegram"
                    />
                )}
            </div>
        </div>
    );
}
