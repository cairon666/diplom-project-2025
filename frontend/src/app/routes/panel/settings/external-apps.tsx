import { createFileRoute } from '@tanstack/react-router';
import { ExternalAppsList } from 'src/features/externalApps';

export const Route = createFileRoute('/panel/settings/external-apps')({
    component: ExternalAppsList,
}); 