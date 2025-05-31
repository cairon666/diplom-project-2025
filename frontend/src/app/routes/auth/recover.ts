import { createFileRoute } from '@tanstack/react-router';
import { Recover } from 'src/features/auth';

export const Route = createFileRoute('/auth/recover')({
    component: Recover,
});
