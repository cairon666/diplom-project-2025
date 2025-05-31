import { useNavigate, useSearch } from '@tanstack/react-router';
import { useCallback } from 'react';
import { Route as panelRouter } from 'src/app/routes/panel/route';
import { setAccessToken, setUserId } from 'src/entities/user';
import { useAppDispatch } from 'src/store';

import { LoginResponse } from '../api';

export function useOnSuccessLogin() {
    const search: { redirect: string } = useSearch({ strict: false });
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const onSuccessLogin = useCallback(
        (loginData: LoginResponse) => {
            dispatch(setUserId(loginData.id));
            dispatch(setAccessToken(loginData.access_token));
            navigate({ to: search.redirect || panelRouter.to });
        },
        [dispatch, navigate, search.redirect],
    );

    return { onSuccessLogin };
}
