import { Outlet } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import {
    selectCurrentUser,
    setCurrentUser,
    useGetUserQuery,
} from 'src/entities/user';
import { useAppDispatch } from 'src/store';
import { LoadingSpinner } from 'ui/loading-spinner';

import { Header } from './Header';

export function PanelPage() {
    const dispatch = useAppDispatch();
    const currentUser = useSelector(selectCurrentUser);

    const { isLoading, data } = useGetUserQuery();

    useEffect(() => {
        if (data) {
            dispatch(setCurrentUser(data));
        }
    }, [data, dispatch]);

    if (isLoading || !currentUser) {
        return <LoadingSpinner />;
    }

    return (
        <div className="w-full h-screen flex flex-col">
            <Header />
            <div className="flex flex-1">
                <Outlet />
            </div>
        </div>
    );
}
