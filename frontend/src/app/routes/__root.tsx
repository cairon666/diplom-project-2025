import { Outlet, createRootRoute, useNavigate } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';
import { useEffect } from 'react';
import { Route as loginRouter } from 'src/app/routes/auth/login';
import { selectCurrentUserId } from 'src/entities/user';
import { useAppSelector } from 'src/store';

function RootComponent() {
    const navigate = useNavigate();
    const currentUserId = useAppSelector(selectCurrentUserId);

    useEffect(() => {
        if (!currentUserId) {
            navigate({ to: loginRouter.to });
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <>
            <Outlet />
            {process.env.NODE_ENV === 'development' && (
                <TanStackRouterDevtools />
            )}
        </>
    );
}

export const Route = createRootRoute({
    component: RootComponent,
});
