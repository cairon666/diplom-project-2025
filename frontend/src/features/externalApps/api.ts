import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

export interface ExternalAppListItem {
    id: string;
    owner_id: string;
    name: string;
    created_at: string;
    roles: string[];
}

export interface GetExternalAppsListResponse {
    external_apps: ExternalAppListItem[];
}

export interface CreateExternalAppRequest {
    name: string;
    roles: string[];
}

export interface CreateExternalAppResponse {
    api_key: string;
    id_external_app: string;
}

export interface DeleteExternalAppRequest {
    id: string;
}

export const externalAppsApi = createApi({
    reducerPath: 'externalAppsApi',
    baseQuery: baseQueryWithReauth,
    tagTypes: ['ExternalApp'],
    endpoints: (builder) => ({
        getExternalAppsList: builder.query<GetExternalAppsListResponse, void>({
            query: () => ({
                url: `/v1/user/external-apps`,
                method: 'GET',
            }),
            providesTags: ['ExternalApp'],
        }),
        createExternalApp: builder.mutation<CreateExternalAppResponse, CreateExternalAppRequest>({
            query: (data) => ({
                url: `/v1/user/external-apps`,
                method: 'POST',
                body: data,
            }),
            invalidatesTags: ['ExternalApp'],
        }),
        deleteExternalApp: builder.mutation<unknown, DeleteExternalAppRequest>({
            query: ({ id }) => ({
                url: `/v1/user/external-apps/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: ['ExternalApp'],
        }),
    }),
});

export const {
    useGetExternalAppsListQuery,
    useCreateExternalAppMutation,
    useDeleteExternalAppMutation,
} = externalAppsApi; 