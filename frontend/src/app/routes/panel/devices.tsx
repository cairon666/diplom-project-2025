import { createFileRoute } from '@tanstack/react-router';
import { DevicesPage } from 'src/pages/DevicesPage';

export const Route = createFileRoute('/panel/devices')({
    component: DevicesPage,
}); 