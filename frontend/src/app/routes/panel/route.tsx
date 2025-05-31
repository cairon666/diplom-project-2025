import { createFileRoute } from '@tanstack/react-router';
import { PanelPage } from 'src/pages/PanelPage';

export const Route = createFileRoute('/panel')({
    component: PanelPage,
});
