import { createFileRoute } from '@tanstack/react-router';
import { GeneralSettings } from 'src/features/userSettings/ui/GeneralSettings';

export const Route = createFileRoute('/panel/settings/')({
    component: GeneralSettings,
}); 