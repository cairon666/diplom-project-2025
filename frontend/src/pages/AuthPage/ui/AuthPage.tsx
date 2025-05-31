import { Navigate, Outlet } from '@tanstack/react-router';
import { Route as panelRouter } from 'src/app/routes/panel/route';
import { selectCurrentUserId } from 'src/entities/user';
import { useAppSelector } from 'src/store';

import styles from './AuthPage.module.scss';

export function AuthPage() {
    const currentUserId = useAppSelector(selectCurrentUserId);

    if (currentUserId) {
        return <Navigate to={panelRouter.to} />;
    }

    return (
        <div className={styles.authLayout}>
            <div>123</div>
            <div className={styles.authLayoutForm}>
                <Outlet />
            </div>
        </div>
    );
}
