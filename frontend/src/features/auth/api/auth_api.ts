import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    access_token: string;
    id: string;
}

export interface RegisterRequest {
    first_name: string;
    second_name: string;
    email: string;
    password: string;
}

export interface TelegramAuthInfo {
    auth_date: number;
    first_name: string;
    last_name: string;
    username: string;
    hash: string;
    photo_url: string;
    id: number;
}

export interface TelegramConfirmMutationRequest {
    temp_id: string;
    email: string;
    first_name: string;
    last_name: string;
}

export const authApi = createApi({
    reducerPath: 'authApi',
    baseQuery: baseQueryWithReauth,
    endpoints: (builder) => ({
        login: builder.mutation<LoginResponse, LoginRequest>({
            query: (credentials) => ({
                url: '/v1/auth/login',
                method: 'POST',
                body: credentials,
            }),
        }),
        register: builder.mutation<unknown, RegisterRequest>({
            query: (credentials) => ({
                url: '/v1/auth/register',
                method: 'POST',
                body: credentials,
            }),
        }),
        telegramLogin: builder.mutation<LoginResponse, TelegramAuthInfo>({
            query: (credentials) => ({
                url: '/v1/auth/telegram/login',
                method: 'POST',
                body: credentials,
            }),
        }),
        telegramConfirm: builder.mutation<
            LoginResponse,
            TelegramConfirmMutationRequest
        >({
            query: (credentials) => ({
                url: '/v1/auth/telegram/confirm-register',
                method: 'POST',
                body: credentials,
            }),
        }),
    }),
});

export const {
    useLoginMutation,
    useRegisterMutation,
    useTelegramConfirmMutation,
    useTelegramLoginMutation,
} = authApi;
