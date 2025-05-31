import { createFileRoute } from '@tanstack/react-router';
import { Login } from 'src/features/auth';

export const Route = createFileRoute('/auth/login')({
    component: Login,
});
