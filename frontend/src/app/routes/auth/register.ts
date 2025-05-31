import { createFileRoute } from '@tanstack/react-router';
import { Register } from 'src/features/auth';

export const Route = createFileRoute('/auth/register')({
    component: Register,
});
