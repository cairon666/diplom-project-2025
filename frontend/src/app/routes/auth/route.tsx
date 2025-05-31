import { createFileRoute } from '@tanstack/react-router';
import { AuthPage } from 'src/pages/AuthPage';

export const Route = createFileRoute('/auth')({
    component: AuthPage,
});
