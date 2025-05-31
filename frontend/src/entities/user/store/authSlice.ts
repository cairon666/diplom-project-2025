import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import {
    clearUserIdLocalStorage,
    clearAccessTokenLocalStorage,
    getAccessTokenLocalStorage,
    getUserIdLocalStorage,
    setAccessTokenLocalStorage,
    setUserIdLocalStorage,
} from '../lib/localStorage';

export interface AuthState {
    user_id: string | null;
    accessToken: string | null;
}

const initialState: AuthState = {
    user_id: getUserIdLocalStorage(),
    accessToken: getAccessTokenLocalStorage(),
};

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setAccessToken: (state, action: PayloadAction<string>) => {
            setAccessTokenLocalStorage(action.payload);
            state.accessToken = action.payload;
        },
        setUserId: (state, action: PayloadAction<string>) => {
            setUserIdLocalStorage(action.payload);
            state.user_id = action.payload;
        },
        logout: (state) => {
            clearUserIdLocalStorage();
            clearAccessTokenLocalStorage();
            state.accessToken = null;
            state.user_id = null;
        },
    },
    selectors: {
        selectCurrentUserId: (state) => state.user_id!,
    },
});

export const { setAccessToken, setUserId, logout } = authSlice.actions;
export const { selectCurrentUserId } = authSlice.selectors;
