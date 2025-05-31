import { createFileRoute, Outlet } from '@tanstack/react-router';
import { SettingsLayout } from 'src/features/userSettings/ui/SettingsLayout';

function SettingsRoute() {
    return (
        <SettingsLayout>
            <Outlet />
        </SettingsLayout>
    );
}

export const Route = createFileRoute('/panel/settings')({
    component: SettingsRoute,
}); 