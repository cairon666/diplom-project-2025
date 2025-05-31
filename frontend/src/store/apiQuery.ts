import {
    BaseQueryFn,
    FetchArgs,
    fetchBaseQuery,
    FetchBaseQueryError,
} from '@reduxjs/toolkit/query/react';
import { router } from 'src/app/App';
import { logout, setAccessToken } from 'src/entities/user';
import { RootState } from 'src/store';

import { AppBaseError } from './apiErrors';

const baseQuery = fetchBaseQuery({
    baseUrl: '/api',
    prepareHeaders: (headers, { getState }) => {
        const { accessToken } = (getState() as RootState).auth;
        if (accessToken) {
            headers.set('authorization', `Bearer ${accessToken}`);
        }
        return headers;
    },
}) as BaseQueryFn<
    string | FetchArgs,
    unknown,
    FetchBaseQueryError & { data: AppBaseError }
>;

interface RefreshResponse {
    access_token: string;
    refresh_token: string;
}

export const baseQueryWithReauth: BaseQueryFn<
    string | FetchArgs,
    unknown,
    FetchBaseQueryError
> = async (args, api, extraOptions) => {
    let result = await baseQuery(args, api, extraOptions);

    if (result.error && result.error.status === 401) {
        const refreshResult = await baseQuery(
            {
                url: '/v1/auth/refresh',
                method: 'POST',
            },
            api,
            extraOptions,
        );
        const refreshResponse = refreshResult.data as
            | RefreshResponse
            | undefined;

        if (refreshResponse) {
            api.dispatch(setAccessToken(refreshResponse.access_token));
            result = await baseQuery(args, api, extraOptions);
        } else {
            api.dispatch(logout());
            router.navigate({ to: '/auth/login' });
        }
    }

    return result;
};
