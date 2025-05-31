import { LoadingSpinner } from 'ui/loading-spinner';
import { Separator } from 'ui/separator';

import { useSettingQuery } from '../api';
import { EmailSetting } from './EmailSetting';
import { PasswordSetting } from './PasswordSetting';
import { ProfileSetting } from './ProfileSetting';
import { TelegramSetting } from './TelegramSetting';

export function GeneralSettings() {
    const { data, isLoading } = useSettingQuery();

    if (isLoading || !data) {
        return (
            <div className="flex-1 p-6">
                <div className="max-w-7xl mx-auto">
                    <div className="flex items-center gap-2">
                        <LoadingSpinner />
                        <span>Загрузка настроек...</span>
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
                        <h2 className="text-2xl font-bold mb-2">Основные настройки</h2>
                        <p className="text-muted-foreground mb-6">
                            Управление основной информацией профиля, безопасностью и подключенными сервисами.
                        </p>
                    </div>
                    
                    <div className="flex flex-col w-full gap-4 items-center w-full">
                        <ProfileSetting />
                        <Separator />
                        <EmailSetting />
                        <Separator />
                        <PasswordSetting />
                        <Separator />
                        <TelegramSetting />
                    </div>
                </div>
            </div>
        </div>
    );
} 