import { configureStore } from '@reduxjs/toolkit';
import { authSlice, currentUserSlice, currentUserApi } from 'src/entities/user';
import { authApi } from 'src/features/auth';
import { userSettingApi } from 'src/features/userSettings';
import { deviceApi } from 'src/features/devices';
import { externalAppsApi } from 'src/features/externalApps';
import { rrIntervalsApi, rrAnalyticsApi } from 'src/features/health';

export const store = configureStore({
    reducer: {
        [authSlice.reducerPath]: authSlice.reducer,
        [authApi.reducerPath]: authApi.reducer,
        [currentUserSlice.reducerPath]: currentUserSlice.reducer,
        [currentUserApi.reducerPath]: currentUserApi.reducer,
        [userSettingApi.reducerPath]: userSettingApi.reducer,
        [deviceApi.reducerPath]: deviceApi.reducer,
        [externalAppsApi.reducerPath]: externalAppsApi.reducer,
        [rrIntervalsApi.reducerPath]: rrIntervalsApi.reducer,
        [rrAnalyticsApi.reducerPath]: rrAnalyticsApi.reducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(
            authApi.middleware,
            currentUserApi.middleware,
            userSettingApi.middleware,
            deviceApi.middleware,
            externalAppsApi.middleware,
            rrIntervalsApi.middleware,
            rrAnalyticsApi.middleware,
        ),
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch;
