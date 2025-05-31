import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

export interface GetSettingResponse {
    email: string;
    has_password: boolean;
    has_telegram: boolean;
    first_name: string;
    last_name: string;
}

export interface GetUserByIdResponse {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
}

export interface UpdateProfileSettingRequest {
    first_name: string;
    last_name: string;
}

export interface UpdateEmailSettingRequest {
    email: string;
}

export interface UpdatePasswordSettingRequest {
    password: string;
}

export interface UpdateTelegramSettingRequest {
    id: number;
    first_name: string;
    last_name: string;
    username: string;
    photo_url: string;
    auth_date: number;
    hash: string;
}

export const userSettingApi = createApi({
    reducerPath: 'userSettingApi',
    baseQuery: baseQueryWithReauth,
    endpoints: (builder) => ({
        getUserById: builder.query<GetUserByIdResponse, void>({
            query: () => ({
                url: `/v1/user`,
                method: 'GET',
            }),
        }),
        setting: builder.query<GetSettingResponse, void>({
            query: () => ({
                url: `/v1/user/setting`,
                method: 'GET',
            }),
        }),
        updateProfileSetting: builder.mutation<
            unknown,
            UpdateProfileSettingRequest
        >({
            query: (data) => ({
                url: `/v1/user/setting/profile`,
                method: 'PATCH',
                body: data,
            }),
        }),
        updateEmailSetting: builder.mutation<
            unknown,
            UpdateEmailSettingRequest
        >({
            query: (data) => ({
                url: `/v1/user/setting/email`,
                method: 'PATCH',
                body: data,
            }),
        }),
        updatePasswordSetting: builder.mutation<
            unknown,
            UpdatePasswordSettingRequest
        >({
            query: (data) => ({
                url: `/v1/user/setting/password`,
                method: 'PATCH',
                body: data,
            }),
        }),
        updateTelegramSetting: builder.mutation<
            unknown,
            UpdateTelegramSettingRequest
        >({
            query: (data) => ({
                url: `/v1/user/setting/telegram`,
                method: 'PATCH',
                body: data,
            }),
        }),
        unlinkTelegramSetting: builder.mutation<unknown, void>({
            query: () => ({
                url: `/v1/user/setting/telegram`,
                method: 'DELETE',
            }),
        }),
    }),
});

export const {
    useGetUserByIdQuery,
    useSettingQuery,
    useUpdateProfileSettingMutation,
    useUpdateEmailSettingMutation,
    useUpdatePasswordSettingMutation,
    useUpdateTelegramSettingMutation,
    useUnlinkTelegramSettingMutation,
} = userSettingApi;
