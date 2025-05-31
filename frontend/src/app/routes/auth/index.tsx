import { createFileRoute, Navigate } from '@tanstack/react-router';

function AuthIndexPage() {
    return <Navigate to="/auth/login" />;
}

export const Route = createFileRoute('/auth/')({
    component: AuthIndexPage,
}); 