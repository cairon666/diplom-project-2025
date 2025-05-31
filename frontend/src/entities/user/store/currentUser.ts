import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { IUser } from '../model';

export interface CurrentUserState {
    user: IUser | null;
}

const initialState: CurrentUserState = {
    user: null,
};

export const currentUserSlice = createSlice({
    name: 'currentUser',
    initialState,
    reducers: {
        setCurrentUser: (state, action: PayloadAction<IUser>) => {
            state.user = action.payload;
        },
        setPartialCurrentUser: (
            state,
            action: PayloadAction<Partial<IUser>>,
        ) => {
            if (state.user) {
                Object.entries(action.payload).forEach(([key, value]) => {
                    // @ts-ignore
                    state.user[key] = value;
                });
            }
        },
    },
    selectors: {
        selectCurrentUser: (state) => state.user!,
    },
});

export const { setCurrentUser, setPartialCurrentUser } =
    currentUserSlice.actions;
export const { selectCurrentUser } = currentUserSlice.selectors;
