import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

import { IUser } from '../model';

export const currentUserApi = createApi({
    reducerPath: 'currentUserApi',
    baseQuery: baseQueryWithReauth,
    endpoints: (builder) => ({
        getUser: builder.query<IUser, void>({
            query: () => ({
                url: `/v1/user`,
                method: 'GET',
            }),
        }),
    }),
});

export const { useGetUserQuery } = currentUserApi;
