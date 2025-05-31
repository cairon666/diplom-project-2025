import { createFileRoute, Navigate } from '@tanstack/react-router';

function RouteComponent() {
    return <Navigate to="/panel" />;
}

export const Route = createFileRoute('/')({
    component: RouteComponent,
});
