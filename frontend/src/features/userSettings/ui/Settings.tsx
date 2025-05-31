import { GeneralSettings } from './GeneralSettings';
import { SettingsLayout } from './SettingsLayout';

export function Settings() {
    return (
        <SettingsLayout>
            <GeneralSettings />
        </SettingsLayout>
    );
}

// Экспорт с опечаткой для обратной совместимости
export { Settings as Settigs };
